package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallPortAliasDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallPortAliasDataSource)(nil)
)

type FirewallPortAliasDataSourceModel struct {
	FirewallPortAliasModel
}

func NewFirewallPortAliasDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallPortAliasDataSource{}
}

type FirewallPortAliasDataSource struct {
	client *pfsense.Client
}

func (d *FirewallPortAliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_port_alias", req.ProviderTypeName)
}

func (d *FirewallPortAliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall port alias by name. Aliases act as placeholders for real hosts, networks, or ports and can be used to minimize the number of changes that have to be made if a host, network, or port changes.",
		MarkdownDescription: "Retrieves a single firewall [port alias](https://docs.netgate.com/pfsense/en/latest/firewall/aliases.html) by name. Aliases act as placeholders for real hosts, networks, or ports and can be used to minimize the number of changes that have to be made if a host, network, or port changes.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallPortAliasModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallPortAliasModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: FirewallPortAliasModel{}.descriptions()["entries"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.StringAttribute{
							Description: FirewallPortAliasEntryModel{}.descriptions()["port"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallPortAliasEntryModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallPortAliasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallPortAliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallPortAliasDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portAlias, err := d.client.GetFirewallPortAlias(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get port alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *portAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
