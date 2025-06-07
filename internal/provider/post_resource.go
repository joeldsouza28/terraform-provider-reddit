package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &postResource{}

type postResource struct {
	client *redditClient
}

type postResourceModel struct {
	Title     types.String `tfsdk:"title"`
	Text      types.String `tfsdk:"text"`
	Subreddit types.String `tfsdk:"subreddit"`
	PostID    types.String `tfsdk:"post_id"`
	Flair     types.String `tfsdk:"flair"`
	NSFW      types.Bool   `tfsdk:"nsfw"`
}

func NewPostResource() resource.Resource {
	return &postResource{}
}

func (r *postResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post"
}

func (r *postResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*redditClient)
}

func (r *postResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				Required: true,
				Computed: false,
			},
			"text": schema.StringAttribute{
				Optional: true,
			},
			"subreddit": schema.StringAttribute{
				Required: true,
			},
			"post_id": schema.StringAttribute{
				Computed: true,
			},
			"flair": schema.StringAttribute{
				Optional: true,
			},
			"nsfw": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

func (r *postResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data postResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetToken()
	if err != nil {
		resp.Diagnostics.AddError("Token Error", err.Error())
		return
	}

	postID, err := SubmitPost(token, data.Subreddit.ValueString(), data.Title.ValueString(), data.Text.ValueString(), data.Flair.ValueString(), data.NSFW.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Post Creation Error", err.Error())
		return
	}

	data.PostID = types.StringValue(postID)
	resp.State.Set(ctx, data)
}

func (r *postResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Not implemented
}

func (r *postResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data postResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetToken()
	if err != nil {
		resp.Diagnostics.AddError("Token Error", err.Error())
		return
	}

	err = UpdatePostText(token, data.PostID.ValueString(), data.Text.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}
	data.PostID = types.StringValue(data.PostID.String())
	resp.State.Set(ctx, data)
}

func (r *postResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	postID := state.PostID.ValueString()
	if !strings.HasPrefix(postID, "t3_") {
		postID = "t3_" + postID
	}

	token, err := r.client.GetToken()
	if err != nil {
		resp.Diagnostics.AddError("Token Error", err.Error())
		return
	}

	err = DeletePost(token, postID)
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}
}
