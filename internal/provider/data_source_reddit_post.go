package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type PostDataSource struct {
	client *redditClient
}

type PostDataSourceModel struct {
	PostID    types.String `tfsdk:"post_id"`
	Title     types.String `tfsdk:"title"`
	Text      types.String `tfsdk:"text"`
	Subreddit types.String `tfsdk:"subreddit"`
}

var _ datasource.DataSource = &PostDataSource{}

func NewPostDataSource() datasource.DataSource {
	return &PostDataSource{}
}

func (d *PostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post"
}

func (d *PostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*redditClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

}

func (d *PostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reddit data source",

		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				MarkdownDescription: "Title",
				Computed:            true,
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Text",
				Computed:            true,
			},
			"subreddit": schema.StringAttribute{
				MarkdownDescription: "Title",
				Computed:            true,
			},
			"post_id": schema.StringAttribute{
				MarkdownDescription: "Post Id",
				Required:            true,
			},
		},
	}
}

func (d *PostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PostDataSourceModel

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := d.client.GetToken()

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)

	}

	postID := data.PostID.String()
	post, err := FetchPostByID(token, postID)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting post %s", err)

	}
	data.Title = types.StringValue(post.Title)
	data.Text = types.StringValue(post.Text)
	data.Subreddit = types.StringValue(post.Subreddit)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
