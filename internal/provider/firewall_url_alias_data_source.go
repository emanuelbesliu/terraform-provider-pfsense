package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallURLAliasDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallURLAliasDataSource)(nil)
)

type FirewallURLAliasDataSourceModel struct {
	FirewallURLAliasModel
}

func NewFirewallURLAliasDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallURLAliasDataSource{}
}

type FirewallURLAliasDataSource struct {
	client *pfsense.Client
}

func (d *FirewallURLAliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_url_alias", req.ProviderTypeName)
}

func (d *FirewallURLAliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall URL alias by name. URL aliases import hosts, networks, or ports from URLs and can be periodically updated. Aliases can be referenced by firewall rules, port forwards, outbound NAT rules, and other places in the firewall.",
		MarkdownDescription: "Retrieves a single firewall [URL alias](https://docs.netgate.com/pfsense/en/latest/firewall/aliases.html) by name. URL aliases import hosts, networks, or ports from URLs and can be periodically updated. Aliases can be referenced by firewall rules, port forwards, outbound NAT rules, and other places in the firewall.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description:         FirewallURLAliasModel{}.descriptions()["type"].Description,
				MarkdownDescription: FirewallURLAliasModel{}.descriptions()["type"].MarkdownDescription,
				Computed:            true,
			},
			"update_frequency": schema.Int64Attribute{
				Description: FirewallURLAliasModel{}.descriptions()["update_frequency"].Description,
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["entries"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: FirewallURLAliasEntryModel{}.descriptions()["url"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallURLAliasEntryModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallURLAliasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallURLAliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallURLAliasDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlAlias, err := d.client.GetFirewallURLAlias(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get URL alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *urlAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
