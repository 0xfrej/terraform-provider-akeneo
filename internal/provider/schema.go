package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// immutableIdentifier marks a value Akeneo uses as an identifier.
// The API will not rename it, so Terraform has to replace the resource.
func immutableIdentifier() []planmodifier.String {
	return []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}
}
