package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNATNPtRuleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATNPtRuleDataSource)(nil)
)

type FirewallNATNPtRuleDataSourceModel struct {
	FirewallNATNPtRuleModel
}

func NewFirewallNATNPtRuleDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATNPtRuleDataSource{}
}

type FirewallNATNPtRuleDataSource struct {
	client *pfsense.Client
}

func (d *FirewallNATNPtRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_npt", req.ProviderTypeName)
}

func (d *FirewallNATNPtRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single NAT NPt rule by its description.",
		MarkdownDescription: "Retrieves a single [NAT NPt](https://docs.netgate.com/pfsense/en/latest/nat/npt.html) rule by its description.",
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
				Required:    true,
			},
		},
	}
}

func (d *FirewallNATNPtRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATNPtRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATNPtRuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetNATNPt(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Unable to get NAT NPt rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
