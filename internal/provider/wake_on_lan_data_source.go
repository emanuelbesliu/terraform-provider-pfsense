package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*WakeOnLanDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*WakeOnLanDataSource)(nil)
)

type WakeOnLanDataSourceModel struct {
	WakeOnLanModel
}

func NewWakeOnLanDataSource() datasource.DataSource { //nolint:ireturn
	return &WakeOnLanDataSource{}
}

type WakeOnLanDataSource struct {
	client *pfsense.Client
}

func (d *WakeOnLanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_wake_on_lan", req.ProviderTypeName)
}

func (d *WakeOnLanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single Wake-on-LAN entry by its MAC address.",
		MarkdownDescription: "Retrieves a single [Wake-on-LAN](https://docs.netgate.com/pfsense/en/latest/services/wake-on-lan.html) entry by its MAC address.",
		Attributes: map[string]schema.Attribute{
			"mac": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["mac"].Description,
				Required:    true,
			},
			"interface": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["interface"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *WakeOnLanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *WakeOnLanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WakeOnLanDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := d.client.GetWakeOnLanEntry(ctx, data.MAC.ValueString())
	if addError(&resp.Diagnostics, "Unable to get wake on lan entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
