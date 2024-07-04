package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &ZonefileProvider{}
var _ provider.ProviderWithFunctions = &ZonefileProvider{}

type ZonefileProvider struct {
	// version is set to the provider version on release, "dev" when the provider
	// is built and ran locally, and "test" when running acceptance testing.
	version string
}

type ZonefileProviderModel struct{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZonefileProvider{
			version: version,
		}
	}
}

func (p *ZonefileProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zonefile"
	resp.Version = p.version
}

func (p *ZonefileProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *ZonefileProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ZonefileProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *ZonefileProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRecordsDataSource,
		NewRecordSetsDataSource,
	}
}

func (p *ZonefileProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *ZonefileProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}
