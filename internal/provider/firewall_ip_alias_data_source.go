package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallIPAliasDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallIPAliasDataSource)(nil)
)

type FirewallIPAliasDataSourceModel struct {
	FirewallIPAliasModel
}

func NewFirewallIPAliasDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallIPAliasDataSource{}
}

type FirewallIPAliasDataSource struct {
	client *pfsense.Client
}

func (d *FirewallIPAliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_ip_alias", req.ProviderTypeName)
}

func (d *FirewallIPAliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall IP alias by name. Aliases act as placeholders for real hosts, networks, or ports and can be used to minimize the number of changes that have to be made if a host, network, or port changes.",
		MarkdownDescription: "Retrieves a single firewall [IP alias](https://docs.netgate.com/pfsense/en/latest/firewall/aliases.html) by name. Aliases act as placeholders for real hosts, networks, or ports and can be used to minimize the number of changes that have to be made if a host, network, or port changes.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallIPAliasModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallIPAliasModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description:         FirewallIPAliasModel{}.descriptions()["type"].Description,
				MarkdownDescription: FirewallIPAliasModel{}.descriptions()["type"].MarkdownDescription,
				Computed:            true,
			},
			"entries": schema.ListNestedAttribute{
				Description: FirewallIPAliasModel{}.descriptions()["entries"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description: FirewallIPAliasEntryModel{}.descriptions()["ip"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallIPAliasEntryModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallIPAliasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallIPAliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallIPAliasDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ipAlias, err := d.client.GetFirewallIPAlias(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get IP alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ipAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
