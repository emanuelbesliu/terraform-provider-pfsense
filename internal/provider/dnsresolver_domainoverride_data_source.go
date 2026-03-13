package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DNSResolverDomainOverrideDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverDomainOverrideDataSource)(nil)
)

type DNSResolverDomainOverrideDataSourceModel struct {
	DNSResolverDomainOverrideModel
}

func NewDNSResolverDomainOverrideDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverDomainOverrideDataSource{}
}

type DNSResolverDomainOverrideDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverDomainOverrideDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_domainoverride", req.ProviderTypeName)
}

func (d *DNSResolverDomainOverrideDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single DNS resolver domain override by domain. Domains for which the resolver's standard DNS lookup process should be overridden and a different (non-standard) lookup server should be queried instead.",
		MarkdownDescription: "Retrieves a single DNS resolver [domain override](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver-domain-overrides.html) by domain. Domains for which the resolver's standard DNS lookup process should be overridden and a different (non-standard) lookup server should be queried instead.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: DNSResolverDomainOverrideModel{}.descriptions()["domain"].Description,
				Required:    true,
			},
			"ip_address": schema.StringAttribute{
				Description: DNSResolverDomainOverrideModel{}.descriptions()["ip_address"].Description,
				Computed:    true,
			},
			"tls_queries": schema.BoolAttribute{
				Description:         DNSResolverDomainOverrideModel{}.descriptions()["tls_queries"].Description,
				MarkdownDescription: DNSResolverDomainOverrideModel{}.descriptions()["tls_queries"].MarkdownDescription,
				Computed:            true,
			},
			"tls_hostname": schema.StringAttribute{
				Description: DNSResolverDomainOverrideModel{}.descriptions()["tls_hostname"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: DNSResolverDomainOverrideModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *DNSResolverDomainOverrideDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverDomainOverrideDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverDomainOverrideDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainOverride, err := d.client.GetDNSResolverDomainOverride(ctx, data.Domain.ValueString())
	if addError(&resp.Diagnostics, "Unable to get domain override", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *domainOverride)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
