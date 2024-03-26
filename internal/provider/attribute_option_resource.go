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
var _ resource.Resource = &AttributeOptionResource{}
var _ resource.ResourceWithImportState = &AttributeOptionResource{}
var _ resource.ResourceWithConfigure = &AttributeOptionResource{}

func NewAttributeOptionResource() resource.Resource {
	return &AttributeOptionResource{}
}

// AttributeOptionResource defines the resource implementation.
type AttributeOptionResource struct {
	client *akeneox.AttributeService
}

// AttributeOptionResourceModel describes the resource data model.
type AttributeOptionResourceModel struct {
	Code      types.String `tfsdk:"code"`
	Attribute types.String `tfsdk:"attribute"`
	SortOrder types.Int64  `tfsdk:"sort_order"`
	Labels    types.Map    `tfsdk:"labels"`
}

func (r *AttributeOptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute_option"
}

func (r *AttributeOptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo attribute option resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Attribute option code",
				Required:    true,
			},
			"attribute": schema.StringAttribute{
				Description: "Parent attribute code",
				Required:    true,
			},
			"sort_order": schema.Int64Attribute{
				Description: "Order of the attribute option",
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

func (r *AttributeOptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AttributeOptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AttributeOptionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateAttributeOption(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating an attribute option",
			"An unexpected error occurred when creating attribute option. \n\n"+
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

func (r *AttributeOptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AttributeOptionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrData, err := r.client.GetAttributeOption(data.Attribute.ValueString(), data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading an attribute option",
			"An unexpected error occurred when reading attribute option. \n\n"+
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

func (r *AttributeOptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AttributeOptionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAttributeOption(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating an attribute option",
			"An unexpected error occurred when updating attribute option. \n\n"+
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

func (r *AttributeOptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data AttributeOptionResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for attribute options.",
	)
}

func (r *AttributeOptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *AttributeOptionResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *AttributeOptionResourceModel) *goakeneo.AttributeOption {
	a := goakeneo.AttributeOption{
		Code:      data.Code.ValueString(),
		Attribute: data.Attribute.ValueString(),
	}

	if !(data.SortOrder.IsNull() || data.SortOrder.IsUnknown()) {
		v := int(data.SortOrder.ValueInt64())
		a.SortOrder = &v
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

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *AttributeOptionResource) mapToTfObject(respDiags *diag.Diagnostics, data *AttributeOptionResourceModel, attrData *goakeneo.AttributeOption) {
	data.Code = types.StringValue(attrData.Code)
	data.Attribute = types.StringValue(attrData.Attribute)

	if attrData.Labels != nil {
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

	if attrData.SortOrder != nil {
		data.SortOrder = types.Int64Value(int64(*attrData.SortOrder))
	}
}
