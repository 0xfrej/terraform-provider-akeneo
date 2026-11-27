package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestResourceIdentifiersRequireReplace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := &AkeneoProvider{}

	extraIdentifiers := map[string][]string{
		"akeneo_attribute_option": {"attribute"},
		"akeneo_family_variant":   {"family_code"},
	}

	for _, factory := range p.Resources(ctx) {
		res := factory()

		meta := resource.MetadataResponse{}
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "akeneo"}, &meta)

		schemaResp := resource.SchemaResponse{}
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		if schemaResp.Diagnostics.HasError() {
			t.Fatalf("%s schema: %v", meta.TypeName, schemaResp.Diagnostics)
		}

		names := append([]string{"code"}, extraIdentifiers[meta.TypeName]...)
		for _, name := range names {
			attr, ok := schemaResp.Schema.Attributes[name].(schema.StringAttribute)
			if !ok {
				t.Errorf("%s.%s is missing or is not a string attribute", meta.TypeName, name)
				continue
			}
			if !requiresReplace(ctx, attr) {
				t.Errorf("%s.%s must force a new resource when it changes", meta.TypeName, name)
			}
		}
	}
}

func requiresReplace(ctx context.Context, attr schema.StringAttribute) bool {
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{},
		},
	}
	raw := func(value string) tftypes.Value {
		tfValue, err := types.StringValue(value).ToTerraformValue(ctx)
		if err != nil {
			panic(err)
		}
		return tftypes.NewValue(
			testSchema.Type().TerraformType(ctx),
			map[string]tftypes.Value{"code": tfValue},
		)
	}

	resp := planmodifier.StringResponse{
		PlanValue: types.StringValue("new"),
	}
	req := planmodifier.StringRequest{
		Path:        path.Root("code"),
		ConfigValue: types.StringValue("new"),
		PlanValue:   types.StringValue("new"),
		StateValue:  types.StringValue("old"),
		Plan: tfsdk.Plan{
			Schema: testSchema,
			Raw:    raw("new"),
		},
		State: tfsdk.State{
			Schema: testSchema,
			Raw:    raw("old"),
		},
	}
	for _, modifier := range attr.PlanModifiers {
		modifier.PlanModifyString(ctx, req, &resp)
	}
	return resp.RequiresReplace
}
