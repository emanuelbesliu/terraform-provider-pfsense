package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemAdvancedNetworkingDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemAdvancedNetworkingDataSource)(nil)
)

type SystemAdvancedNetworkingDataSourceModel struct {
	SystemAdvancedNetworkingModel
}

func NewSystemAdvancedNetworkingDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemAdvancedNetworkingDataSource{}
}

type SystemAdvancedNetworkingDataSource struct {
	client *pfsense.Client
}

func (d *SystemAdvancedNetworkingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_networking", req.ProviderTypeName)
}

func (d *SystemAdvancedNetworkingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := SystemAdvancedNetworkingModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Retrieves the system advanced networking configuration including DHCP options, IPv6 settings, and network interface offloading.",
		MarkdownDescription: "Retrieves the [system advanced networking](https://docs.netgate.com/pfsense/en/latest/config/advanced/networking.html) configuration including DHCP options, IPv6 settings, and network interface offloading.",
		Attributes: map[string]schema.Attribute{
			// DHCP Options
			"dhcp_backend": schema.StringAttribute{
				Description:         descriptions["dhcp_backend"].Description,
				MarkdownDescription: descriptions["dhcp_backend"].MarkdownDescription,
				Computed:            true,
			},
			"ignore_isc_warning": schema.BoolAttribute{
				Description: descriptions["ignore_isc_warning"].Description,
				Computed:    true,
			},
			"radvd_debug": schema.BoolAttribute{
				Description: descriptions["radvd_debug"].Description,
				Computed:    true,
			},
			"dhcp6_debug": schema.BoolAttribute{
				Description: descriptions["dhcp6_debug"].Description,
				Computed:    true,
			},
			"dhcp6_no_release": schema.BoolAttribute{
				Description: descriptions["dhcp6_no_release"].Description,
				Computed:    true,
			},
			"global_v6_duid": schema.StringAttribute{
				Description: descriptions["global_v6_duid"].Description,
				Computed:    true,
			},

			// IPv6 Options
			"ipv6_allow": schema.BoolAttribute{
				Description: descriptions["ipv6_allow"].Description,
				Computed:    true,
			},
			"ipv6_nat_enable": schema.BoolAttribute{
				Description: descriptions["ipv6_nat_enable"].Description,
				Computed:    true,
			},
			"ipv6_nat_ip_address": schema.StringAttribute{
				Description: descriptions["ipv6_nat_ip_address"].Description,
				Computed:    true,
			},
			"prefer_ipv4": schema.BoolAttribute{
				Description: descriptions["prefer_ipv4"].Description,
				Computed:    true,
			},
			"ipv6_dont_create_local_dns": schema.BoolAttribute{
				Description: descriptions["ipv6_dont_create_local_dns"].Description,
				Computed:    true,
			},

			// Network Interfaces
			"disable_checksum_offloading": schema.BoolAttribute{
				Description: descriptions["disable_checksum_offloading"].Description,
				Computed:    true,
			},
			"disable_segmentation_offloading": schema.BoolAttribute{
				Description: descriptions["disable_segmentation_offloading"].Description,
				Computed:    true,
			},
			"disable_large_receive_offloading": schema.BoolAttribute{
				Description: descriptions["disable_large_receive_offloading"].Description,
				Computed:    true,
			},
			"hn_altq_enable": schema.BoolAttribute{
				Description: descriptions["hn_altq_enable"].Description,
				Computed:    true,
			},
			"suppress_arp_messages": schema.BoolAttribute{
				Description: descriptions["suppress_arp_messages"].Description,
				Computed:    true,
			},
			"ip_change_kill_states": schema.BoolAttribute{
				Description: descriptions["ip_change_kill_states"].Description,
				Computed:    true,
			},
			"use_if_pppoe": schema.BoolAttribute{
				Description: descriptions["use_if_pppoe"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemAdvancedNetworkingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemAdvancedNetworkingDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemAdvancedNetworkingDataSourceModel

	a, err := d.client.GetAdvancedNetworking(ctx)
	if addError(&resp.Diagnostics, "Unable to get system advanced networking settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
