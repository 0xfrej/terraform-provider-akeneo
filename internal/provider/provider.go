package provider

import (
	"context"
	"fmt"
	goakeneo "github.com/ezifyio/go-akeneo"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AkeneoProvider satisfies various provider interfaces.
var _ provider.Provider = &AkeneoProvider{}

// AkeneoProvider defines the provider implementation.
type AkeneoProvider struct {
	version string
}

// AkeneoProviderModel describes the provider data model.
type AkeneoProviderModel struct {
	Host                types.String `tfsdk:"host"`
	UnsecureApi         types.Bool   `tfsdk:"unsecure_api"`
	ApiUsername         types.String `tfsdk:"api_username"`
	ApiPassword         types.String `tfsdk:"api_password"`
	ApiClientId         types.String `tfsdk:"api_client_id"`
	ApiSecret           types.String `tfsdk:"api_client_secret"`
	ExtraAttributeTypes types.List   `tfsdk:"extra_attribute_types"`
}

type DataSourceData struct {
	Client *goakeneo.Client
}

type ResourceData struct {
	Client *goakeneo.Client
	//ExtraAttributeTypes *[]string
	//AvaialableLocales   []string
}

func (p *AkeneoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "akeneo"
	resp.Version = p.version
}

func (p *AkeneoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Akeneo host (optionally with port separated by double colon)",
				Required:            true,
				Sensitive:           true,
			},
			"unsecure_api": schema.BoolAttribute{
				MarkdownDescription: "Use http calls to the API",
				Optional:            true,
			},
			"api_username": schema.StringAttribute{
				MarkdownDescription: "Akeneo API client username",
				Required:            true,
				Sensitive:           true,
			},
			"api_password": schema.StringAttribute{
				MarkdownDescription: "Akeneo API client password",
				Required:            true,
				Sensitive:           true,
			},
			"api_client_id": schema.StringAttribute{
				MarkdownDescription: "Akeneo API client ID",
				Required:            true,
				Sensitive:           true,
			},
			"api_client_secret": schema.StringAttribute{
				MarkdownDescription: "Akeneo API client secret",
				Required:            true,
				Sensitive:           true,
			},
			"extra_attribute_types": schema.ListAttribute{
				MarkdownDescription: "Extra attribute types that are not supported by default",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
					),
				},
			},
		},
	}
}

func (p *AkeneoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AkeneoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connector := goakeneo.Connector{
		ClientID: data.ApiClientId.ValueString(),
		Secret:   data.ApiSecret.ValueString(),
		UserName: data.ApiUsername.ValueString(),
		Password: data.ApiPassword.ValueString(),
	}

	var proto string
	if data.UnsecureApi.ValueBool() {
		proto = "http"
	} else {
		proto = "https"
	}

	opts := []goakeneo.Option{
		goakeneo.WithBaseURL(fmt.Sprintf("%s://%s", proto, data.Host.ValueString())),
	}

	// TODO: implement rate limit
	// opts = append(opts, goakeneo.WithRateLimit(10, 1*time.Second))
	client, err := goakeneo.NewClient(connector, opts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"An unexpected error occurred when creating the Akeneo API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Akeneo API Client Error: "+err.Error(),
		)
		return
	}
	//TODO: add akeneo version validation - support 6, 7

	//var extraAttrs *[]string
	//if !data.ExtraAttributeTypes.IsNull() {
	//	elements := make([]types.String, 0, len(data.ExtraAttributeTypes.Elements()))
	//	resp.Diagnostics.Append(data.ExtraAttributeTypes.ElementsAs(ctx, &elements, false)...)
	//	res := make([]string, len(elements))
	//	for i, ext := range elements {
	//		res[i] = ext.ValueString()
	//	}
	//	extraAttrs = &res
	//}

	// todo: list all pages
	//locales, _, err := client.Locale.ListWithPagination(nil)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Unable to retrieve available locales",
	//		"An unexpected error occurred when retrieving the available locales from Akeneo. "+
	//			"If the error is not clear, please contact the provider developers.\n\n"+
	//			"Akeneo API Error: "+err.Error(),
	//	)
	//	return
	//}

	//availableLocales := make([]string, len(locales))
	//for i, locale := range locales {
	//	availableLocales[i] = locale.Code
	//}

	resp.DataSourceData = &DataSourceData{
		Client: client,
	}
	resp.ResourceData = &ResourceData{
		Client: client,
		//ExtraAttributeTypes: extraAttrs,
		//AvaialableLocales:   availableLocales,
	}
}

func (p *AkeneoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAttributeResource,
		NewAttributeOptionResource,
		NewAttributeGroupResource,
		NewFamilyResource,
		NewFamilyVariantResource,
	}
}

func (p *AkeneoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AkeneoProvider{
			version: version,
		}
	}
}
