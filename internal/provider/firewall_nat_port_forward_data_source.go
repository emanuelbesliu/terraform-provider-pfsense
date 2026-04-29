package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNATPortForwardDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNATPortForwardDataSource)(nil)
)

type FirewallNATPortForwardDataSourceModel struct {
	FirewallNATPortForwardModel
}

func NewFirewallNATPortForwardDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNATPortForwardDataSource{}
}

type FirewallNATPortForwardDataSource struct {
	client *pfsense.Client
}

func (d *FirewallNATPortForwardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_port_forward", req.ProviderTypeName)
}

func (d *FirewallNATPortForwardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single NAT port forward rule by its description.",
		MarkdownDescription: "Retrieves a single [NAT port forward](https://docs.netgate.com/pfsense/en/latest/nat/port-forwards.html) rule by its description.",
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
				Required:    true,
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
	}
}

func (d *FirewallNATPortForwardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallNATPortForwardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNATPortForwardDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetNATPortForward(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Unable to get NAT port forward", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
