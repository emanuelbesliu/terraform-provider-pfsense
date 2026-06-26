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
	_ datasource.DataSource              = (*FirewallNATNPtRulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATNPtRulesDataSource)(nil)
)

type FirewallNATNPtRulesModel struct {
	Rules types.List `tfsdk:"rules"`
}

func NewFirewallNATNPtRulesDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATNPtRulesDataSource{}
}

type FirewallNATNPtRulesDataSource struct {
	client *pfsense.Client
}

func (m *FirewallNATNPtRulesModel) Set(ctx context.Context, rules pfsense.NATNPts) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []FirewallNATNPtRuleModel{}
	for _, r := range rules {
		var model FirewallNATNPtRuleModel
		diags.Append(model.Set(ctx, r)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallNATNPtRuleModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.Rules = listValue

	return diags
}

func (d *FirewallNATNPtRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_npt_rules", req.ProviderTypeName)
}

func (d *FirewallNATNPtRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all NAT NPt rules.",
		MarkdownDescription: "Retrieves all [NAT NPt](https://docs.netgate.com/pfsense/en/latest/nat/npt.html) rules.",
		Attributes: map[string]schema.Attribute{
			"rules": schema.ListNestedAttribute{
				Description: "List of NAT NPt rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface": schema.StringAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["interface"].Description,
							Computed:    true,
						},
						"source_prefix": schema.StringAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["source_prefix"].Description,
							Computed:    true,
						},
						"source_not": schema.BoolAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["source_not"].Description,
							Computed:    true,
						},
						"destination_prefix": schema.StringAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["destination_prefix"].Description,
							Computed:    true,
						},
						"destination_not": schema.BoolAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["destination_not"].Description,
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["disabled"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallNATNPtRuleModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallNATNPtRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATNPtRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATNPtRulesModel

	rules, err := d.client.GetNATNPts(ctx)
	if addError(&resp.Diagnostics, "Unable to get NAT NPt rules", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rules)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
