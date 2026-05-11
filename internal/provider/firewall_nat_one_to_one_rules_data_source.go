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
	_ datasource.DataSource              = (*FirewallNATOneToOneRulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATOneToOneRulesDataSource)(nil)
)

type FirewallNATOneToOneRulesModel struct {
	Rules types.List `tfsdk:"rules"`
}

func NewFirewallNATOneToOneRulesDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATOneToOneRulesDataSource{}
}

type FirewallNATOneToOneRulesDataSource struct {
	client *pfsense.Client
}

func (m *FirewallNATOneToOneRulesModel) Set(ctx context.Context, rules pfsense.NATOneToOnes) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []FirewallNATOneToOneRuleModel{}
	for _, r := range rules {
		var model FirewallNATOneToOneRuleModel
		diags.Append(model.Set(ctx, r)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallNATOneToOneRuleModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.Rules = listValue

	return diags
}

func (d *FirewallNATOneToOneRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_one_to_one_rules", req.ProviderTypeName)
}

func (d *FirewallNATOneToOneRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all NAT 1:1 rules.",
		MarkdownDescription: "Retrieves all [NAT 1:1](https://docs.netgate.com/pfsense/en/latest/nat/1-1.html) rules.",
		Attributes: map[string]schema.Attribute{
			"rules": schema.ListNestedAttribute{
				Description: "List of NAT 1:1 rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["interface"].Description,
							Computed:    true,
						},
						"external": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["external"].Description,
							Computed:    true,
						},
						"ipprotocol": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["ipprotocol"].Description,
							Computed:    true,
						},
						"source_address": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["source_address"].Description,
							Computed:    true,
						},
						"source_not": schema.BoolAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["source_not"].Description,
							Computed:    true,
						},
						"destination_address": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["destination_address"].Description,
							Computed:    true,
						},
						"destination_not": schema.BoolAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["destination_not"].Description,
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["disabled"].Description,
							Computed:    true,
						},
						"no_binat": schema.BoolAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["no_binat"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"nat_reflection": schema.StringAttribute{
							Description: FirewallNATOneToOneRuleModel{}.descriptions()["nat_reflection"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallNATOneToOneRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATOneToOneRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATOneToOneRulesModel

	rules, err := d.client.GetNATOneToOnes(ctx)
	if addError(&resp.Diagnostics, "Unable to get NAT 1:1 rules", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rules)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
