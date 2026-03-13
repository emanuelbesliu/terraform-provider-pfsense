package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*InterfaceDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfaceDataSource)(nil)
)

type InterfaceDataSourceModel struct {
	InterfaceModel
}

func NewInterfaceDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfaceDataSource{}
}

type InterfaceDataSource struct {
	client *pfsense.Client
}

func (d *InterfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface", req.ProviderTypeName)
}

func (d *InterfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single interface assignment by logical name (e.g. 'wan', 'lan', 'opt1').",
		MarkdownDescription: "Retrieves a single [interface assignment](https://docs.netgate.com/pfsense/en/latest/interfaces/index.html) by logical name (e.g. `wan`, `lan`, `opt1`).",
		Attributes: map[string]schema.Attribute{
			"logical_name": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["logical_name"].Description,
				Required:    true,
			},
			"port": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["port"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["enabled"].Description,
				Computed:    true,
			},
			"ipv4_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv4_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv4_type"].MarkdownDescription,
				Computed:            true,
			},
			"ipv4_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_address"].Description,
				Computed:    true,
			},
			"ipv4_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_subnet"].Description,
				Computed:    true,
			},
			"ipv4_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_gateway"].Description,
				Computed:    true,
			},
			"ipv6_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv6_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv6_type"].MarkdownDescription,
				Computed:            true,
			},
			"ipv6_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_address"].Description,
				Computed:    true,
			},
			"ipv6_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_subnet"].Description,
				Computed:    true,
			},
			"ipv6_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_gateway"].Description,
				Computed:    true,
			},
			"spoof_mac": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["spoof_mac"].Description,
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mtu"].Description,
				Computed:    true,
			},
			"mss": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mss"].Description,
				Computed:    true,
			},
			"block_private": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_private"].Description,
				Computed:    true,
			},
			"block_bogons": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_bogons"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *InterfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := d.client.GetInterface(ctx, data.LogicalName.ValueString())
	if addError(&resp.Diagnostics, "Unable to get interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
