package stringvalidatorx

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	pimConversionOpTypes = []string{
		"add",
		"sub",
		"mul",
		"div",
	}
)

type isConversionOperatorValidator struct {
}

func (i isConversionOperatorValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %s", pimConversionOpTypes)
}

func (i isConversionOperatorValidator) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i isConversionOperatorValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	val := request.ConfigValue.ValueString()
	for _, t := range pimConversionOpTypes {
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

func IsPimConversionOperator() validator.String {
	v := isConversionOperatorValidator{}
	return v
}
