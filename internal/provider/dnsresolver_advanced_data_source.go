package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DNSResolverAdvancedDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverAdvancedDataSource)(nil)
)

type DNSResolverAdvancedDataSourceModel struct {
	DNSResolverAdvancedModel
}

func NewDNSResolverAdvancedDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverAdvancedDataSource{}
}

type DNSResolverAdvancedDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverAdvancedDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_advanced", req.ProviderTypeName)
}

func (d *DNSResolverAdvancedDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves the DNS resolver (Unbound) advanced settings configuration.",
		MarkdownDescription: "Retrieves the [DNS resolver (Unbound)](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html) advanced settings configuration.",
		Attributes: map[string]schema.Attribute{
			"hide_identity": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["hide_identity"].Description,
				Computed:    true,
			},
			"hide_version": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["hide_version"].Description,
				Computed:    true,
			},
			"prefetch": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["prefetch"].Description,
				Computed:    true,
			},
			"prefetch_key": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["prefetch_key"].Description,
				Computed:    true,
			},
			"dnssec_stripped": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dnssec_stripped"].Description,
				Computed:    true,
			},
			"aggressive_nsec": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["aggressive_nsec"].Description,
				Computed:    true,
			},
			"qname_minimisation": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["qname_minimisation"].Description,
				Computed:    true,
			},
			"qname_minimisation_strict": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["qname_minimisation_strict"].Description,
				Computed:    true,
			},
			"use_caps": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["use_caps"].Description,
				Computed:    true,
			},
			"dns_record_cache": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns_record_cache"].Description,
				Computed:    true,
			},
			"disable_auto_access_control": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["disable_auto_access_control"].Description,
				Computed:    true,
			},
			"disable_auto_host_entries": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["disable_auto_host_entries"].Description,
				Computed:    true,
			},
			"dns64": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64"].Description,
				Computed:    true,
			},
			"dns64_prefix": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64_prefix"].Description,
				Computed:    true,
			},
			"dns64_netbits": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64_netbits"].Description,
				Computed:    true,
			},
			"msg_cache_size": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["msg_cache_size"].Description,
				Computed:    true,
			},
			"outgoing_num_tcp": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["outgoing_num_tcp"].Description,
				Computed:    true,
			},
			"incoming_num_tcp": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["incoming_num_tcp"].Description,
				Computed:    true,
			},
			"edns_buffer_size": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["edns_buffer_size"].Description,
				Computed:    true,
			},
			"num_queries_per_thread": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["num_queries_per_thread"].Description,
				Computed:    true,
			},
			"jostle_timeout": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["jostle_timeout"].Description,
				Computed:    true,
			},
			"cache_max_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["cache_max_ttl"].Description,
				Computed:    true,
			},
			"cache_min_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["cache_min_ttl"].Description,
				Computed:    true,
			},
			"infra_keep_probing": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_keep_probing"].Description,
				Computed:    true,
			},
			"infra_host_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_host_ttl"].Description,
				Computed:    true,
			},
			"infra_cache_num_hosts": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_cache_num_hosts"].Description,
				Computed:    true,
			},
			"unwanted_reply_threshold": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["unwanted_reply_threshold"].Description,
				Computed:    true,
			},
			"log_verbosity": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["log_verbosity"].Description,
				Computed:    true,
			},
			"sock_queue_timeout": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["sock_queue_timeout"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *DNSResolverAdvancedDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverAdvancedDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverAdvancedDataSourceModel

	da, err := d.client.GetDNSResolverAdvanced(ctx)
	if addError(&resp.Diagnostics, "Unable to get DNS resolver advanced settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *da)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
