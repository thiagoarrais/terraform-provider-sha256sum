package sha256sum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &SHA256SumProvider{}

type SHA256SumProvider struct {
	Version string
}

// Metadata satisfies the provider.Provider interface for ExampleCloudProvider
func (p *SHA256SumProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sha256sum"
}

// Schema satisfies the provider.Provider interface for ExampleCloudProvider.
func (p *SHA256SumProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Provider specific implementation.
		},
	}
}

// Configure satisfies the provider.Provider interface for ExampleCloudProvider.
func (p *SHA256SumProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Provider specific implementation.
}

// DataSources satisfies the provider.Provider interface for ExampleCloudProvider.
func (p *SHA256SumProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Provider specific implementation
	}
}

// Resources satisfies the provider.Provider interface for ExampleCloudProvider.
func (p *SHA256SumProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &FileResource{} },
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SHA256SumProvider{
			Version: version,
		}
	}
}
