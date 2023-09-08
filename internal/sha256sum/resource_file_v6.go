package sha256sum

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FileResource{}

type FileResource struct {
}

type FileModel struct {
	ID       types.String `tfsdk:"id"`
	Path     types.String `tfsdk:"path"`
	Contents types.String `tfsdk:"contents"`
}

func (*FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_v6"
}

func (*FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contents": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (*FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	contents, err := base64.StdEncoding.DecodeString(data.Contents.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to decode base64", err.Error())
		return
	}

	err = os.WriteFile(data.Path.ValueString(), contents, 0644)
	if err != nil {
		resp.Diagnostics.AddError("Unable to write file", err.Error())
		return
	}

	data.ID = data.Path
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (*FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	var data FileModel
	data.ID = id
	data.Path = id

	binaryContents, err := os.ReadFile(id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read file", err.Error())
		return
	}
	contents := base64.StdEncoding.EncodeToString(binaryContents)
	data.Contents = types.StringValue(contents)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (*FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	contents, err := base64.RawStdEncoding.DecodeString(data.Contents.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to decode base64", err.Error())
		return
	}

	err = os.WriteFile(data.Path.ValueString(), contents, 0644)
	if err != nil {
		resp.Diagnostics.AddError("Unable to write file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (*FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	err := os.Remove(id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to remove file", err.Error())
		return
	}
}
