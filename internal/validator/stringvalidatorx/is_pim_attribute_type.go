package stringvalidatorx

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	pimTypes = []string{
		"pim_catalog_identifier",
		"pim_catalog_text",
		"pim_catalog_textarea",
		"pim_catalog_simpleselect",
		"pim_catalog_multiselect",
		"pim_catalog_boolean",
		"pim_catalog_date",
		"pim_catalog_number",
		"pim_catalog_metric",
		"pim_catalog_price_collection",
		"pim_catalog_image",
		"pim_catalog_file",
		"pim_catalog_asset_collection",
		"akeneo_reference_entity",
		"akeneo_reference_entity_collection",
		"pim_reference_data_simpleselect",
		"pim_reference_data_multiselect",
		"pim_catalog_table",
	}
)

type isPimAttributeTypeValidator struct {
	extraTypes []string
}

func (i isPimAttributeTypeValidator) getTypes() []string {
	return append(pimTypes, i.extraTypes...)
}

func (i isPimAttributeTypeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %s", i.getTypes())
}

func (i isPimAttributeTypeValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i isPimAttributeTypeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	val := request.ConfigValue.String()
	for _, t := range i.getTypes() {
		if val == t {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		i.Description(ctx),
		val,
	))
}

func IsPimAttributeType(extraTypes *[]string) validator.String {
	v := isPimAttributeTypeValidator{}
	if extraTypes != nil {
		v.extraTypes = *extraTypes
	}
	return v
}
