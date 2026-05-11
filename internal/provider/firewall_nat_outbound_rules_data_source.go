package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNATOutboundRulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATOutboundRulesDataSource)(nil)
)

type FirewallNATOutboundRulesModel struct {
	Rules types.List   `tfsdk:"rules"`
	Mode  types.String `tfsdk:"mode"`
}

func NewFirewallNATOutboundRulesDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATOutboundRulesDataSource{}
}

type FirewallNATOutboundRulesDataSource struct {
	client *pfsense.Client
}

func (m *FirewallNATOutboundRulesModel) Set(ctx context.Context, rules pfsense.NATOutboundRules, mode string) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []FirewallNATOutboundRuleModel{}
	for _, r := range rules {
		var model FirewallNATOutboundRuleModel
		diags.Append(model.Set(ctx, r)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallNATOutboundRuleModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.Rules = listValue
	m.Mode = types.StringValue(mode)

	return diags
}

func (d *FirewallNATOutboundRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_outbound_rules", req.ProviderTypeName)
}

func (d *FirewallNATOutboundRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all NAT outbound rules and the outbound NAT mode.",
		MarkdownDescription: "Retrieves all [NAT outbound](https://docs.netgate.com/pfsense/en/latest/nat/outbound.html) rules and the outbound NAT mode.",
		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				Description: "Outbound NAT mode. Values: 'automatic', 'hybrid', 'advanced', 'disabled'.",
				Computed:    true,
			},
			"rules": schema.ListNestedAttribute{
				Description: "List of NAT outbound rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallNATOutboundRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATOutboundRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATOutboundRulesModel

	rules, err := d.client.GetNATOutboundRules(ctx)
	if addError(&resp.Diagnostics, "Unable to get NAT outbound rules", err) {
		return
	}

	mode, err := d.client.GetNATOutboundMode(ctx)
	if addError(&resp.Diagnostics, "Unable to get NAT outbound mode", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rules, mode)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
