package provider

import (
	"context"
	"fmt"
	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
	goakeneo "github.com/ezifyio/go-akeneo"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
var _ resource.Resource = &FamilyVariantResource{}
var _ resource.ResourceWithImportState = &FamilyVariantResource{}
var _ resource.ResourceWithConfigure = &FamilyVariantResource{}

func NewFamilyVariantResource() resource.Resource {
	return &FamilyVariantResource{}
}

// FamilyVariantResource defines the resource implementation.
type FamilyVariantResource struct {
	client *akeneox.FamilyService
}

type VariantAttributeSetModel struct {
	Level      types.Int64 `tfsdk:"level"`
	Axes       types.List  `tfsdk:"axes"`
	Attributes types.List  `tfsdk:"attributes"`
}

// FamilyVariantResourceModel describes the resource data model.
type FamilyVariantResourceModel struct {
	FamilyCode           types.String               `tfsdk:"family_code"`
	Code                 types.String               `tfsdk:"code"`
	Labels               types.Map                  `tfsdk:"labels"`
	VariantAttributeSets []VariantAttributeSetModel `tfsdk:"variant_attribute_sets"`
}

func (r *FamilyVariantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_family_variant"
}

func (r *FamilyVariantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo family variant resource",

		Attributes: map[string]schema.Attribute{
			"family_code": schema.StringAttribute{
				Description: "Family code to which this variant belongs",
				Required:    true,
			},
			"code": schema.StringAttribute{
				Description: "Family variant code",
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
			"variant_attribute_sets": schema.ListNestedAttribute{
				Description: "Attribute distributions according to the enrichment level.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"level": schema.Int64Attribute{
							Description: "Enrichment level",
							Required:    true,
							Validators: []validator.Int64{
								int64validator.Between(1, 2147483647),
							},
						},
						"axes": schema.ListAttribute{
							Description: "Codes of attributes used as variant axes",
							ElementType: types.StringType,
							Optional:    true,
						},
						"attributes": schema.ListAttribute{
							Description: "Codes of attributes bind to this enrichment level",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *FamilyVariantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FamilyVariantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FamilyVariantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateOrCreate(data.FamilyCode.ValueString(), data.Code.ValueString(), *apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a family variant",
			"An unexpected error occurred when creating family variant. \n\n"+
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

func (r *FamilyVariantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FamilyVariantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData, err := r.client.GetFamilyVariant(data.FamilyCode.ValueString(), data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a family variant",
			"An unexpected error occurred when reading family variant. \n\n"+
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

func (r *FamilyVariantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FamilyVariantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateOrCreate(data.FamilyCode.ValueString(), data.Code.ValueString(), *apiData)
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

func (r *FamilyVariantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data FamilyVariantResourceModel
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

func (r *FamilyVariantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *FamilyVariantResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *FamilyVariantResourceModel) *goakeneo.FamilyVariant {
	a := goakeneo.FamilyVariant{
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

	sets := make([]goakeneo.VariantAttributeSet, len(data.VariantAttributeSets))
	for i, set := range data.VariantAttributeSets {
		s := goakeneo.VariantAttributeSet{
			Level: int(set.Level.ValueInt64()),
		}

		if !set.Axes.IsNull() {
			elements := make([]types.String, len(set.Axes.Elements()))
			diags.Append(set.Axes.ElementsAs(ctx, &elements, false)...)
			axes := make([]string, len(elements))
			for j, axis := range elements {
				axes[j] = axis.ValueString()
			}
			s.Axes = axes
		}

		if !set.Attributes.IsNull() {
			elements := make([]types.String, len(set.Attributes.Elements()))
			diags.Append(set.Attributes.ElementsAs(ctx, &elements, false)...)
			attributes := make([]string, len(elements))
			for j, a := range elements {
				attributes[j] = a.ValueString()
			}
			s.Attributes = attributes
		}

		sets[i] = s
	}
	a.VariantAttributeSets = sets

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *FamilyVariantResource) mapToTfObject(respDiags *diag.Diagnostics, data *FamilyVariantResourceModel, apiData *goakeneo.FamilyVariant) {
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

	if len(apiData.VariantAttributeSets) > 0 {
		sets := make([]VariantAttributeSetModel, len(apiData.VariantAttributeSets))
		for i, set := range apiData.VariantAttributeSets {
			s := VariantAttributeSetModel{
				Level: types.Int64Value(int64(set.Level)),
			}

			if len(set.Axes) > 0 {
				axes := make([]attr.Value, len(set.Axes))
				for j, axis := range set.Axes {
					axes[j] = types.StringValue(axis)
				}
				listVal, diags := types.ListValue(types.StringType, axes)
				if diags.HasError() {
					respDiags.Append(diags...)
				}
				s.Axes = listVal
			}

			if len(set.Attributes) > 0 {
				attributes := make([]attr.Value, len(set.Attributes))
				for j, a := range set.Attributes {
					attributes[j] = types.StringValue(a)
				}
				listVal, diags := types.ListValue(types.StringType, attributes)
				if diags.HasError() {
					respDiags.Append(diags...)
				}
				s.Attributes = listVal
			}

			sets[i] = s
		}
		data.VariantAttributeSets = sets
	}
}
