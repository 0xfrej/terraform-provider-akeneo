package provider

import (
	"context"
	"fmt"
	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
	goakeneo "github.com/ezifyio/go-akeneo"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FamilyResource{}
var _ resource.ResourceWithImportState = &FamilyResource{}
var _ resource.ResourceWithConfigure = &FamilyResource{}

func NewFamilyResource() resource.Resource {
	return &FamilyResource{}
}

// FamilyResource defines the resource implementation.
type FamilyResource struct {
	client *akeneox.FamilyService
}

// FamilyResourceModel describes the resource data model.
type FamilyResourceModel struct {
	Code                  types.String `tfsdk:"code"`
	Labels                types.Map    `tfsdk:"labels"`
	Attributes            types.List   `tfsdk:"attributes"`
	AttributeAsLabel      types.String `tfsdk:"attribute_as_label"`
	AttributeAsImage      types.String `tfsdk:"attribute_as_image"`
	AttributeRequirements types.Map    `tfsdk:"attribute_requirements"`
}

func (r *FamilyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_family"
}

func (r *FamilyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo family resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Family code",
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
			"attributes": schema.ListAttribute{
				Description: "Attributes assigned to the family",
				Optional:    true,
				ElementType: types.StringType,
			},
			"attribute_as_label": schema.StringAttribute{
				Description: "Attribute used as product label for the family",
				Optional:    true,
			},
			"attribute_as_image": schema.StringAttribute{
				Description: "Attribute used as product image for the family",
				Optional:    true,
			},
			"attribute_requirements": schema.MapAttribute{
				Description: "Attribute codes of the family that are required for the completeness calculation for each channel.",
				Optional:    true,
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
			},
		},
	}
}

func (r *FamilyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = akeneox.NewFamilyClient(data.Client)
}

func (r *FamilyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FamilyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateFamily(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a family",
			"An unexpected error occurred when creating family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FamilyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FamilyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData, err := r.client.GetFamily(data.Code.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a family",
			"An unexpected error occurred when reading family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	r.mapToTfObject(&resp.Diagnostics, &data, apiData)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FamilyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FamilyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateFamily(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating a family",
			"An unexpected error occurred when updating family. \n\n"+
				"Akeneo API Error: "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FamilyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data FamilyResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for attributes.",
	)
}

func (r *FamilyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *FamilyResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *FamilyResourceModel) *goakeneo.Family {
	a := goakeneo.Family{
		Code: data.Code.ValueString(),
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

	if !(data.Attributes.IsNull() || data.Attributes.IsUnknown()) {
		elements := make([]types.String, 0, len(data.Attributes.Elements()))
		diags.Append(data.Attributes.ElementsAs(ctx, &elements, false)...)
		attrs := make([]string, len(elements))
		for i, ext := range elements {
			attrs[i] = ext.ValueString()
		}
		a.Attributes = attrs
	}

	if !(data.AttributeAsLabel.IsNull() || data.AttributeAsLabel.IsUnknown()) {
		a.AttributeAsLabel = data.AttributeAsLabel.ValueString()
	}

	if !(data.AttributeAsImage.IsNull() || data.AttributeAsImage.IsUnknown()) {
		a.AttributeAsImage = data.AttributeAsImage.ValueString()
	}

	if !(data.AttributeRequirements.IsNull() || data.AttributeRequirements.IsUnknown()) {
		elements := make(map[string]types.List, len(data.AttributeRequirements.Elements()))
		diags.Append(data.AttributeRequirements.ElementsAs(ctx, &elements, false)...)
		reqs := make(map[string][]string)
		for channel, attrs := range elements {
			innerElements := make([]types.String, 0, len(attrs.Elements()))
			diags.Append(attrs.ElementsAs(ctx, &innerElements, false)...)
			reqAttrs := make([]string, len(innerElements))
			for i, a := range innerElements {
				reqAttrs[i] = a.ValueString()
			}
			reqs[channel] = reqAttrs
		}
		a.AttributeRequirements = reqs
	}

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *FamilyResource) mapToTfObject(respDiags *diag.Diagnostics, data *FamilyResourceModel, apiData *goakeneo.Family) {
	// todo map only non empties
	data.Code = types.StringValue(apiData.Code)

	if len(apiData.Labels) > 0 {
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

	if len(apiData.Attributes) > 0 {
		elements := make([]attr.Value, len(apiData.Attributes))

		for k, v := range apiData.Attributes {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Attributes = listVal
	}

	data.AttributeAsLabel = types.StringValue(apiData.AttributeAsLabel)
	data.AttributeAsImage = types.StringValue(apiData.AttributeAsImage)

	if len(apiData.AttributeRequirements) > 0 {
		elements := make(map[string]attr.Value, len(apiData.AttributeRequirements))

		for k, v := range apiData.AttributeRequirements {
			innerElements := make([]attr.Value, len(v))
			for i, a := range v {
				innerElements[i] = types.StringValue(a)
			}
			listVal, diags := types.ListValue(types.StringType, innerElements)
			if diags.HasError() {
				respDiags.Append(diags...)
			}
			elements[k] = listVal
		}

		mapVal, diags := types.MapValue(types.ListType{ElemType: types.StringType}, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.AttributeRequirements = mapVal
	}
}
