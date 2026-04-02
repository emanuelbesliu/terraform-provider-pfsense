package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DNSResolverGeneralDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverGeneralDataSource)(nil)
)

type DNSResolverGeneralDataSourceModel struct {
	DNSResolverGeneralModel
}

func NewDNSResolverGeneralDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverGeneralDataSource{}
}

type DNSResolverGeneralDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverGeneralDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_general", req.ProviderTypeName)
}

func (d *DNSResolverGeneralDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves the DNS resolver (Unbound) general settings configuration.",
		MarkdownDescription: "Retrieves the [DNS resolver (Unbound)](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html) general settings configuration.",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["enable"].Description,
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: DNSResolverGeneralModel{}.descriptions()["port"].Description,
				Computed:    true,
			},
			"enable_ssl": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["enable_ssl"].Description,
				Computed:    true,
			},
			"tls_port": schema.Int64Attribute{
				Description: DNSResolverGeneralModel{}.descriptions()["tls_port"].Description,
				Computed:    true,
			},
			"ssl_cert_ref": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["ssl_cert_ref"].Description,
				Computed:    true,
			},
			"active_interfaces": schema.ListAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["active_interfaces"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["active_interfaces"].MarkdownDescription,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"outgoing_interfaces": schema.ListAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["outgoing_interfaces"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["outgoing_interfaces"].MarkdownDescription,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"system_domain_local_zone_type": schema.StringAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["system_domain_local_zone_type"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["system_domain_local_zone_type"].MarkdownDescription,
				Computed:            true,
			},
			"dnssec": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["dnssec"].Description,
				Computed:    true,
			},
			"forwarding": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["forwarding"].Description,
				Computed:    true,
			},
			"forward_tls_upstream": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["forward_tls_upstream"].Description,
				Computed:    true,
			},
			"register_dhcp_leases": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_dhcp_leases"].Description,
				Computed:    true,
			},
			"register_dhcp_static_maps": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_dhcp_static_maps"].Description,
				Computed:    true,
			},
			"register_openvpn_clients": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_openvpn_clients"].Description,
				Computed:    true,
			},
			"strict_outgoing_interface": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["strict_outgoing_interface"].Description,
				Computed:    true,
			},
			"python": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python"].Description,
				Computed:    true,
			},
			"python_order": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python_order"].Description,
				Computed:    true,
			},
			"python_script": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python_script"].Description,
				Computed:    true,
			},
			"custom_options": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["custom_options"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *DNSResolverGeneralDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverGeneralDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverGeneralDataSourceModel

	dg, err := d.client.GetDNSResolverGeneral(ctx)
	if addError(&resp.Diagnostics, "Unable to get DNS resolver general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *dg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
