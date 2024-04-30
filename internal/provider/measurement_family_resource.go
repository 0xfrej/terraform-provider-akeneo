package provider

import (
	"context"
	"fmt"
	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MeasurementFamilyResource{}
var _ resource.ResourceWithImportState = &MeasurementFamilyResource{}
var _ resource.ResourceWithConfigure = &MeasurementFamilyResource{}

func NewMeasurementFamilyResource() resource.Resource {
	return &MeasurementFamilyResource{}
}

// MeasurementFamilyResource defines the resource implementation.
type MeasurementFamilyResource struct {
	client *akeneox.MeasurementFamilyService
}

//TODO: use map for units just as the api does because akeneo returns objects in unpredictable order and tf alwasy thinks there is a change

type MeasurementFamilyResourceUnitConversionModel struct {
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

type MeasurementFamilyResourceUnitModel struct {
	Code                types.String                                   `tfsdk:"code"`
	Symbol              types.String                                   `tfsdk:"symbol"`
	Labels              types.Map                                      `tfsdk:"labels"`
	ConvertFromStandard []MeasurementFamilyResourceUnitConversionModel `tfsdk:"convert_from_standard"`
}

// MeasurementFamilyResourceModel describes the resource data model.
type MeasurementFamilyResourceModel struct {
	Code             types.String                         `tfsdk:"code"`
	StandardUnitCode types.String                         `tfsdk:"standard_unit_code"`
	Labels           types.Map                            `tfsdk:"labels"`
	Units            []MeasurementFamilyResourceUnitModel `tfsdk:"units"`
}

func (r *MeasurementFamilyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_measurement_family"
}

func (r *MeasurementFamilyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo measurement family resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Measurement family code (preferred uppercase values to follow Akeneo's convention)",
				Required:    true,
			},
			"standard_unit_code": schema.StringAttribute{
				Description: "Unit code used as the standard unit for this measurement family",
				Required:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Label definition per locale",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidatorx.IsLocaleCode()),
				},
			},
			"units": schema.ListNestedAttribute{
				Description: "Unit definitions",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Description: "Measurement unit code.",
							Required:    true,
						},
						"labels": schema.MapAttribute{
							Description: "Label definition per locale",
							Optional:    true,
							ElementType: types.StringType,
							Validators: []validator.Map{
								mapvalidator.KeysAre(stringvalidatorx.IsLocaleCode()),
							},
						},
						"symbol": schema.StringAttribute{
							Description: "Measurement unit symbol.",
							Required:    true,
						},
						"convert_from_standard": schema.ListNestedAttribute{
							Description: "Calculation to convert the unit from the standard unit.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"operator": schema.StringAttribute{
										Description: "The operator for a conversion operation to convert a unit from the standard unit.",
										Required:    true,
										Validators: []validator.String{
											stringvalidatorx.IsPimConversionOperator(),
										},
									},
									"value": schema.StringAttribute{
										Description: "The value for a conversion operation to convert the unit from the standard unit.",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^-?\d*(\.\d+)?$`),
												"must only contain decimal values",
											),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

func (r *MeasurementFamilyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*ResourceData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ResourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if data.Client == nil {
		resp.Diagnostics.AddError(
			"Missing client instance",
			"Client instance pointer passed to Configure is required, got nil",
		)
		return
	}

	r.client = akeneox.NewMeasurementFamilyClient(data.Client)
}

func (r *MeasurementFamilyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MeasurementFamilyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateMeasurementFamilies([]akeneox.MeasurementFamily{*apiData})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a measurement family",
			"An unexpected error occurred when creating measurement family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if result != nil {
		for _, r := range *result {
			if r.StatusCode > 299 {
				if r.Message != "" {
					resp.Diagnostics.AddError(
						"Error while creating a measurement family",
						"An unexpected error occurred when creating measurement family. \n\n"+
							"Akeneo API Error: "+r.Message,
					)
				}

				for _, e := range r.Errors {
					resp.Diagnostics.AddError(
						"Error while creating a measurement family",
						"A validation error was returned from the Akeneo API. \n\n"+
							"Validation Error: "+e.Message+"\n"+
							"On property: "+e.Property+"\n",
					)
				}
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeasurementFamilyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MeasurementFamilyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrData, err := r.client.GetMeasurementFamily(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading a measurement family",
			"An unexpected error occurred when reading measurement family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	r.mapToTfObject(&resp.Diagnostics, &data, attrData)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeasurementFamilyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MeasurementFamilyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateMeasurementFamilies([]akeneox.MeasurementFamily{*apiData})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating a measurement family",
			"An unexpected error occurred when updating measurement family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if result != nil {
		for _, r := range *result {
			if r.StatusCode > 299 {
				if r.Message != "" {
					resp.Diagnostics.AddError(
						"Error while updating a measurement family",
						"An unexpected error occurred when creating measurement family. \n\n"+
							"Akeneo API Error: "+r.Message,
					)
				}

				for _, e := range r.Errors {
					resp.Diagnostics.AddError(
						"Error while updating a measurement family",
						"A validation error was returned from the Akeneo API. \n\n"+
							"Validation Error: "+e.Message+"\n"+
							"On property: "+e.Property+"\n",
					)
				}
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeasurementFamilyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data MeasurementFamilyResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for measurement family.",
	)
}

func (r *MeasurementFamilyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *MeasurementFamilyResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *MeasurementFamilyResourceModel) *akeneox.MeasurementFamily {
	a := akeneox.MeasurementFamily{
		Code:             data.Code.ValueString(),
		StandardUnitCode: data.StandardUnitCode.ValueString(),
	}

	if !(data.Labels.IsNull() || data.Labels.IsUnknown()) {
		elements := make(map[string]types.String, len(data.Labels.Elements()))
		diags.Append(data.Labels.ElementsAs(ctx, &elements, false)...)
		labels := make(map[string]string)
		for locale, label := range elements {
			labels[locale] = label.ValueString()
		}
		a.Labels = labels
	}

	units := make(map[string]akeneox.MeasurementUnit, len(data.Units))
	for _, unit := range data.Units {
		u := akeneox.MeasurementUnit{
			Code:   unit.Code.ValueString(),
			Symbol: unit.Symbol.ValueString(),
		}

		if !(unit.Labels.IsNull() || unit.Labels.IsUnknown()) {
			elements := make(map[string]types.String, len(unit.Labels.Elements()))
			diags.Append(unit.Labels.ElementsAs(ctx, &elements, false)...)
			labels := make(map[string]string)
			for locale, label := range elements {
				labels[locale] = label.ValueString()
			}
			u.Labels = labels
		}

		conversions := make([]akeneox.MeasurementUnitConversion, len(unit.ConvertFromStandard))
		for i, conversion := range unit.ConvertFromStandard {
			conversions[i] = akeneox.MeasurementUnitConversion{
				Operator: conversion.Operator.ValueString(),
				Value:    conversion.Value.ValueString(),
			}
		}
		u.ConvertFromStandard = conversions

		units[u.Code] = u
	}
	a.Units = units

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *MeasurementFamilyResource) mapToTfObject(respDiags *diag.Diagnostics, data *MeasurementFamilyResourceModel, apiData *akeneox.MeasurementFamily) {
	data.Code = types.StringValue(apiData.Code)
	data.StandardUnitCode = types.StringValue(apiData.StandardUnitCode)

	if apiData.Labels != nil {
		elements := make(map[string]attr.Value, len(apiData.Labels))

		for k, v := range apiData.Labels {
			elements[k] = types.StringValue(v)
		}

		mapVal, diags := types.MapValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Labels = mapVal
	}

	if len(apiData.Units) > 0 {
		units := make([]MeasurementFamilyResourceUnitModel, len(apiData.Units))
		i := 0
		for _, unit := range apiData.Units {
			u := MeasurementFamilyResourceUnitModel{
				Code:   types.StringValue(unit.Code),
				Symbol: types.StringValue(unit.Symbol),
			}

			if len(unit.Labels) > 0 {
				elements := make(map[string]attr.Value, len(unit.Labels))
				for k, v := range unit.Labels {
					elements[k] = types.StringValue(v)
				}
				mapVal, diags := types.MapValue(types.StringType, elements)
				if diags.HasError() {
					respDiags.Append(diags...)
				}
				u.Labels = mapVal
			}

			if len(unit.ConvertFromStandard) > 0 {
				conversions := make([]MeasurementFamilyResourceUnitConversionModel, len(unit.ConvertFromStandard))
				for i, conversion := range unit.ConvertFromStandard {
					c := MeasurementFamilyResourceUnitConversionModel{
						Operator: types.StringValue(conversion.Operator),
						Value:    types.StringValue(conversion.Value),
					}
					conversions[i] = c
				}
				u.ConvertFromStandard = conversions
			}

			units[i] = u
			i++
		}
		data.Units = units
	}
}
