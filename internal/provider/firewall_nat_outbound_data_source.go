package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNATOutboundRuleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATOutboundRuleDataSource)(nil)
)

type FirewallNATOutboundRuleDataSourceModel struct {
	FirewallNATOutboundRuleModel
}

func NewFirewallNATOutboundRuleDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATOutboundRuleDataSource{}
}

type FirewallNATOutboundRuleDataSource struct {
	client *pfsense.Client
}

func (d *FirewallNATOutboundRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_outbound", req.ProviderTypeName)
}

func (d *FirewallNATOutboundRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single NAT outbound rule by its description.",
		MarkdownDescription: "Retrieves a single [NAT outbound](https://docs.netgate.com/pfsense/en/latest/nat/outbound.html) rule by its description.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["interface"].Description,
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["protocol"].Description,
				Computed:    true,
			},
			"source_address": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_address"].Description,
				Computed:    true,
			},
			"source_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_port"].Description,
				Computed:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_not"].Description,
				Computed:    true,
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_address"].Description,
				Computed:    true,
			},
			"destination_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_port"].Description,
				Computed:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_not"].Description,
				Computed:    true,
			},
			"target": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target"].Description,
				Computed:    true,
			},
			"target_ip": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target_ip"].Description,
				Computed:    true,
			},
			"target_ip_subnet": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target_ip_subnet"].Description,
				Computed:    true,
			},
			"nat_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["nat_port"].Description,
				Computed:    true,
			},
			"pool_options": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["pool_options"].Description,
				Computed:    true,
			},
			"source_hash_key": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_hash_key"].Description,
				Computed:    true,
			},
			"static_nat_port": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["static_nat_port"].Description,
				Computed:    true,
			},
			"no_sync": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["no_sync"].Description,
				Computed:    true,
			},
			"no_nat": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["no_nat"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["disabled"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["description"].Description,
				Required:    true,
			},
		},
	}
}

func (d *FirewallNATOutboundRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATOutboundRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATOutboundRuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetNATOutboundRule(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Unable to get NAT outbound rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
