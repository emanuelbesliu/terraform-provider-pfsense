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
	_ datasource.DataSource              = (*FirewallNATPortForwardsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATPortForwardsDataSource)(nil)
)

type FirewallNATPortForwardsModel struct {
	PortForwards types.List `tfsdk:"port_forwards"`
}

func NewFirewallNATPortForwardsDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATPortForwardsDataSource{}
}

type FirewallNATPortForwardsDataSource struct {
	client *pfsense.Client
}

func (m *FirewallNATPortForwardsModel) Set(ctx context.Context, rules pfsense.NATPortForwards) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []FirewallNATPortForwardModel{}
	for _, r := range rules {
		var model FirewallNATPortForwardModel
		diags.Append(model.Set(ctx, r)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallNATPortForwardModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.PortForwards = listValue

	return diags
}

func (d *FirewallNATPortForwardsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_port_forwards", req.ProviderTypeName)
}

func (d *FirewallNATPortForwardsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all NAT port forward rules.",
		MarkdownDescription: "Retrieves all [NAT port forward](https://docs.netgate.com/pfsense/en/latest/nat/port-forwards.html) rules.",
		Attributes: map[string]schema.Attribute{
			"port_forwards": schema.ListNestedAttribute{
				Description: "List of NAT port forward rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["interface"].Description,
							Computed:    true,
						},
						"ipprotocol": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["ipprotocol"].Description,
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["protocol"].Description,
							Computed:    true,
						},
						"source_address": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["source_address"].Description,
							Computed:    true,
						},
						"source_port": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["source_port"].Description,
							Computed:    true,
						},
						"source_not": schema.BoolAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["source_not"].Description,
							Computed:    true,
						},
						"destination_address": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["destination_address"].Description,
							Computed:    true,
						},
						"destination_port": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["destination_port"].Description,
							Computed:    true,
						},
						"destination_not": schema.BoolAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["destination_not"].Description,
							Computed:    true,
						},
						"target": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["target"].Description,
							Computed:    true,
						},
						"local_port": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["local_port"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["disabled"].Description,
							Computed:    true,
						},
						"no_rdr": schema.BoolAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["no_rdr"].Description,
							Computed:    true,
						},
						"nat_reflection": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["nat_reflection"].Description,
							Computed:    true,
						},
						"associated_rule_id": schema.StringAttribute{
							Description: FirewallNATPortForwardModel{}.descriptions()["associated_rule_id"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallNATPortForwardsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATPortForwardsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATPortForwardsModel

	rules, err := d.client.GetNATPortForwards(ctx)
	if addError(&resp.Diagnostics, "Unable to get NAT port forwards", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rules)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
