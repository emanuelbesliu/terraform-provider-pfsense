package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallRuleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallRuleDataSource)(nil)
)

type FirewallRuleDataSourceModel struct {
	FirewallRuleModel
}

func NewFirewallRuleDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallRuleDataSource{}
}

type FirewallRuleDataSource struct {
	client *pfsense.Client
}

func (d *FirewallRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_rule", req.ProviderTypeName)
}

func (d *FirewallRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall rule by tracker ID. Firewall rules control traffic flow through the pfSense firewall.",
		MarkdownDescription: "Retrieves a single [firewall rule](https://docs.netgate.com/pfsense/en/latest/firewall/index.html) by tracker ID. Firewall rules control traffic flow through the pfSense firewall.",
		Attributes: map[string]schema.Attribute{
			"tracker": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["tracker"].Description,
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["type"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["type"].MarkdownDescription,
				Computed:            true,
			},
			"interface": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["interface"].Description,
				Computed:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["protocol"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["protocol"].MarkdownDescription,
				Computed:            true,
			},
			"source_address": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_address"].Description,
				Computed:    true,
			},
			"source_port": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_port"].Description,
				Computed:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_not"].Description,
				Computed:    true,
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_address"].Description,
				Computed:    true,
			},
			"destination_port": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_port"].Description,
				Computed:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_not"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["disabled"].Description,
				Computed:    true,
			},
			"log": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["log"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *FirewallRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallRuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetFirewallRule(ctx, data.Tracker.ValueString())
	if addError(&resp.Diagnostics, "Unable to get firewall rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
