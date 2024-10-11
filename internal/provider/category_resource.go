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
var _ resource.Resource = &CategoryResource{}
var _ resource.ResourceWithImportState = &CategoryResource{}
var _ resource.ResourceWithConfigure = &CategoryResource{}

func NewCategoryResource() resource.Resource {
	return &CategoryResource{}
}

// CategoryResource defines the resource implementation.
type CategoryResource struct {
	client *akeneox.CategoryService
}

// CategoryResourceModel describes the resource data model.
type CategoryResourceModel struct {
	Code   types.String `tfsdk:"code"`
	Parent types.String `tfsdk:"parent"`
	Labels types.Map    `tfsdk:"labels"`
}

func (r *CategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_category"
}

func (r *CategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo category resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Category code",
				Required:    true,
			},
			"parent": schema.StringAttribute{
				Description: "Category parent",
				Optional:    true,
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

func (r *CategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = akeneox.NewCategoryClient(data.Client)
}

func (r *CategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CategoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateCategory(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a category",
			"An unexpected error occurred when creating category. \n\n"+
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

func (r *CategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CategoryResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData, err := r.client.GetCategory(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a category",
			"An unexpected error occurred when reading category. \n\n"+
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

func (r *CategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CategoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateCategory(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating a category",
			"An unexpected error occurred when updating category. \n\n"+
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

func (r *CategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data CategoryResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for categories.",
	)
}

func (r *CategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *CategoryResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *CategoryResourceModel) *goakeneo.Category {
	a := goakeneo.Category{
		Code:   data.Code.ValueString(),
		Parent: data.Parent.ValueStringPointer(),
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

func (r *CategoryResource) mapToTfObject(respDiags *diag.Diagnostics, data *CategoryResourceModel, apiData *goakeneo.Category) {
	data.Code = types.StringValue(apiData.Code)
	data.Parent = types.StringPointerValue(apiData.Parent)

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
}
