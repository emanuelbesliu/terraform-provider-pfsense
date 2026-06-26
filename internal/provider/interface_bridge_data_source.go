package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*InterfaceBridgeDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfaceBridgeDataSource)(nil)
)

type InterfaceBridgeDataSourceModel struct {
	InterfaceBridgeModel
}

func NewInterfaceBridgeDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfaceBridgeDataSource{}
}

type InterfaceBridgeDataSource struct {
	client *pfsense.Client
}

func interfaceBridgeDataSourceAttributes(bridgeIfRequired bool) map[string]schema.Attribute {
	descriptions := InterfaceBridgeModel{}.descriptions()

	stringList := func(description string) schema.ListAttribute {
		return schema.ListAttribute{
			Description: description,
			ElementType: types.StringType,
			Computed:    true,
		}
	}

	return map[string]schema.Attribute{
		"bridge_if": schema.StringAttribute{
			Description: descriptions["bridge_if"].Description,
			Required:    bridgeIfRequired,
			Computed:    !bridgeIfRequired,
		},
		"members":              stringList(descriptions["members"].Description),
		"description":          schema.StringAttribute{Description: descriptions["description"].Description, Computed: true},
		"enable_stp":           schema.BoolAttribute{Description: descriptions["enable_stp"].Description, Computed: true},
		"ip6_link_local":       schema.BoolAttribute{Description: descriptions["ip6_link_local"].Description, Computed: true},
		"protocol":             schema.StringAttribute{Description: descriptions["protocol"].Description, Computed: true},
		"priority":             schema.Int64Attribute{Description: descriptions["priority"].Description, Computed: true},
		"hello_time":           schema.Int64Attribute{Description: descriptions["hello_time"].Description, Computed: true},
		"forward_delay":        schema.Int64Attribute{Description: descriptions["forward_delay"].Description, Computed: true},
		"max_age":              schema.Int64Attribute{Description: descriptions["max_age"].Description, Computed: true},
		"hold_count":           schema.Int64Attribute{Description: descriptions["hold_count"].Description, Computed: true},
		"max_addresses":        schema.Int64Attribute{Description: descriptions["max_addresses"].Description, Computed: true},
		"cache_expire":         schema.Int64Attribute{Description: descriptions["cache_expire"].Description, Computed: true},
		"stp_interfaces":       stringList(descriptions["stp_interfaces"].Description),
		"static_interfaces":    stringList(descriptions["static_interfaces"].Description),
		"private_interfaces":   stringList(descriptions["private_interfaces"].Description),
		"span_interfaces":      stringList(descriptions["span_interfaces"].Description),
		"edge_interfaces":      stringList(descriptions["edge_interfaces"].Description),
		"auto_edge_interfaces": stringList(descriptions["auto_edge_interfaces"].Description),
		"ptp_interfaces":       stringList(descriptions["ptp_interfaces"].Description),
		"auto_ptp_interfaces":  stringList(descriptions["auto_ptp_interfaces"].Description),
	}
}

func (d *InterfaceBridgeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_bridge", req.ProviderTypeName)
}

func (d *InterfaceBridgeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single bridge interface by its assigned interface name.",
		MarkdownDescription: "Retrieves a single [bridge interface](https://docs.netgate.com/pfsense/en/latest/interfaces/bridges.html) by its assigned interface name.",
		Attributes:          interfaceBridgeDataSourceAttributes(true),
	}
}

func (d *InterfaceBridgeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfaceBridgeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceBridgeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bridge, err := d.client.GetBridge(ctx, data.BridgeIf.ValueString())
	if addError(&resp.Diagnostics, "Unable to get bridge", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *bridge)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
