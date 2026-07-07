package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// openVPNServerComputedAttributes returns the schema attributes for an OpenVPN
// server as computed values, used by both the singular and plural data sources.
// When vpnIDRequired is true the vpn_id attribute is marked required (used to
// look up a single server), otherwise it is computed.
func openVPNServerComputedAttributes(vpnIDRequired bool) map[string]schema.Attribute {
	descriptions := OpenVPNServerModel{}.descriptions()

	computedString := func(key string) schema.StringAttribute {
		return schema.StringAttribute{
			Description: descriptions[key].Description,
			Computed:    true,
		}
	}

	computedBool := func(key string) schema.BoolAttribute {
		return schema.BoolAttribute{
			Description: descriptions[key].Description,
			Computed:    true,
		}
	}

	computedStringList := func(key string) schema.ListAttribute {
		return schema.ListAttribute{
			Description: descriptions[key].Description,
			Computed:    true,
			ElementType: types.StringType,
		}
	}

	attributes := map[string]schema.Attribute{
		"vpn_id": schema.StringAttribute{
			Description: descriptions["vpn_id"].Description,
			Required:    vpnIDRequired,
			Computed:    !vpnIDRequired,
		},
		"disable":                 computedBool("disable"),
		"mode":                    computedString("mode"),
		"auth_mode":               computedStringList("auth_mode"),
		"dev_mode":                computedString("dev_mode"),
		"protocol":                computedString("protocol"),
		"interface":               computedString("interface"),
		"ip_address":              computedString("ip_address"),
		"local_port":              computedString("local_port"),
		"description":             computedString("description"),
		"custom_options":          computedString("custom_options"),
		"tls":                     computedString("tls"),
		"tls_type":                computedString("tls_type"),
		"tls_auth_keydir":         computedString("tls_auth_keydir"),
		"ca_ref":                  computedString("ca_ref"),
		"crl_ref":                 computedString("crl_ref"),
		"ocsp_check":              computedBool("ocsp_check"),
		"ocsp_url":                computedString("ocsp_url"),
		"cert_ref":                computedString("cert_ref"),
		"dh_length":               computedString("dh_length"),
		"ecdh_curve":              computedString("ecdh_curve"),
		"cert_depth":              computedString("cert_depth"),
		"strict_user_cn":          computedBool("strict_user_cn"),
		"remote_cert_tls":         computedBool("remote_cert_tls"),
		"shared_key":              schema.StringAttribute{Description: descriptions["shared_key"].Description, Computed: true, Sensitive: true},
		"data_ciphers":            computedStringList("data_ciphers"),
		"data_ciphers_fallback":   computedString("data_ciphers_fallback"),
		"digest":                  computedString("digest"),
		"tunnel_network":          computedString("tunnel_network"),
		"tunnel_network_v6":       computedString("tunnel_network_v6"),
		"local_network":           computedString("local_network"),
		"local_network_v6":        computedString("local_network_v6"),
		"remote_network":          computedString("remote_network"),
		"remote_network_v6":       computedString("remote_network_v6"),
		"gw_redir":                computedBool("gw_redir"),
		"gw_redir_v6":             computedBool("gw_redir_v6"),
		"topology":                computedString("topology"),
		"max_clients":             computedString("max_clients"),
		"connection_limit":        computedString("connection_limit"),
		"client_to_client":        computedBool("client_to_client"),
		"duplicate_cn":            computedBool("duplicate_cn"),
		"dynamic_ip":              computedBool("dynamic_ip"),
		"compression":             computedString("compression"),
		"compression_push":        computedBool("compression_push"),
		"allow_compression":       computedString("allow_compression"),
		"pass_tos":                computedBool("pass_tos"),
		"dns_domain_enable":       computedBool("dns_domain_enable"),
		"dns_domain":              computedString("dns_domain"),
		"dns_server_enable":       computedBool("dns_server_enable"),
		"dns_server1":             computedString("dns_server1"),
		"dns_server2":             computedString("dns_server2"),
		"dns_server3":             computedString("dns_server3"),
		"dns_server4":             computedString("dns_server4"),
		"ntp_server_enable":       computedBool("ntp_server_enable"),
		"ntp_server1":             computedString("ntp_server1"),
		"ntp_server2":             computedString("ntp_server2"),
		"push_register_dns":       computedBool("push_register_dns"),
		"push_block_outside_dns":  computedBool("push_block_outside_dns"),
		"username_as_common_name": computedBool("username_as_common_name"),
		"create_gw":               computedString("create_gw"),
		"verbosity_level":         computedString("verbosity_level"),
	}

	return attributes
}

var (
	_ datasource.DataSource              = (*OpenVPNServerDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*OpenVPNServerDataSource)(nil)
)

type OpenVPNServerDataSourceModel struct {
	OpenVPNServerModel
}

func NewOpenVPNServerDataSource() datasource.DataSource { //nolint:ireturn
	return &OpenVPNServerDataSource{}
}

type OpenVPNServerDataSource struct {
	client *pfsense.Client
}

func (d *OpenVPNServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_openvpn_server", req.ProviderTypeName)
}

func (d *OpenVPNServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single OpenVPN server instance by its vpn_id.",
		MarkdownDescription: "Retrieves a single OpenVPN [server](https://docs.netgate.com/pfsense/en/latest/vpn/openvpn/index.html) instance by its `vpn_id`.",
		Attributes:          openVPNServerComputedAttributes(true),
	}
}

func (d *OpenVPNServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *OpenVPNServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OpenVPNServerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := d.client.GetOpenVPNServer(ctx, data.VPNID.ValueString())
	if addError(&resp.Diagnostics, "Unable to get OpenVPN server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
