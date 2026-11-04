package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource              = &DhcpBackendResource{}
	_ resource.ResourceWithConfigure = &DhcpBackendResource{}
)

type DhcpBackendResourceModel struct {
	Backend types.String `tfsdk:"backend"`
}

type DhcpBackendResource struct {
	client *pfsense.Client
}

func NewDhcpBackendResource() resource.Resource { //nolint:ireturn
	return &DhcpBackendResource{}
}

func (r *DhcpBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dhcp_backend", req.ProviderTypeName)
}

func (r *DhcpBackendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Global DHCP server backend selector ([`isc` or `kea`](https://docs.netgate.com/pfsense/en/latest/services/dhcp/index.html)), shown under System > Advanced, Networking (Server Backend). ISC DHCP is deprecated (end of life); Kea is the modern replacement. Both backends read the same DHCP server configuration managed by this provider. Restart DHCP services after changes if leases must be served immediately.",
		MarkdownDescription: "Global DHCP server backend selector ([`isc` or `kea`](https://docs.netgate.com/pfsense/en/latest/services/dhcp/index.html)), shown under **System > Advanced > Networking** (Server Backend). ISC DHCP is deprecated (end-of-life upstream); Kea is the modern replacement. Both backends read the same DHCP server configuration managed by this provider. Restart DHCP services after changes if leases must be served immediately.",
		Attributes: map[string]schema.Attribute{
			"backend": schema.StringAttribute{
				Description: "DHCP server backend: `isc` or `kea`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("isc", "kea"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DhcpBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *DhcpBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DhcpBackendResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	backend, err := r.client.CreateDhcpBackend(ctx, data.Backend.ValueString())
	if addError(&resp.Diagnostics, "Error creating DHCP backend", err) {
		return
	}

	data.Backend = types.StringValue(backend)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DhcpBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DhcpBackendResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	backend, err := r.client.GetDhcpBackend(ctx)
	if addError(&resp.Diagnostics, "Error reading DHCP backend", err) {
		return
	}

	data.Backend = types.StringValue(backend)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DhcpBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DhcpBackendResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	backend, err := r.client.UpdateDhcpBackend(ctx, data.Backend.ValueString())
	if addError(&resp.Diagnostics, "Error updating DHCP backend", err) {
		return
	}

	data.Backend = types.StringValue(backend)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DhcpBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DhcpBackendResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteDhcpBackend(ctx, data.Backend.ValueString()); err != nil {
		addError(&resp.Diagnostics, "Error deleting DHCP backend", err)
	}
}

var _ = path.Root
