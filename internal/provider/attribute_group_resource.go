package provider

import (
	"context"
	"fmt"
	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
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
var _ resource.Resource = &AttributeGroupResource{}
var _ resource.ResourceWithImportState = &AttributeGroupResource{}
var _ resource.ResourceWithConfigure = &AttributeGroupResource{}

func NewAttributeGroupResource() resource.Resource {
	return &AttributeGroupResource{}
}

// AttributeGroupResource defines the resource implementation.
type AttributeGroupResource struct {
	client *akeneox.AttributeService
}

// AttributeGroupResourceModel describes the resource data model.
type AttributeGroupResourceModel struct {
	Code      types.String `tfsdk:"code"`
	SortOrder types.Int64  `tfsdk:"sort_order"`
	Labels    types.Map    `tfsdk:"labels"`
}

func (r *AttributeGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute_group"
}

func (r *AttributeGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo attribute group resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Attribute group code",
				Required:    true,
			},
			"sort_order": schema.Int64Attribute{
				Description: "Order of the attribute group",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647),
				},
			},
			"labels": schema.MapAttribute{
				Description: "Label definition per locale",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidatorx.IsLocaleCode()),
				},
			},
		},
	}
}

func (r *AttributeGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = akeneox.NewAttributeClient(data.Client)
}

func (r *AttributeGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AttributeGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateAttributeGroup(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating an attribute group",
			"An unexpected error occurred when creating attribute group. \n\n"+
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

func (r *AttributeGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AttributeGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrData, err := r.client.GetAttributeGroup(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading an attribute group",
			"An unexpected error occurred when reading attribute group. \n\n"+
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

func (r *AttributeGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AttributeGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAttributeGroup(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating an attribute group",
			"An unexpected error occurred when updating attribute group. \n\n"+
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

func (r *AttributeGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data AttributeGroupResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for attribute group.",
	)
}

func (r *AttributeGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *AttributeGroupResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *AttributeGroupResourceModel) *akeneox.AttributeGroup {
	a := akeneox.AttributeGroup{
		Code: data.Code.ValueString(),
	}

	if !data.SortOrder.IsNull() {
		a.SortOrder = int(data.SortOrder.ValueInt64())
	}

	if !data.Labels.IsNull() {
		elements := make(map[string]types.String, len(data.Labels.Elements()))
		diags.Append(data.Labels.ElementsAs(ctx, &elements, false)...)
		labels := make(map[string]string)
		for locale, label := range elements {
			labels[locale] = label.ValueString()
		}
		a.Labels = labels
	}

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *AttributeGroupResource) mapToTfObject(respDiags *diag.Diagnostics, data *AttributeGroupResourceModel, attrData *akeneox.AttributeGroup) {
	data.Code = types.StringValue(attrData.Code)

	if len(attrData.Labels) > 0 {
		elements := make(map[string]attr.Value, len(attrData.Labels))

		for k, v := range attrData.Labels {
			elements[k] = types.StringValue(v)
		}

		mapVal, diags := types.MapValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Labels = mapVal
	}

	data.SortOrder = types.Int64Value(int64(attrData.SortOrder))
}
