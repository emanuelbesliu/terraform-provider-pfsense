package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemAdvancedNetworkingResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemAdvancedNetworkingResource)(nil)
	_ resource.ResourceWithImportState = (*SystemAdvancedNetworkingResource)(nil)
)

type SystemAdvancedNetworkingResourceModel struct {
	SystemAdvancedNetworkingModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemAdvancedNetworkingResource() resource.Resource { //nolint:ireturn
	return &SystemAdvancedNetworkingResource{}
}

type SystemAdvancedNetworkingResource struct {
	client *pfsense.Client
}

func (r *SystemAdvancedNetworkingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_networking", req.ProviderTypeName)
}

func (r *SystemAdvancedNetworkingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := SystemAdvancedNetworkingModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "System advanced networking configuration including DHCP options, IPv6 settings, and network interface offloading. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System advanced networking](https://docs.netgate.com/pfsense/en/latest/config/advanced/networking.html) configuration including DHCP options, IPv6 settings, and network interface offloading. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			// DHCP Options
			"dhcp_backend": schema.StringAttribute{
				Description:         descriptions["dhcp_backend"].Description,
				MarkdownDescription: descriptions["dhcp_backend"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedNetworkingDHCPBackend),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedNetworking{}.DHCPBackendOptions()...),
				},
			},
			"ignore_isc_warning": schema.BoolAttribute{
				Description: descriptions["ignore_isc_warning"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"radvd_debug": schema.BoolAttribute{
				Description: descriptions["radvd_debug"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dhcp6_debug": schema.BoolAttribute{
				Description: descriptions["dhcp6_debug"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dhcp6_no_release": schema.BoolAttribute{
				Description: descriptions["dhcp6_no_release"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"global_v6_duid": schema.StringAttribute{
				Description: descriptions["global_v6_duid"].Description,
				Optional:    true,
			},

			// IPv6 Options
			"ipv6_allow": schema.BoolAttribute{
				Description: descriptions["ipv6_allow"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ipv6_nat_enable": schema.BoolAttribute{
				Description: descriptions["ipv6_nat_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ipv6_nat_ip_address": schema.StringAttribute{
				Description: descriptions["ipv6_nat_ip_address"].Description,
				Optional:    true,
			},
			"prefer_ipv4": schema.BoolAttribute{
				Description: descriptions["prefer_ipv4"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ipv6_dont_create_local_dns": schema.BoolAttribute{
				Description: descriptions["ipv6_dont_create_local_dns"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Network Interfaces
			"disable_checksum_offloading": schema.BoolAttribute{
				Description: descriptions["disable_checksum_offloading"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_segmentation_offloading": schema.BoolAttribute{
				Description: descriptions["disable_segmentation_offloading"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_large_receive_offloading": schema.BoolAttribute{
				Description: descriptions["disable_large_receive_offloading"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"hn_altq_enable": schema.BoolAttribute{
				Description: descriptions["hn_altq_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"suppress_arp_messages": schema.BoolAttribute{
				Description: descriptions["suppress_arp_messages"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ip_change_kill_states": schema.BoolAttribute{
				Description: descriptions["ip_change_kill_states"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"use_if_pppoe": schema.BoolAttribute{
				Description: descriptions["use_if_pppoe"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Apply
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *SystemAdvancedNetworkingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemAdvancedNetworkingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemAdvancedNetworkingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedNetworking
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedNetworking(ctx, aReq)
	if addError(&resp.Diagnostics, "Error creating system advanced networking settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNetworkingChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced networking changes", err)
	}
}

func (r *SystemAdvancedNetworkingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemAdvancedNetworkingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetAdvancedNetworking(ctx)
	if addError(&resp.Diagnostics, "Error reading system advanced networking settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemAdvancedNetworkingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemAdvancedNetworkingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedNetworking
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedNetworking(ctx, aReq)
	if addError(&resp.Diagnostics, "Error updating system advanced networking settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNetworkingChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced networking changes", err)
	}
}

func (r *SystemAdvancedNetworkingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemAdvancedNetworkingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultNetworking := pfsense.AdvancedNetworking{
		DHCPBackend: pfsense.DefaultAdvancedNetworkingDHCPBackend,
	}

	_, err := r.client.UpdateAdvancedNetworking(ctx, defaultNetworking)
	if addError(&resp.Diagnostics, "Error resetting system advanced networking settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNetworkingChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced networking changes", err)
	}
}

func (r *SystemAdvancedNetworkingResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	a, err := r.client.GetAdvancedNetworking(ctx)
	if addError(&resp.Diagnostics, "Error importing system advanced networking settings", err) {
		return
	}

	var data SystemAdvancedNetworkingResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
