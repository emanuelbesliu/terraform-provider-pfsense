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
	_ datasource.DataSource              = (*DNSResolverHostOverrideDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverHostOverrideDataSource)(nil)
)

type DNSResolverHostOverrideDataSourceModel struct {
	DNSResolverHostOverrideModel
}

func NewDNSResolverHostOverrideDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverHostOverrideDataSource{}
}

type DNSResolverHostOverrideDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverHostOverrideDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_hostoverride", req.ProviderTypeName)
}

func (d *DNSResolverHostOverrideDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single DNS resolver host override by fully qualified domain name (FQDN). Hosts for which the resolver's standard DNS lookup process should be overridden and a specific IPv4 or IPv6 address should automatically be returned by the resolver.",
		MarkdownDescription: "Retrieves a single DNS resolver [host override](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver-host-overrides.html) by fully qualified domain name (FQDN). Hosts for which the resolver's standard DNS lookup process should be overridden and a specific IPv4 or IPv6 address should automatically be returned by the resolver.",
		Attributes: map[string]schema.Attribute{
			"fqdn": schema.StringAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["fqdn"].Description,
				Required:    true,
			},
			"host": schema.StringAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["host"].Description,
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["domain"].Description,
				Computed:    true,
			},
			"ip_addresses": schema.ListAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["ip_addresses"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"aliases": schema.ListNestedAttribute{
				Description: DNSResolverHostOverrideModel{}.descriptions()["aliases"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Description: DNSResolverHostOverrideAliasModel{}.descriptions()["host"].Description,
							Computed:    true,
						},
						"domain": schema.StringAttribute{
							Description: DNSResolverHostOverrideAliasModel{}.descriptions()["domain"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: DNSResolverHostOverrideAliasModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *DNSResolverHostOverrideDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverHostOverrideDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverHostOverrideDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostOverride, err := d.client.GetDNSResolverHostOverride(ctx, data.FQDN.ValueString())
	if addError(&resp.Diagnostics, "Unable to get host override", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *hostOverride)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
