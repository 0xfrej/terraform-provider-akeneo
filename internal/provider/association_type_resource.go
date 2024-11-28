package provider

import (
	"context"
	"fmt"

	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
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
var _ resource.Resource = &AssociationTypeResource{}
var _ resource.ResourceWithImportState = &AssociationTypeResource{}
var _ resource.ResourceWithConfigure = &AssociationTypeResource{}

func NewAssociationTypeResource() resource.Resource {
	return &AssociationTypeResource{}
}

// AssociationTypeResource defines the resource implementation.
type AssociationTypeResource struct {
	client *akeneox.AssociationTypeService
}

// AssociationTypeResourceModel describes the resource data model.
type AssociationTypeResourceModel struct {
	Code         types.String `tfsdk:"code"`
	Labels       types.Map    `tfsdk:"labels"`
	IsQuantified types.Bool   `tfsdk:"is_quantified"`
	IsTwoWay     types.Bool   `tfsdk:"is_two_way"`
}

func (r *AssociationTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_association_type"
}

func (r *AssociationTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo association type resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Association type code",
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
			"is_quantified": schema.BoolAttribute{
				Description: "Whether the association type is a quantified association",
				Optional:    true,
			},
			"is_two_way": schema.BoolAttribute{
				Description: "Whether the association type is a two-way association",
				Optional:    true,
			},
		},
	}
}

func (r *AssociationTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = akeneox.NewAssociationTypeClient(data.Client)
}

func (r *AssociationTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AssociationTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateAssociationTypes(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a association type",
			"An unexpected error occurred when creating association type. \n\n"+
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

func (r *AssociationTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AssociationTypeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrData, err := r.client.GetAssociationType(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading an association type",
			"An unexpected error occurred when reading association type. \n\n"+
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

func (r *AssociationTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AssociationTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.UpdateAssociationTypes(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a association type",
			"An unexpected error occurred when creating association type. \n\n"+
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

func (r *AssociationTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data AssociationTypeResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for association types.",
	)
}

func (r *AssociationTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *AssociationTypeResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *AssociationTypeResourceModel) *akeneox.AssociationType {
	a := akeneox.AssociationType{
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

	if !(data.IsQuantified.IsNull() || data.IsQuantified.IsUnknown()) {
		v := data.IsQuantified.ValueBool()
		a.IsQuantified = &v
	}

	if !(data.IsTwoWay.IsNull() || data.IsTwoWay.IsUnknown()) {
		v := data.IsTwoWay.ValueBool()
		a.IsTwoWay = &v
	}

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *AssociationTypeResource) mapToTfObject(respDiags *diag.Diagnostics, data *AssociationTypeResourceModel, apiData *akeneox.AssociationType) {
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

	if apiData.IsQuantified != nil {
		data.IsQuantified = types.BoolValue(*apiData.IsQuantified)
	}

	if apiData.IsTwoWay != nil {
		data.IsQuantified = types.BoolValue(*apiData.IsTwoWay)
	}
}
