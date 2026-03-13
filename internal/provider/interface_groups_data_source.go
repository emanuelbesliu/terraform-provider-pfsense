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
	_ datasource.DataSource              = (*InterfaceGroupsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfaceGroupsDataSource)(nil)
)

type InterfaceGroupsModel struct {
	InterfaceGroups types.List `tfsdk:"interface_groups"`
}

func NewInterfaceGroupsDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfaceGroupsDataSource{}
}

type InterfaceGroupsDataSource struct {
	client *pfsense.Client
}

func (m *InterfaceGroupsModel) Set(ctx context.Context, groups pfsense.InterfaceGroups) diag.Diagnostics {
	var diags diag.Diagnostics

	groupModels := []InterfaceGroupModel{}
	for _, g := range groups {
		var groupModel InterfaceGroupModel
		diags.Append(groupModel.Set(ctx, g)...)
		groupModels = append(groupModels, groupModel)
	}

	groupsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InterfaceGroupModel{}.AttrTypes()}, groupModels)
	diags.Append(newDiags...)
	m.InterfaceGroups = groupsValue

	return diags
}

func (d *InterfaceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_groups", req.ProviderTypeName)
}

func (d *InterfaceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all interface groups. Interface groups allow applying firewall rules to multiple interfaces at once.",
		MarkdownDescription: "Retrieves all [interface groups](https://docs.netgate.com/pfsense/en/latest/interfaces/groups.html). Interface groups allow applying firewall rules to multiple interfaces at once.",
		Attributes: map[string]schema.Attribute{
			"interface_groups": schema.ListNestedAttribute{
				Description: "List of interface groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: InterfaceGroupModel{}.descriptions()["name"].Description,
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *InterfaceGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfaceGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceGroupsModel

	groups, err := d.client.GetInterfaceGroups(ctx)
	if addError(&resp.Diagnostics, "Unable to get interface groups", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *groups)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
