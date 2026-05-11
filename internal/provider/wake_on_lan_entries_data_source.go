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
	_ datasource.DataSource              = (*WakeOnLanEntriesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*WakeOnLanEntriesDataSource)(nil)
)

type WakeOnLanEntriesModel struct {
	Entries types.List `tfsdk:"entries"`
}

func NewWakeOnLanEntriesDataSource() datasource.DataSource { //nolint:ireturn
	return &WakeOnLanEntriesDataSource{}
}

type WakeOnLanEntriesDataSource struct {
	client *pfsense.Client
}

func (m *WakeOnLanEntriesModel) Set(ctx context.Context, entries pfsense.WakeOnLanEntries) diag.Diagnostics {
	var diags diag.Diagnostics

	entryModels := []WakeOnLanModel{}
	for _, e := range entries {
		var entryModel WakeOnLanModel
		diags.Append(entryModel.Set(ctx, e)...)
		entryModels = append(entryModels, entryModel)
	}

	entriesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: WakeOnLanModel{}.AttrTypes()}, entryModels)
	diags.Append(newDiags...)
	m.Entries = entriesValue

	return diags
}

func (d *WakeOnLanEntriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_wake_on_lan_entries", req.ProviderTypeName)
}

func (d *WakeOnLanEntriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all Wake-on-LAN entries.",
		MarkdownDescription: "Retrieves all [Wake-on-LAN](https://docs.netgate.com/pfsense/en/latest/services/wake-on-lan.html) entries.",
		Attributes: map[string]schema.Attribute{
			"entries": schema.ListNestedAttribute{
				Description: "List of Wake-on-LAN entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"interface": schema.StringAttribute{
							Description: WakeOnLanModel{}.descriptions()["interface"].Description,
							Computed:    true,
						},
						"mac": schema.StringAttribute{
							Description: WakeOnLanModel{}.descriptions()["mac"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: WakeOnLanModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WakeOnLanEntriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *WakeOnLanEntriesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WakeOnLanEntriesModel

	entries, err := d.client.GetWakeOnLanEntries(ctx)
	if addError(&resp.Diagnostics, "Unable to get wake on lan entries", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entries)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
