package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DHCPv4ServerDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DHCPv4ServerDataSource)(nil)
)

type DHCPv4ServerDataSourceModel struct {
	DHCPv4ServerModel
}

func NewDHCPv4ServerDataSource() datasource.DataSource { //nolint:ireturn
	return &DHCPv4ServerDataSource{}
}

type DHCPv4ServerDataSource struct {
	client *pfsense.Client
}

func (d *DHCPv4ServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dhcpv4_server", req.ProviderTypeName)
}

func (d *DHCPv4ServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves the DHCPv4 server configuration for a single network interface.",
		MarkdownDescription: "Retrieves the DHCPv4 [server configuration](https://docs.netgate.com/pfsense/en/latest/services/dhcp/ipv4.html) for a single network interface.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["enable"].Description,
				Computed:    true,
			},
			"range_from": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["range_from"].Description,
				Computed:    true,
			},
			"range_to": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["range_to"].Description,
				Computed:    true,
			},
			"dns_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["dns_servers"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"gateway": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["gateway"].Description,
				Computed:    true,
			},
			"domain_name": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["domain_name"].Description,
				Computed:    true,
			},
			"domain_search_list": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["domain_search_list"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_lease_time": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["default_lease_time"].Description,
				Computed:    true,
				CustomType:  timetypes.GoDurationType{},
			},
			"maximum_lease_time": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["maximum_lease_time"].Description,
				Computed:    true,
				CustomType:  timetypes.GoDurationType{},
			},
			"wins_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["wins_servers"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"ntp_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ntp_servers"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tftp_server": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["tftp_server"].Description,
				Computed:    true,
			},
			"ldap_server": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ldap_server"].Description,
				Computed:    true,
			},
			"mac_allow": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["mac_allow"].Description,
				Computed:    true,
			},
			"mac_deny": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["mac_deny"].Description,
				Computed:    true,
			},
			"deny_unknown": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["deny_unknown"].Description,
				Computed:    true,
			},
			"ignore_client_uids": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ignore_client_uids"].Description,
				Computed:    true,
			},
			"static_arp": schema.BoolAttribute{
				Description:         DHCPv4ServerModel{}.descriptions()["static_arp"].Description,
				MarkdownDescription: DHCPv4ServerModel{}.descriptions()["static_arp"].MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (d *DHCPv4ServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DHCPv4ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DHCPv4ServerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	srv, err := d.client.GetDHCPv4Server(ctx, data.Interface.ValueString())
	if addError(&resp.Diagnostics, "Unable to get DHCPv4 server configuration", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *srv)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
