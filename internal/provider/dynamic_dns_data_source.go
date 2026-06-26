package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DynamicDNSDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DynamicDNSDataSource)(nil)
)

type DynamicDNSDataSourceModel struct {
	DynamicDNSModel
}

func NewDynamicDNSDataSource() datasource.DataSource { //nolint:ireturn
	return &DynamicDNSDataSource{}
}

type DynamicDNSDataSource struct {
	client *pfsense.Client
}

func (d *DynamicDNSDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dynamic_dns", req.ProviderTypeName)
}

func (d *DynamicDNSDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	desc := DynamicDNSModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves a dynamic DNS client entry by its index.",
		MarkdownDescription: "Retrieves a [dynamic DNS](https://docs.netgate.com/pfsense/en/latest/services/dyndns/index.html) client entry by its index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: desc["id"].Description,
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: desc["type"].Description,
				Computed:    true,
			},
			"interface": schema.StringAttribute{
				Description: desc["interface"].Description,
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: desc["host"].Description,
				Computed:    true,
			},
			"domain_name": schema.StringAttribute{
				Description: desc["domain_name"].Description,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: desc["username"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"password": schema.StringAttribute{
				Description: desc["password"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"mx": schema.StringAttribute{
				Description: desc["mx"].Description,
				Computed:    true,
			},
			"wildcard": schema.BoolAttribute{
				Description: desc["wildcard"].Description,
				Computed:    true,
			},
			"proxied": schema.BoolAttribute{
				Description: desc["proxied"].Description,
				Computed:    true,
			},
			"verbose_log": schema.BoolAttribute{
				Description: desc["verbose_log"].Description,
				Computed:    true,
			},
			"curl_ipresolve_v4": schema.BoolAttribute{
				Description: desc["curl_ipresolve_v4"].Description,
				Computed:    true,
			},
			"curl_ssl_verifypeer": schema.BoolAttribute{
				Description: desc["curl_ssl_verifypeer"].Description,
				Computed:    true,
			},
			"zone_id": schema.StringAttribute{
				Description: desc["zone_id"].Description,
				Computed:    true,
			},
			"ttl": schema.StringAttribute{
				Description: desc["ttl"].Description,
				Computed:    true,
			},
			"max_cache_age": schema.StringAttribute{
				Description: desc["max_cache_age"].Description,
				Computed:    true,
			},
			"update_url": schema.StringAttribute{
				Description: desc["update_url"].Description,
				Computed:    true,
			},
			"result_match": schema.StringAttribute{
				Description: desc["result_match"].Description,
				Computed:    true,
			},
			"request_interface": schema.StringAttribute{
				Description: desc["request_interface"].Description,
				Computed:    true,
			},
			"curl_proxy": schema.StringAttribute{
				Description: desc["curl_proxy"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: desc["description"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: desc["disabled"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *DynamicDNSDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DynamicDNSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DynamicDNSDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := d.client.GetDynamicDNS(ctx, int(data.ID.ValueInt64()))
	if addError(&resp.Diagnostics, "Unable to get dynamic DNS entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
