package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemTunableDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemTunableDataSource)(nil)
)

type SystemTunableDataSourceModel struct {
	SystemTunableModel
}

func NewSystemTunableDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemTunableDataSource{}
}

type SystemTunableDataSource struct {
	client *pfsense.Client
}

func (d *SystemTunableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_tunable", req.ProviderTypeName)
}

func (d *SystemTunableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single system tunable (sysctl) by name. Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		MarkdownDescription: "Retrieves a single system [tunable](https://docs.netgate.com/pfsense/en/latest/system/advanced-tunables.html) (sysctl) by name. Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		Attributes: map[string]schema.Attribute{
			"tunable": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["tunable"].Description,
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["value"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemTunableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemTunableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemTunableDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tunable, err := d.client.GetSystemTunable(ctx, data.Tunable.ValueString())
	if addError(&resp.Diagnostics, "Unable to get system tunable", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *tunable)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
