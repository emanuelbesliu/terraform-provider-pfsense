package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNATOneToOneRuleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATOneToOneRuleDataSource)(nil)
)

type FirewallNATOneToOneRuleDataSourceModel struct {
	FirewallNATOneToOneRuleModel
}

func NewFirewallNATOneToOneRuleDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATOneToOneRuleDataSource{}
}

type FirewallNATOneToOneRuleDataSource struct {
	client *pfsense.Client
}

func (d *FirewallNATOneToOneRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_one_to_one", req.ProviderTypeName)
}

func (d *FirewallNATOneToOneRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single NAT 1:1 rule by its description.",
		MarkdownDescription: "Retrieves a single [NAT 1:1](https://docs.netgate.com/pfsense/en/latest/nat/1-1.html) rule by its description.",
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
				Required:    true,
			},
			"nat_reflection": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["nat_reflection"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *FirewallNATOneToOneRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATOneToOneRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATOneToOneRuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetNATOneToOne(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Unable to get NAT 1:1 rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
