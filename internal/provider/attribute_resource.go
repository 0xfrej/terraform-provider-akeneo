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

	_, err := r.client.CreateAttribute(*apiData)
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

	if !data.SortOrder.IsNull() {
		a.SortOrder = int(data.SortOrder.ValueInt64())
	}

	if !data.Localizable.IsNull() {
		a.Localizable = data.Localizable.ValueBool()
	}

	if !data.Scopable.IsNull() {
		a.Scopable = data.Scopable.ValueBool()
	}

	if !data.Unique.IsNull() {
		a.Unique = data.Unique.ValueBool()
	}

	if !data.UseableAsGridFilter.IsNull() {
		a.UseableAsGridFilter = data.UseableAsGridFilter.ValueBool()
	}

	if !data.MaxCharacters.IsNull() {
		a.MaxCharacters = int(data.MaxCharacters.ValueInt64())
	}

	if !data.ValidationRule.IsNull() {
		a.ValidationRule = data.ValidationRule.ValueString()
	}

	if !data.ValidationRegexp.IsNull() {
		a.ValidationRegexp = data.ValidationRegexp.ValueString()
	}

	if !data.WysiwygEnabled.IsNull() {
		a.WysiwygEnabled = data.WysiwygEnabled.ValueBool()
	}

	if !data.NumberMin.IsNull() {
		v, _ := data.NumberMin.ValueBigFloat().Float64()
		a.NumberMin = strconv.FormatFloat(v, 'f', -1, 64)
	}

	if !data.NumberMax.IsNull() {
		v, _ := data.NumberMin.ValueBigFloat().Float64()
		a.NumberMin = strconv.FormatFloat(v, 'f', -1, 64)
	}

	if !data.DecimalsAllowed.IsNull() {
		a.DecimalsAllowed = data.DecimalsAllowed.ValueBool()
	}

	if !data.NegativeAllowed.IsNull() {
		a.NegativeAllowed = data.NegativeAllowed.ValueBool()
	}

	if !data.MetricFamily.IsNull() {
		a.MetricFamily = data.MetricFamily.ValueString()
	}

	if !data.DefaultMetricUnit.IsNull() {
		a.DefaultMetricUnit = data.DefaultMetricUnit.ValueString()
	}

	if !data.DateMin.IsNull() {
		a.DateMin = data.DateMin.ValueString()
	}

	if !data.DateMax.IsNull() {
		a.DateMax = data.DateMax.ValueString()
	}

	if !data.MaxFileSize.IsNull() {
		a.MaxFileSize = strconv.Itoa(int(data.MaxFileSize.ValueInt64()))
	}

	if !data.ReferenceDataName.IsNull() {
		a.ReferenceDataName = data.ReferenceDataName.ValueString()
	}

	if !data.DefaultValue.IsNull() {
		a.DefaultValue = data.DefaultValue.ValueBool()
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

	if !data.GroupLabels.IsNull() {
		elements := make(map[string]types.String, len(data.GroupLabels.Elements()))
		diags.Append(data.GroupLabels.ElementsAs(ctx, &elements, false)...)
		labels := make(map[string]string)
		for locale, label := range elements {
			labels[locale] = label.ValueString()
		}
		a.GroupLabels = labels
	}

	if !data.AvailableLocales.IsNull() {
		elements := make([]types.String, 0, len(data.AvailableLocales.Elements()))
		diags.Append(data.AvailableLocales.ElementsAs(ctx, &elements, false)...)
		locales := make([]string, len(elements))
		for i, locale := range elements {
			locales[i] = locale.ValueString()
		}
		a.AvailableLocales = locales
	}

	if !data.AllowedExtensions.IsNull() {
		elements := make([]types.String, 0, len(data.AllowedExtensions.Elements()))
		diags.Append(data.AllowedExtensions.ElementsAs(ctx, &elements, false)...)
		exts := make([]string, len(elements))
		for i, ext := range elements {
			exts[i] = ext.ValueString()
		}
		a.AllowedExtensions = exts
	}

	if !data.TableConfiguration.IsNull() {
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
	// todo map only non empties
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
	data.SortOrder = types.Int64Value(int64(attrData.SortOrder))
	data.Localizable = types.BoolValue(attrData.Localizable)
	data.Scopable = types.BoolValue(attrData.Scopable)
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
	data.Unique = types.BoolValue(attrData.Unique)
	data.UseableAsGridFilter = types.BoolValue(attrData.UseableAsGridFilter)
	data.MaxCharacters = types.Int64Value(int64(attrData.MaxCharacters))
	data.ValidationRule = types.StringValue(attrData.ValidationRule)
	data.ValidationRegexp = types.StringValue(attrData.ValidationRegexp)
	data.WysiwygEnabled = types.BoolValue(attrData.WysiwygEnabled)
	if attrData.NumberMin != "" {
		v, err := strconv.ParseFloat(attrData.NumberMin, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
			return
		}
		data.NumberMin = types.NumberValue(big.NewFloat(v))
	}
	if attrData.NumberMax != "" {
		v, err := strconv.ParseFloat(attrData.NumberMax, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
			return
		}
		data.NumberMax = types.NumberValue(big.NewFloat(v))
	}
	data.DecimalsAllowed = types.BoolValue(attrData.DecimalsAllowed)
	data.NegativeAllowed = types.BoolValue(attrData.NegativeAllowed)
	data.MetricFamily = types.StringValue(attrData.MetricFamily)
	data.DefaultMetricUnit = types.StringValue(attrData.DefaultMetricUnit)
	data.DateMin = types.StringValue(attrData.DateMin)
	data.DateMax = types.StringValue(attrData.DateMax)
	if len(attrData.AvailableLocales) > 0 {
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
	if attrData.MaxFileSize != "" {
		v, err := strconv.ParseInt(attrData.MaxFileSize, 10, 64)
		if err != nil {
			respDiags.AddError("Error parsing float value", "Error parsing float value. \n\n"+"Error: "+err.Error())
		}
		data.MaxCharacters = types.Int64Value(v)
	}
	data.ReferenceDataName = types.StringValue(attrData.ReferenceDataName)
	data.DefaultValue = types.BoolValue(attrData.DefaultValue)
	if len(attrData.TableConfiguration) > 0 {
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
