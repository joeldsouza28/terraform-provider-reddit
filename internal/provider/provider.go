package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &RedditProvider{}
var _ provider.ProviderWithFunctions = &RedditProvider{}
var _ provider.ProviderWithEphemeralResources = &RedditProvider{}

type RedditProvider struct {
	version string
}

type RedditProviderModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
}

func (p *RedditProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RedditProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := &redditClient{
		ClientID:     data.ClientID.ValueString(),
		ClientSecret: data.ClientSecret.ValueString(),
		Username:     data.Username.ValueString(),
		Password:     data.Password.ValueString(),
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *RedditProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPostDataSource,
	}
}

func (p *RedditProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *RedditProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *RedditProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPostResource,
		NewCommentResource,
	}
}

func (p *RedditProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "reddit"
	resp.Version = p.version
}

func (p *RedditProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "Reddit Bot Client Id",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Description: "Reddit Bot Client Secret",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Reddit Username",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Description: "Reddit Password",
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedditProvider{
			version: version,
		}
	}
}
