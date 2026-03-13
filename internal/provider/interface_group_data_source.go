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
	_ datasource.DataSource              = (*InterfaceGroupDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfaceGroupDataSource)(nil)
)

type InterfaceGroupDataSourceModel struct {
	InterfaceGroupModel
}

func NewInterfaceGroupDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfaceGroupDataSource{}
}

type InterfaceGroupDataSource struct {
	client *pfsense.Client
}

func (d *InterfaceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_group", req.ProviderTypeName)
}

func (d *InterfaceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single interface group by name. Interface groups allow applying firewall rules to multiple interfaces at once.",
		MarkdownDescription: "Retrieves a single [interface group](https://docs.netgate.com/pfsense/en/latest/interfaces/groups.html) by name. Interface groups allow applying firewall rules to multiple interfaces at once.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: InterfaceGroupModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"members": schema.ListAttribute{
				Description: InterfaceGroupModel{}.descriptions()["members"].Description,
				ElementType: types.StringType,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: InterfaceGroupModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *InterfaceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfaceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetInterfaceGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get interface group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
