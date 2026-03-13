package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DHCPv4StaticMappingDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DHCPv4StaticMappingDataSource)(nil)
)

type DHCPv4StaticMappingDataSourceModel struct {
	DHCPv4StaticMappingModel
}

func NewDHCPv4StaticMappingDataSource() datasource.DataSource { //nolint:ireturn
	return &DHCPv4StaticMappingDataSource{}
}

type DHCPv4StaticMappingDataSource struct {
	client *pfsense.Client
}

func (d *DHCPv4StaticMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dhcpv4_staticmapping", req.ProviderTypeName)
}

func (d *DHCPv4StaticMappingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single DHCPv4 static mapping by interface and MAC address. Static mappings express a preference for which IP address will be assigned to a given client based on its MAC address.",
		MarkdownDescription: "Retrieves a single DHCPv4 [static mapping](https://docs.netgate.com/pfsense/en/latest/services/dhcp/ipv4.html#static-mappings) by interface and MAC address. Static mappings express a preference for which IP address will be assigned to a given client based on its MAC address.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"mac_address": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["mac_address"].Description,
				CustomType:  macAddressType{},
				Required:    true,
			},
			"client_identifier": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["client_identifier"].Description,
				Computed:    true,
			},
			"ip_address": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["ip_address"].Description,
				Computed:    true,
			},
			"arp_table_static_entry": schema.BoolAttribute{
				Description:         DHCPv4StaticMappingModel{}.descriptions()["arp_table_static_entry"].Description,
				MarkdownDescription: DHCPv4StaticMappingModel{}.descriptions()["arp_table_static_entry"].MarkdownDescription,
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["hostname"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"wins_servers": schema.ListAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["wins_servers"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"dns_servers": schema.ListAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["dns_servers"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"gateway": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["gateway"].Description,
				Computed:    true,
			},
			"domain_name": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["domain_name"].Description,
				Computed:    true,
			},
			"domain_search_list": schema.ListAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["domain_search_list"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_lease_time": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["default_lease_time"].Description,
				Computed:    true,
				CustomType:  timetypes.GoDurationType{},
			},
			"maximum_lease_time": schema.StringAttribute{
				Description: DHCPv4StaticMappingModel{}.descriptions()["maximum_lease_time"].Description,
				Computed:    true,
				CustomType:  timetypes.GoDurationType{},
			},
		},
	}
}

func (d *DHCPv4StaticMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DHCPv4StaticMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DHCPv4StaticMappingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	macAddress, newDiags := data.MACAddress.parseMACAddress()
	resp.Diagnostics.Append(newDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	staticMapping, err := d.client.GetDHCPv4StaticMapping(ctx, data.Interface.ValueString(), macAddress)
	if addError(&resp.Diagnostics, "Unable to get static mapping", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *staticMapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
