package provider

import (
	"context"
	"fmt"
	"github.com/0xfrej/terraform-provider-akeneo/internal/akeneox"
	"github.com/0xfrej/terraform-provider-akeneo/internal/validator/stringvalidatorx"
	goakeneo "github.com/ezifyio/go-akeneo"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"math/big"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AttributeResource{}
var _ resource.ResourceWithImportState = &AttributeResource{}
var _ resource.ResourceWithConfigure = &AttributeResource{}

func NewAttributeResource() resource.Resource {
	return &AttributeResource{}
}

// AttributeResource defines the resource implementation.
type AttributeResource struct {
	client *akeneox.AttributeService
}

// AttributeResourceModel describes the resource data model.
type AttributeResourceModel struct {
	Code                types.String `tfsdk:"code"`
	Type                types.String `tfsdk:"type"`
	Labels              types.Map    `tfsdk:"labels"`
	Group               types.String `tfsdk:"group"`
	GroupLabels         types.Map    `tfsdk:"group_labels"`
	SortOrder           types.Int64  `tfsdk:"sort_order"`
	Localizable         types.Bool   `tfsdk:"localizable"`
	Scopable            types.Bool   `tfsdk:"scopable"`
	AvailableLocales    types.List   `tfsdk:"available_locales"`
	Unique              types.Bool   `tfsdk:"unique"`
	UseableAsGridFilter types.Bool   `tfsdk:"useable_as_grid_filter"`
	MaxCharacters       types.Int64  `tfsdk:"max_characters"`
	ValidationRule      types.String `tfsdk:"validation_rule"`
	ValidationRegexp    types.String `tfsdk:"validation_regexp"`
	WysiwygEnabled      types.Bool   `tfsdk:"wysiwyg_enabled"`
	NumberMin           types.Number `tfsdk:"number_min"`
	NumberMax           types.Number `tfsdk:"number_max"`
	DecimalsAllowed     types.Bool   `tfsdk:"decimals_allowed"`
	NegativeAllowed     types.Bool   `tfsdk:"negative_allowed"`
	MetricFamily        types.String `tfsdk:"metric_family"`
	DefaultMetricUnit   types.String `tfsdk:"default_metric_unit"`
	DateMin             types.String `tfsdk:"date_min"`
	DateMax             types.String `tfsdk:"date_max"`
	AllowedExtensions   types.List   `tfsdk:"allowed_extensions"`
	MaxFileSize         types.Int64  `tfsdk:"max_file_size"`
	ReferenceDataName   types.String `tfsdk:"reference_data_name"`
	DefaultValue        types.Bool   `tfsdk:"default_value"`
	TableConfiguration  types.List   `tfsdk:"table_configuration"`
}

func (r *AttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute"
}

func (r *AttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Akeneo attribute resource",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				Description: "Attribute code",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Attribute type - see akeneo available akeneo types in the documentation. Example: pim_catalog_file",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidatorx.IsPimAttributeType(nil),
						stringvalidator.LengthAtLeast(1),
					),
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
			"group": schema.StringAttribute{
				Description: "Attribute group",
				Required:    true,
			},
			"group_labels": schema.MapAttribute{
				Description: "Label definition per locale",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{
					mapvalidator.KeysAre(stringvalidatorx.IsLocaleCode()),
				},
			},
			"sort_order": schema.Int64Attribute{
				Description: "Order of the attribute in its group",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647),
				},
			},
			"localizable": schema.BoolAttribute{
				Description: "Whether the attribute is localizable, i.e. can have one value by locale",
				Optional:    true,
			},
			"scopable": schema.BoolAttribute{
				Description: "Whether the attribute is scopable, i.e. can have one value by channel",
				Optional:    true,
			},
			"available_locales": schema.ListAttribute{
				Description: "To make the attribute locale specific, specify here for which locales it is specific",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidatorx.IsLocaleCode()),
				},
			},
			"unique": schema.BoolAttribute{
				Description: " Whether two values for the attribute cannot be the same",
				Optional:    true,
			},
			"useable_as_grid_filter": schema.BoolAttribute{
				Description: "Whether the attribute can be used as a filter for the product grid in the PIM user interface",
				Optional:    true,
			},
			"max_characters": schema.Int64Attribute{
				Description: "Number maximum of characters allowed for the value of the attribute when the attribute type is `pim_catalog_text`, `pim_catalog_textarea` or `pim_catalog_identifier`",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2147483647),
				},
			},
			"validation_rule": schema.StringAttribute{
				Description: "Validation rule type used to validate any attribute value when the attribute type is `pim_catalog_text` or `pim_catalog_identifier`",
				Optional:    true,
			},
			"validation_regexp": schema.StringAttribute{
				Description: "Regexp expression used to validate any attribute value when the attribute type is `pim_catalog_text` or `pim_catalog_identifier`",
				Optional:    true,
			},
			"wysiwyg_enabled": schema.BoolAttribute{
				Description: "Whether the WYSIWYG interface is shown when the attribute type is `pim_catalog_textarea`",
				Optional:    true,
			},
			"number_min": schema.NumberAttribute{
				Description: "Minimum integer value allowed when the attribute type is `pim_catalog_metric`, `pim_catalog_price` or `pim_catalog_number`",
				Optional:    true,
			},
			"number_max": schema.NumberAttribute{
				Description: "Maximum integer value allowed when the attribute type is `pim_catalog_metric`, `pim_catalog_price` or `pim_catalog_number`",
				Optional:    true,
			},
			"decimals_allowed": schema.BoolAttribute{
				Description: "Whether decimals are allowed when the attribute type is `pim_catalog_metric`, `pim_catalog_price` or `pim_catalog_number`",
				Optional:    true,
			},
			"negative_allowed": schema.BoolAttribute{
				Description: "Whether negative values are allowed when the attribute type is `pim_catalog_metric` or `pim_catalog_number`",
				Optional:    true,
			},
			"metric_family": schema.StringAttribute{
				Description: "Metric family when the attribute type is `pim_catalog_metric`",
				Optional:    true,
			},
			"default_metric_unit": schema.StringAttribute{
				Description: "Default metric unit when the attribute type is `pim_catalog_metric`",
				Optional:    true,
			},
			"date_min": schema.StringAttribute{
				Description: "Minimum date allowed when the attribute type is `pim_catalog_date`",
				Optional:    true,
			},
			"date_max": schema.StringAttribute{
				Description: "Maximum date allowed when the attribute type is `pim_catalog_date`",
				Optional:    true,
			},
			"allowed_extensions": schema.ListAttribute{
				Description: "Extensions allowed when the attribute type is `pim_catalog_file` or `pim_catalog_image`",
				Optional:    true,
				ElementType: types.StringType,
			},
			"max_file_size": schema.Int64Attribute{
				Description: "Max file size in MB when the attribute type is `pim_catalog_file` or `pim_catalog_image`",
				Optional:    true,
			},
			"reference_data_name": schema.StringAttribute{
				Description: "Reference entity code when the attribute type is `akeneo_reference_entity` or `akeneo_reference_entity_collection` OR Asset family code when the attribute type is `pim_catalog_asset_collection`",
				Optional:    true,
			},
			"default_value": schema.BoolAttribute{
				Description: "Default value for a Yes/No attribute, applied when creating a new product or product model (only available since the 5.0)",
				Optional:    true,
			},
			"table_configuration": schema.ListAttribute{
				Description: "Configuration of the Table attribute (columns)",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *AttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AttributeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateAttribute(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating an attribute",
			"An unexpected error occurred when creating attribute. \n\n"+
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

func (r *AttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AttributeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attrData, err := r.client.GetAttribute(data.Code.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading an attribute",
			"An unexpected error occurred when reading attribute. \n\n"+
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

func (r *AttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AttributeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiData := r.mapToApiObject(ctx, &resp.Diagnostics, &data)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAttribute(*apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating an attribute",
			"An unexpected error occurred when updating attribute. \n\n"+
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

func (r *AttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data AttributeResourceModel
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

func (r *AttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

func (r *AttributeResource) mapToApiObject(ctx context.Context, diags *diag.Diagnostics, data *AttributeResourceModel) *goakeneo.Attribute {
	a := goakeneo.Attribute{
		Code:  data.Code.ValueString(),
		Type:  data.Type.ValueString(),
		Group: data.Group.ValueString(),
	}

	if !(data.SortOrder.IsNull() || data.SortOrder.IsUnknown()) {
		v := int(data.SortOrder.ValueInt64())
		a.SortOrder = &v
	}

	if !(data.Localizable.IsNull() || data.Localizable.IsUnknown()) {
		v := data.Localizable.ValueBool()
		a.Localizable = &v
	}

	if !(data.Scopable.IsNull() || data.Scopable.IsUnknown()) {
		v := data.Scopable.ValueBool()
		a.Scopable = &v
	}

	if !(data.Unique.IsNull() || data.Unique.IsUnknown()) {
		v := data.Unique.ValueBool()
		a.Unique = &v
	}

	if !(data.UseableAsGridFilter.IsNull() || data.UseableAsGridFilter.IsUnknown()) {
		v := data.UseableAsGridFilter.ValueBool()
		a.UseableAsGridFilter = &v
	}

	if !(data.MaxCharacters.IsNull() || data.MaxCharacters.IsUnknown()) {
		v := int(data.MaxCharacters.ValueInt64())
		a.MaxCharacters = &v
	}

	if !(data.ValidationRule.IsNull() || data.ValidationRule.IsUnknown()) {
		v := data.ValidationRule.ValueString()
		a.ValidationRule = &v
	}

	if !(data.ValidationRegexp.IsNull() || data.ValidationRegexp.IsUnknown()) {
		v := data.ValidationRegexp.ValueString()
		a.ValidationRegexp = &v
	}

	if !(data.WysiwygEnabled.IsNull() || data.WysiwygEnabled.IsUnknown()) {
		v := data.WysiwygEnabled.ValueBool()
		a.WysiwygEnabled = &v
	}

	if !(data.NumberMin.IsNull() || data.NumberMin.IsUnknown()) {
		v, _ := data.NumberMin.ValueBigFloat().Float64()
		r := strconv.FormatFloat(v, 'f', -1, 64)
		a.NumberMin = &r
	}

	if !(data.NumberMax.IsNull() || data.NumberMax.IsUnknown()) {
		v, _ := data.NumberMax.ValueBigFloat().Float64()
		r := strconv.FormatFloat(v, 'f', -1, 64)
		a.NumberMax = &r
	}

	if !(data.DecimalsAllowed.IsNull() || data.DecimalsAllowed.IsUnknown()) {
		v := data.DecimalsAllowed.ValueBool()
		a.DecimalsAllowed = &v
	}

	if !(data.NegativeAllowed.IsNull() || data.NegativeAllowed.IsUnknown()) {
		v := data.NegativeAllowed.ValueBool()
		a.NegativeAllowed = &v
	}

	if !(data.MetricFamily.IsNull() || data.MetricFamily.IsUnknown()) {
		v := data.MetricFamily.ValueString()
		a.MetricFamily = &v
	}

	if !(data.DefaultMetricUnit.IsNull() || data.DefaultMetricUnit.IsUnknown()) {
		v := data.DefaultMetricUnit.ValueString()
		a.DefaultMetricUnit = &v
	}

	if !(data.DateMin.IsNull() || data.DateMin.IsUnknown()) {
		v := data.DateMin.ValueString()
		a.DateMin = &v
	}

	if !(data.DateMax.IsNull() || data.DateMax.IsUnknown()) {
		v := data.DateMax.ValueString()
		a.DateMax = &v
	}

	if !(data.MaxFileSize.IsNull() || data.MaxFileSize.IsUnknown()) {
		v := strconv.Itoa(int(data.MaxFileSize.ValueInt64()))
		a.MaxFileSize = &v
	}

	if !(data.ReferenceDataName.IsNull() || data.ReferenceDataName.IsUnknown()) {
		v := data.ReferenceDataName.ValueString()
		a.ReferenceDataName = &v
	}

	if !(data.DefaultValue.IsNull() || data.DefaultValue.IsUnknown()) {
		v := data.DefaultValue.ValueBool()
		a.DefaultValue = &v
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

	if !(data.GroupLabels.IsNull() || data.GroupLabels.IsUnknown()) {
		elements := make(map[string]types.String, len(data.GroupLabels.Elements()))
		diags.Append(data.GroupLabels.ElementsAs(ctx, &elements, false)...)
		labels := make(map[string]string)
		for locale, label := range elements {
			labels[locale] = label.ValueString()
		}
		a.GroupLabels = labels
	}

	if !(data.AvailableLocales.IsNull() || data.AvailableLocales.IsUnknown()) {
		elements := make([]types.String, 0, len(data.AvailableLocales.Elements()))
		diags.Append(data.AvailableLocales.ElementsAs(ctx, &elements, false)...)
		locales := make([]string, len(elements))
		for i, locale := range elements {
			locales[i] = locale.ValueString()
		}
		a.AvailableLocales = locales
	}

	if !(data.AllowedExtensions.IsNull() || data.AllowedExtensions.IsUnknown()) {
		elements := make([]types.String, 0, len(data.AllowedExtensions.Elements()))
		diags.Append(data.AllowedExtensions.ElementsAs(ctx, &elements, false)...)
		exts := make([]string, len(elements))
		for i, ext := range elements {
			exts[i] = ext.ValueString()
		}
		a.AllowedExtensions = exts
	}

	if !(data.TableConfiguration.IsNull() || data.TableConfiguration.IsUnknown()) {
		elements := make([]types.String, 0, len(data.TableConfiguration.Elements()))
		diags.Append(data.TableConfiguration.ElementsAs(ctx, &elements, false)...)
		cfg := make([]string, len(elements))
		for i, val := range elements {
			cfg[i] = val.ValueString()
		}
		a.TableConfiguration = cfg
	}

	if diags.HasError() {
		return nil
	}

	return &a
}

func (r *AttributeResource) mapToTfObject(respDiags *diag.Diagnostics, data *AttributeResourceModel, attrData *goakeneo.Attribute) {
	data.Code = types.StringValue(attrData.Code)
	data.Type = types.StringValue(attrData.Type)
	data.Group = types.StringValue(attrData.Group)

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

	if len(attrData.GroupLabels) > 0 {
		elements := make(map[string]attr.Value, len(attrData.GroupLabels))

		for k, v := range attrData.GroupLabels {
			elements[k] = types.StringValue(v)
		}

		mapVal, diags := types.MapValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.GroupLabels = mapVal
	}
	if attrData.SortOrder != nil {
		data.SortOrder = types.Int64Value(int64(*attrData.SortOrder))
	}
	if attrData.Localizable != nil {
		data.Localizable = types.BoolValue(*attrData.Localizable)
	}
	if attrData.Scopable != nil {
		data.Scopable = types.BoolValue(*attrData.Scopable)
	}
	if len(attrData.AvailableLocales) > 0 {
		elements := make([]attr.Value, len(attrData.AvailableLocales))

		for k, v := range attrData.AvailableLocales {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.AvailableLocales = listVal
	}
	if attrData.Unique != nil {
		data.Unique = types.BoolValue(*attrData.Unique)
	}
	if attrData.UseableAsGridFilter != nil {
		data.UseableAsGridFilter = types.BoolValue(*attrData.UseableAsGridFilter)
	}
	if attrData.MaxCharacters != nil {
		data.MaxCharacters = types.Int64Value(int64(*attrData.MaxCharacters))
	}
	if attrData.ValidationRule != nil {
		data.ValidationRule = types.StringValue(*attrData.ValidationRule)
	}
	if attrData.ValidationRegexp != nil {
		data.ValidationRegexp = types.StringValue(*attrData.ValidationRegexp)
	}
	if attrData.WysiwygEnabled != nil {
		data.WysiwygEnabled = types.BoolValue(*attrData.WysiwygEnabled)
	}
	if attrData.NumberMin != nil {
		v, err := strconv.ParseFloat(*attrData.NumberMin, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
			return
		}
		data.NumberMin = types.NumberValue(big.NewFloat(v))
	}
	if attrData.NumberMax != nil {
		v, err := strconv.ParseFloat(*attrData.NumberMax, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
			return
		}
		data.NumberMax = types.NumberValue(big.NewFloat(v))
	}
	if attrData.DecimalsAllowed != nil {
		data.DecimalsAllowed = types.BoolValue(*attrData.DecimalsAllowed)
	}
	if attrData.NegativeAllowed != nil {
		data.NegativeAllowed = types.BoolValue(*attrData.NegativeAllowed)
	}
	if attrData.MetricFamily != nil {
		data.MetricFamily = types.StringValue(*attrData.MetricFamily)
	}
	if attrData.DefaultMetricUnit != nil {
		data.DefaultMetricUnit = types.StringValue(*attrData.DefaultMetricUnit)
	}
	if attrData.DateMin != nil {
		data.DateMin = types.StringValue(*attrData.DateMin)
	}
	if attrData.DateMax != nil {
		data.DateMax = types.StringValue(*attrData.DateMax)
	}
	if attrData.AvailableLocales != nil {
		elements := make([]attr.Value, len(attrData.AllowedExtensions))

		for k, v := range attrData.AllowedExtensions {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.AllowedExtensions = listVal
	}
	if attrData.MaxFileSize != nil {
		v, err := strconv.ParseInt(*attrData.MaxFileSize, 10, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
		}
		data.MaxCharacters = types.Int64Value(v)
	}
	if attrData.ReferenceDataName != nil {
		data.ReferenceDataName = types.StringValue(*attrData.ReferenceDataName)
	}
	if attrData.DefaultValue != nil {
		data.DefaultValue = types.BoolValue(*attrData.DefaultValue)
	}
	if attrData.TableConfiguration != nil {
		elements := make([]attr.Value, len(attrData.TableConfiguration))

		for k, v := range attrData.TableConfiguration {
			elements[k] = types.StringValue(v)
		}

		listVal, diags := types.ListValue(types.StringType, elements)
		if diags.HasError() {
			respDiags.Append(diags...)
		}
		data.TableConfiguration = listVal
	}
}
