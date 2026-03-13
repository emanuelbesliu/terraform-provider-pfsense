package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemAdvancedFirewallResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemAdvancedFirewallResource)(nil)
	_ resource.ResourceWithImportState = (*SystemAdvancedFirewallResource)(nil)
)

type SystemAdvancedFirewallResourceModel struct {
	SystemAdvancedFirewallModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemAdvancedFirewallResource() resource.Resource { //nolint:ireturn
	return &SystemAdvancedFirewallResource{}
}

type SystemAdvancedFirewallResource struct {
	client *pfsense.Client
}

func (r *SystemAdvancedFirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_firewall", req.ProviderTypeName)
}

func (r *SystemAdvancedFirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := SystemAdvancedFirewallModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "System advanced firewall and NAT configuration including packet processing, VPN packet processing, advanced options, bogon networks, NAT reflection, and state timeouts. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System advanced firewall and NAT](https://docs.netgate.com/pfsense/en/latest/config/advanced/firewall-nat.html) configuration including packet processing, VPN packet processing, advanced options, bogon networks, NAT reflection, and state timeouts. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			// Packet Processing
			"scrub_no_df": schema.BoolAttribute{
				Description: descriptions["scrub_no_df"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"scrub_random_id": schema.BoolAttribute{
				Description: descriptions["scrub_random_id"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"firewall_optimization": schema.StringAttribute{
				Description:         descriptions["firewall_optimization"].Description,
				MarkdownDescription: descriptions["firewall_optimization"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedFirewallOptimization),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedFirewall{}.OptimizationOptions()...),
				},
			},
			"disable_scrub": schema.BoolAttribute{
				Description: descriptions["disable_scrub"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"adaptive_start": schema.Int64Attribute{
				Description: descriptions["adaptive_start"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"adaptive_end": schema.Int64Attribute{
				Description: descriptions["adaptive_end"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"maximum_states": schema.Int64Attribute{
				Description: descriptions["maximum_states"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"maximum_table_entries": schema.Int64Attribute{
				Description: descriptions["maximum_table_entries"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"maximum_fragment_entries": schema.Int64Attribute{
				Description: descriptions["maximum_fragment_entries"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},

			// VPN Packet Processing
			"vpn_scrub_no_df": schema.BoolAttribute{
				Description: descriptions["vpn_scrub_no_df"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"vpn_fragment_reassemble": schema.BoolAttribute{
				Description: descriptions["vpn_fragment_reassemble"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"max_mss_enable": schema.BoolAttribute{
				Description: descriptions["max_mss_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"max_mss": schema.Int64Attribute{
				Description: descriptions["max_mss"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedFirewallMaxMSS)),
				Validators: []validator.Int64{
					int64validator.Between(576, 65535),
				},
			},

			// Advanced Options
			"disable_firewall": schema.BoolAttribute{
				Description: descriptions["disable_firewall"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"bypass_static_routes": schema.BoolAttribute{
				Description: descriptions["bypass_static_routes"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_vpn_rules": schema.BoolAttribute{
				Description: descriptions["disable_vpn_rules"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_reply_to": schema.BoolAttribute{
				Description: descriptions["disable_reply_to"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_negate": schema.BoolAttribute{
				Description: descriptions["disable_negate"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_apipa_block": schema.BoolAttribute{
				Description: descriptions["no_apipa_block"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"aliases_resolve_interval": schema.Int64Attribute{
				Description: descriptions["aliases_resolve_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"check_aliases_url_cert": schema.BoolAttribute{
				Description: descriptions["check_aliases_url_cert"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Bogon Networks
			"bogons_update_interval": schema.StringAttribute{
				Description:         descriptions["bogons_update_interval"].Description,
				MarkdownDescription: descriptions["bogons_update_interval"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedFirewallBogonsInterval),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedFirewall{}.BogonsIntervalOptions()...),
				},
			},

			// NAT
			"nat_reflection_mode": schema.StringAttribute{
				Description:         descriptions["nat_reflection_mode"].Description,
				MarkdownDescription: descriptions["nat_reflection_mode"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedFirewallNATReflection),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedFirewall{}.NATReflectionOptions()...),
				},
			},
			"nat_reflection_timeout": schema.Int64Attribute{
				Description: descriptions["nat_reflection_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"enable_binat_reflection": schema.BoolAttribute{
				Description: descriptions["enable_binat_reflection"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_nat_reflection_helper": schema.BoolAttribute{
				Description: descriptions["enable_nat_reflection_helper"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tftp_proxy_interfaces": schema.StringAttribute{
				Description: descriptions["tftp_proxy_interfaces"].Description,
				Optional:    true,
			},

			// State Timeouts
			"tcp_first_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_first_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_opening_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_opening_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_established_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_established_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_closing_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_closing_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_fin_wait_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_fin_wait_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_closed_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_closed_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"tcp_tsdiff_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_tsdiff_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"udp_first_timeout": schema.Int64Attribute{
				Description: descriptions["udp_first_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"udp_single_timeout": schema.Int64Attribute{
				Description: descriptions["udp_single_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"udp_multiple_timeout": schema.Int64Attribute{
				Description: descriptions["udp_multiple_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"icmp_first_timeout": schema.Int64Attribute{
				Description: descriptions["icmp_first_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"icmp_error_timeout": schema.Int64Attribute{
				Description: descriptions["icmp_error_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"other_first_timeout": schema.Int64Attribute{
				Description: descriptions["other_first_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"other_single_timeout": schema.Int64Attribute{
				Description: descriptions["other_single_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
			},
			"other_multiple_timeout": schema.Int64Attribute{
				Description: descriptions["other_multiple_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators:  []validator.Int64{int64validator.AtLeast(0)},
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

func (r *SystemAdvancedFirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemAdvancedFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemAdvancedFirewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedFirewall
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedFirewall(ctx, aReq)
	if addError(&resp.Diagnostics, "Error creating system advanced firewall settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedFirewallChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced firewall changes", err)
	}
}

func (r *SystemAdvancedFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemAdvancedFirewallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetAdvancedFirewall(ctx)
	if addError(&resp.Diagnostics, "Error reading system advanced firewall settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemAdvancedFirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemAdvancedFirewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedFirewall
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedFirewall(ctx, aReq)
	if addError(&resp.Diagnostics, "Error updating system advanced firewall settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedFirewallChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced firewall changes", err)
	}
}

func (r *SystemAdvancedFirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemAdvancedFirewallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultFirewall := pfsense.AdvancedFirewall{
		Optimization:   pfsense.DefaultAdvancedFirewallOptimization,
		MaxMSS:         pfsense.DefaultAdvancedFirewallMaxMSS,
		BogonsInterval: pfsense.DefaultAdvancedFirewallBogonsInterval,
		NATReflection:  pfsense.DefaultAdvancedFirewallNATReflection,
	}

	_, err := r.client.UpdateAdvancedFirewall(ctx, defaultFirewall)
	if addError(&resp.Diagnostics, "Error resetting system advanced firewall settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedFirewallChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced firewall changes", err)
	}
}

func (r *SystemAdvancedFirewallResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	a, err := r.client.GetAdvancedFirewall(ctx)
	if addError(&resp.Diagnostics, "Error importing system advanced firewall settings", err) {
		return
	}

	var data SystemAdvancedFirewallResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
