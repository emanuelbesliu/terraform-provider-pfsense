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
	_ datasource.DataSource              = (*FirewallRulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallRulesDataSource)(nil)
)

type FirewallRulesModel struct {
	FirewallRules types.List `tfsdk:"firewall_rules"`
}

func NewFirewallRulesDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallRulesDataSource{}
}

type FirewallRulesDataSource struct {
	client *pfsense.Client
}

func (m *FirewallRulesModel) Set(ctx context.Context, rules pfsense.FirewallRules) diag.Diagnostics {
	var diags diag.Diagnostics

	ruleModels := []FirewallRuleModel{}
	for _, r := range rules {
		var ruleModel FirewallRuleModel
		diags.Append(ruleModel.Set(ctx, r)...)
		ruleModels = append(ruleModels, ruleModel)
	}

	rulesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallRuleModel{}.AttrTypes()}, ruleModels)
	diags.Append(newDiags...)
	m.FirewallRules = rulesValue

	return diags
}

func (d *FirewallRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_rules", req.ProviderTypeName)
}

func (d *FirewallRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all firewall rules. Firewall rules control traffic flow through the pfSense firewall.",
		MarkdownDescription: "Retrieves all [firewall rules](https://docs.netgate.com/pfsense/en/latest/firewall/index.html). Firewall rules control traffic flow through the pfSense firewall.",
		Attributes: map[string]schema.Attribute{
			"firewall_rules": schema.ListNestedAttribute{
				Description: "List of firewall rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tracker": schema.StringAttribute{
							Description: FirewallRuleModel{}.descriptions()["tracker"].Description,
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *FirewallRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallRulesModel

	rules, err := d.client.GetFirewallRules(ctx)
	if addError(&resp.Diagnostics, "Unable to get firewall rules", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rules)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
