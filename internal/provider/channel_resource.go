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
var _ resource.Resource = &ChannelResource{}
var _ resource.ResourceWithImportState = &ChannelResource{}
var _ resource.ResourceWithConfigure = &ChannelResource{}

func NewChannelResource() resource.Resource {
	return &ChannelResource{}
}

// ChannelResource defines the resource implementation.
type ChannelResource struct {
	client *akeneox.ChannelService
}

// ChannelResourceModel describes the resource data model.
type ChannelResourceModel struct {
	Code            types.String `tfsdk:"code"`
	Labels          types.Map    `tfsdk:"labels"`
	Locales         types.List   `tfsdk:"locales"`
	Currencies      types.List   `tfsdk:"currencies"`
	CategoryTree    types.String `tfsdk:"category_tree"`
	ConversionUnits types.Map    `tfsdk:"conversion_units"`
}

func (r *ChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *ChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo channel resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Channel code",
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
			"locales": schema.ListAttribute{
				Description: "Locales assigned to the channel",
				Required:    true,
				ElementType: types.StringType,
			},
			"currencies": schema.ListAttribute{
				Description: "Currencies assigned to the channel",
				Required:    true,
				ElementType: types.StringType,
			},
			"category_tree": schema.StringAttribute{
				Description: "Category tree assigned to the channel",
				Required:    true,
			},
			"conversion_units": schema.MapAttribute{
				Description: "Converion units assigned to the chennel",
				Optional:    true,
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
			},
		},
	}
}

func (r *ChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = akeneox.NewChannelClient(data.Client)
}

func (r *ChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ChannelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateChannel(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a channel",
			"An unexpected error occurred when creating channel. \n\n"+
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

func (r *ChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ChannelResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData, err := r.client.GetChannel(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating a channel",
			"An unexpected error occurred when reading channel. \n\n"+
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

func (r *ChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ChannelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateChannel(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating a channel",
			"An unexpected error occurred when updating channel. \n\n"+
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

func (r *ChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data ChannelResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.AddError(
		"This resource does not support deletes",
		"This resource does not support deletes. The Akeneo API does not support deletes for channels.",
	)
}

func (r *ChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *ChannelResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *ChannelResourceModel) *goakeneo.Channel {
	a := goakeneo.Channel{
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

	if !(data.Locales.IsNull() || data.Locales.IsUnknown()) {
		elements := make([]types.String, 0, len(data.Locales.Elements()))
		diags.Append(data.Locales.ElementsAs(ctx, &elements, false)...)
		locales := make([]string, len(elements))
		for i, ext := range elements {
			locales[i] = ext.ValueString()
		}
		a.Locales = locales
	}

	if !(data.Currencies.IsNull() || data.Currencies.IsUnknown()) {
		elements := make([]types.String, 0, len(data.Currencies.Elements()))
		diags.Append(data.Currencies.ElementsAs(ctx, &elements, false)...)
		currencies := make([]string, len(elements))
		for i, ext := range elements {
			currencies[i] = ext.ValueString()
		}
		a.Currencies = currencies
	}

	if !(data.CategoryTree.IsNull() || data.CategoryTree.IsUnknown()) {
		a.CategoryTree = data.CategoryTree.ValueString()
	}

	if !(data.ConversionUnits.IsNull() || data.ConversionUnits.IsUnknown()) {
		elements := make(map[string]types.String, len(data.ConversionUnits.Elements()))
		diags.Append(data.ConversionUnits.ElementsAs(ctx, &elements, false)...)
		units := make(map[string]string)
		for code, unit := range elements {
			units[code] = unit.ValueString()
		}
		a.ConversionUnits = units
	}

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *ChannelResource) mapToTfObject(respDiags *diag.Diagnostics, data *ChannelResourceModel, apiData *goakeneo.Channel) {
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

	if len(apiData.Locales) > 0 {
		elements := make([]attr.Value, len(apiData.Locales))

		for k, v := range apiData.Locales {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Locales = listVal
	}

	if len(apiData.Currencies) > 0 {
		elements := make([]attr.Value, len(apiData.Currencies))

		for k, v := range apiData.Currencies {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Currencies = listVal
	}

	data.CategoryTree = types.StringValue(apiData.CategoryTree)

	if len(apiData.ConversionUnits) > 0 {
		elements := make(map[string]attr.Value, len(apiData.ConversionUnits))

		for k, v := range apiData.ConversionUnits {
			elements[k] = types.StringValue(v)
		}

		mapVal, diags := types.MapValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.Labels = mapVal
	}
}
