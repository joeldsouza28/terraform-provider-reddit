package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &commentResource{}

type commentResource struct {
	client *redditClient
}

type commentResourceModel struct {
	PostID    types.String `tfsdk:"post_id"`
	Comment   types.String `tfsdk:"comment"`
	CommentID types.String `tfsdk:"comment_id"`
}

func NewCommentResource() resource.Resource {
	return &commentResource{}
}

func (r *commentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_comment"
}

func (r *commentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*redditClient)
}

func (r *commentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"post_id": schema.StringAttribute{
				Required: true,
			},
			"comment": schema.StringAttribute{
				Optional: true,
			},
			"comment_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *commentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data commentResourceModel
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

	fullname := data.PostID.ValueString()
	if !strings.HasPrefix(fullname, "t3_") {
		fullname = "t3_" + fullname
	}

	commentID, err := AddComment(token, fullname, data.Comment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Add Comment Error", err.Error())
		return
	}

	data.CommentID = types.StringValue(commentID)
	resp.State.Set(ctx, &data)
}

func (r *commentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Not implemented (you can implement a fetch API if needed)
}

func (r *commentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data commentResourceModel
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

	commentID := data.CommentID.ValueString()
	if err := UpdatePostText(token, commentID, data.Comment.ValueString()); err != nil {
		resp.Diagnostics.AddError("Update Comment Error", err.Error())
		return
	}

	data.CommentID = types.StringValue(commentID)
	resp.State.Set(ctx, &data)
}

func (r *commentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state commentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.client.GetToken()
	if err != nil {
		resp.Diagnostics.AddError("Token Error", err.Error())
		return
	}

	commentID := state.CommentID.ValueString()
	if !strings.HasPrefix(commentID, "t1_") {
		commentID = "t1_" + commentID
	}

	if err := DeletePost(token, commentID); err != nil {
		resp.Diagnostics.AddError("Delete Comment Error", err.Error())
		return
	}
}
