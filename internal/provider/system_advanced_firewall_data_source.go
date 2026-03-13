package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemAdvancedFirewallDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemAdvancedFirewallDataSource)(nil)
)

type SystemAdvancedFirewallDataSourceModel struct {
	SystemAdvancedFirewallModel
}

func NewSystemAdvancedFirewallDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemAdvancedFirewallDataSource{}
}

type SystemAdvancedFirewallDataSource struct {
	client *pfsense.Client
}

func (d *SystemAdvancedFirewallDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_firewall", req.ProviderTypeName)
}

func (d *SystemAdvancedFirewallDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := SystemAdvancedFirewallModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Retrieves the system advanced firewall and NAT configuration including packet processing, VPN packet processing, advanced options, bogon networks, NAT reflection, and state timeouts.",
		MarkdownDescription: "Retrieves the [system advanced firewall and NAT](https://docs.netgate.com/pfsense/en/latest/config/advanced/firewall-nat.html) configuration including packet processing, VPN packet processing, advanced options, bogon networks, NAT reflection, and state timeouts.",
		Attributes: map[string]schema.Attribute{
			// Packet Processing
			"scrub_no_df": schema.BoolAttribute{
				Description: descriptions["scrub_no_df"].Description,
				Computed:    true,
			},
			"scrub_random_id": schema.BoolAttribute{
				Description: descriptions["scrub_random_id"].Description,
				Computed:    true,
			},
			"firewall_optimization": schema.StringAttribute{
				Description:         descriptions["firewall_optimization"].Description,
				MarkdownDescription: descriptions["firewall_optimization"].MarkdownDescription,
				Computed:            true,
			},
			"disable_scrub": schema.BoolAttribute{
				Description: descriptions["disable_scrub"].Description,
				Computed:    true,
			},
			"adaptive_start": schema.Int64Attribute{
				Description: descriptions["adaptive_start"].Description,
				Computed:    true,
			},
			"adaptive_end": schema.Int64Attribute{
				Description: descriptions["adaptive_end"].Description,
				Computed:    true,
			},
			"maximum_states": schema.Int64Attribute{
				Description: descriptions["maximum_states"].Description,
				Computed:    true,
			},
			"maximum_table_entries": schema.Int64Attribute{
				Description: descriptions["maximum_table_entries"].Description,
				Computed:    true,
			},
			"maximum_fragment_entries": schema.Int64Attribute{
				Description: descriptions["maximum_fragment_entries"].Description,
				Computed:    true,
			},

			// VPN Packet Processing
			"vpn_scrub_no_df": schema.BoolAttribute{
				Description: descriptions["vpn_scrub_no_df"].Description,
				Computed:    true,
			},
			"vpn_fragment_reassemble": schema.BoolAttribute{
				Description: descriptions["vpn_fragment_reassemble"].Description,
				Computed:    true,
			},
			"max_mss_enable": schema.BoolAttribute{
				Description: descriptions["max_mss_enable"].Description,
				Computed:    true,
			},
			"max_mss": schema.Int64Attribute{
				Description: descriptions["max_mss"].Description,
				Computed:    true,
			},

			// Advanced Options
			"disable_firewall": schema.BoolAttribute{
				Description: descriptions["disable_firewall"].Description,
				Computed:    true,
			},
			"bypass_static_routes": schema.BoolAttribute{
				Description: descriptions["bypass_static_routes"].Description,
				Computed:    true,
			},
			"disable_vpn_rules": schema.BoolAttribute{
				Description: descriptions["disable_vpn_rules"].Description,
				Computed:    true,
			},
			"disable_reply_to": schema.BoolAttribute{
				Description: descriptions["disable_reply_to"].Description,
				Computed:    true,
			},
			"disable_negate": schema.BoolAttribute{
				Description: descriptions["disable_negate"].Description,
				Computed:    true,
			},
			"no_apipa_block": schema.BoolAttribute{
				Description: descriptions["no_apipa_block"].Description,
				Computed:    true,
			},
			"aliases_resolve_interval": schema.Int64Attribute{
				Description: descriptions["aliases_resolve_interval"].Description,
				Computed:    true,
			},
			"check_aliases_url_cert": schema.BoolAttribute{
				Description: descriptions["check_aliases_url_cert"].Description,
				Computed:    true,
			},

			// Bogon Networks
			"bogons_update_interval": schema.StringAttribute{
				Description:         descriptions["bogons_update_interval"].Description,
				MarkdownDescription: descriptions["bogons_update_interval"].MarkdownDescription,
				Computed:            true,
			},

			// NAT
			"nat_reflection_mode": schema.StringAttribute{
				Description:         descriptions["nat_reflection_mode"].Description,
				MarkdownDescription: descriptions["nat_reflection_mode"].MarkdownDescription,
				Computed:            true,
			},
			"nat_reflection_timeout": schema.Int64Attribute{
				Description: descriptions["nat_reflection_timeout"].Description,
				Computed:    true,
			},
			"enable_binat_reflection": schema.BoolAttribute{
				Description: descriptions["enable_binat_reflection"].Description,
				Computed:    true,
			},
			"enable_nat_reflection_helper": schema.BoolAttribute{
				Description: descriptions["enable_nat_reflection_helper"].Description,
				Computed:    true,
			},
			"tftp_proxy_interfaces": schema.StringAttribute{
				Description: descriptions["tftp_proxy_interfaces"].Description,
				Computed:    true,
			},

			// State Timeouts
			"tcp_first_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_first_timeout"].Description,
				Computed:    true,
			},
			"tcp_opening_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_opening_timeout"].Description,
				Computed:    true,
			},
			"tcp_established_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_established_timeout"].Description,
				Computed:    true,
			},
			"tcp_closing_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_closing_timeout"].Description,
				Computed:    true,
			},
			"tcp_fin_wait_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_fin_wait_timeout"].Description,
				Computed:    true,
			},
			"tcp_closed_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_closed_timeout"].Description,
				Computed:    true,
			},
			"tcp_tsdiff_timeout": schema.Int64Attribute{
				Description: descriptions["tcp_tsdiff_timeout"].Description,
				Computed:    true,
			},
			"udp_first_timeout": schema.Int64Attribute{
				Description: descriptions["udp_first_timeout"].Description,
				Computed:    true,
			},
			"udp_single_timeout": schema.Int64Attribute{
				Description: descriptions["udp_single_timeout"].Description,
				Computed:    true,
			},
			"udp_multiple_timeout": schema.Int64Attribute{
				Description: descriptions["udp_multiple_timeout"].Description,
				Computed:    true,
			},
			"icmp_first_timeout": schema.Int64Attribute{
				Description: descriptions["icmp_first_timeout"].Description,
				Computed:    true,
			},
			"icmp_error_timeout": schema.Int64Attribute{
				Description: descriptions["icmp_error_timeout"].Description,
				Computed:    true,
			},
			"other_first_timeout": schema.Int64Attribute{
				Description: descriptions["other_first_timeout"].Description,
				Computed:    true,
			},
			"other_single_timeout": schema.Int64Attribute{
				Description: descriptions["other_single_timeout"].Description,
				Computed:    true,
			},
			"other_multiple_timeout": schema.Int64Attribute{
				Description: descriptions["other_multiple_timeout"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemAdvancedFirewallDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemAdvancedFirewallDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemAdvancedFirewallDataSourceModel

	a, err := d.client.GetAdvancedFirewall(ctx)
	if addError(&resp.Diagnostics, "Unable to get system advanced firewall settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
