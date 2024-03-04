package stringvalidatorx

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var (
	isLocaleRegexp = regexp.MustCompile(`^[a-z]{2}_[A-Z]{2}$`)
)

type isLocaleCode struct {
}

func (i isLocaleCode) Description(_ context.Context) string {
	return "value must be valid locale code (example 'en_US')"
}

func (i isLocaleCode) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i isLocaleCode) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	val := request.ConfigValue.ValueString()
	if !isLocaleRegexp.MatchString(val) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			i.Description(ctx),
			val,
		))
	}
}

func IsLocaleCode() validator.String {
	return isLocaleCode{}
}
