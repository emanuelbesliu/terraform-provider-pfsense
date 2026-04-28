package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallVirtualIPResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallVirtualIPResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallVirtualIPResource)(nil)
)

type FirewallVirtualIPResourceModel struct {
	FirewallVirtualIPModel
}

func NewFirewallVirtualIPResource() resource.Resource { //nolint:ireturn
	return &FirewallVirtualIPResource{}
}

type FirewallVirtualIPResource struct {
	client *pfsense.Client
}

func (r *FirewallVirtualIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_virtual_ip", req.ProviderTypeName)
}

func (r *FirewallVirtualIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall Virtual IP. Virtual IPs allow binding additional IP addresses to interfaces for use with NAT, CARP high availability, or other services.",
		MarkdownDescription: "[Firewall Virtual IP](https://docs.netgate.com/pfsense/en/latest/firewall/virtual-ip-addresses.html). Virtual IPs allow binding additional IP addresses to interfaces for use with NAT, CARP high availability, or other services.",
		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["mode"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.ValidVirtualIPModes...),
				},
			},
			"interface": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"vhid": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["vhid"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinVHID), int64(pfsense.MaxVHID)),
				},
			},
			"advskew": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["advskew"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinAdvSkew), int64(pfsense.MaxAdvSkew)),
				},
			},
			"advbase": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["advbase"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinAdvBase), int64(pfsense.MaxAdvBase)),
				},
			},
			"password": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["password"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"subnet": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["subnet"].Description,
				Required:    true,
			},
			"subnet_bits": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["subnet_bits"].Description,
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 128),
				},
			},
			"description": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"unique_id": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["unique_id"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *FirewallVirtualIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallVirtualIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallVirtualIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var vipReq pfsense.VirtualIP
	resp.Diagnostics.Append(data.Value(ctx, &vipReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vip, err := r.client.CreateVirtualIP(ctx, vipReq)
	if addError(&resp.Diagnostics, "Error creating Virtual IP", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vip)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallVirtualIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallVirtualIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vip, err := r.client.GetVirtualIP(ctx, data.UniqueID.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading Virtual IP", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vip)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallVirtualIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallVirtualIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *FirewallVirtualIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var vipReq pfsense.VirtualIP
	resp.Diagnostics.Append(data.Value(ctx, &vipReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vipReq.UniqueID = state.UniqueID.ValueString()

	vip, err := r.client.UpdateVirtualIP(ctx, vipReq)
	if addError(&resp.Diagnostics, "Error updating Virtual IP", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vip)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallVirtualIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallVirtualIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVirtualIP(ctx, data.UniqueID.ValueString())
	if addError(&resp.Diagnostics, "Error deleting Virtual IP", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallVirtualIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("unique_id"), req, resp)
}
