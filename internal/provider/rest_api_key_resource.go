package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource              = (*RESTAPIKeyResource)(nil)
	_ resource.ResourceWithConfigure = (*RESTAPIKeyResource)(nil)
)

type RESTAPIKeyResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	HashAlgo    types.String `tfsdk:"hash_algo"`
	LengthBytes types.Int64  `tfsdk:"length_bytes"`
	Key         types.String `tfsdk:"key"`
	Hash        types.String `tfsdk:"hash"`
	Username    types.String `tfsdk:"username"`
}

func NewRESTAPIKeyResource() resource.Resource { //nolint:ireturn
	return &RESTAPIKeyResource{}
}

type RESTAPIKeyResource struct {
	client *pfsense.Client
}

func (r *RESTAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_rest_api_key", req.ProviderTypeName)
}

func (r *RESTAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates a pfSense REST API key via the REST API v2. The key value is only available at creation time and stored in Terraform state. Import is not supported.",
		MarkdownDescription: "Creates a pfSense REST API key via the REST API v2. The key value is only available at creation time and stored in Terraform state. Import is not supported.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"hash_algo": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("sha256"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("sha256", "sha384", "sha512"),
				},
			},
			"length_bytes": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(24),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(16, 24, 32, 64),
				},
			},
			"key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hash": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *RESTAPIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *RESTAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RESTAPIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	opts := pfsense.RESTAPIKeyCreateRequest{
		HashAlgo:    data.HashAlgo.ValueString(),
		LengthBytes: int(data.LengthBytes.ValueInt64()),
	}

	if !data.Description.IsNull() {
		opts.Description = data.Description.ValueString()
	}

	key, err := r.client.CreateRESTAPIKey(ctx, opts)
	if addError(&resp.Diagnostics, "Error creating REST API key", err) {
		return
	}

	data.ID = types.Int64Value(int64(key.ID))
	data.Key = types.StringValue(key.Key)
	data.Hash = types.StringValue(key.Hash)
	data.Username = types.StringValue(key.Username)
	data.HashAlgo = types.StringValue(key.HashAlgo)
	data.LengthBytes = types.Int64Value(int64(key.LengthBytes))

	if key.Description != "" {
		data.Description = types.StringValue(key.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RESTAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RESTAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetRESTAPIKeyByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		if errors.Is(err, pfsense.ErrNotFound) {
			resp.State.RemoveResource(ctx)

			return
		}

		addError(&resp.Diagnostics, "Error reading REST API key", err)

		return
	}

	data.Username = types.StringValue(key.Username)
	data.HashAlgo = types.StringValue(key.HashAlgo)
	data.LengthBytes = types.Int64Value(int64(key.LengthBytes))

	if key.Description != "" {
		data.Description = types.StringValue(key.Description)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RESTAPIKeyResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *RESTAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RESTAPIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRESTAPIKey(ctx, int(data.ID.ValueInt64()))
	if addError(&resp.Diagnostics, "Error deleting REST API key", err) {
		return
	}
}
