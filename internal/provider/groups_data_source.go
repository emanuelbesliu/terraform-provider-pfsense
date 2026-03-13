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
	_ datasource.DataSource              = (*GroupsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*GroupsDataSource)(nil)
)

type GroupsModel struct {
	Groups types.List `tfsdk:"groups"`
}

func NewGroupsDataSource() datasource.DataSource { //nolint:ireturn
	return &GroupsDataSource{}
}

type GroupsDataSource struct {
	client *pfsense.Client
}

func (m *GroupsModel) Set(ctx context.Context, groups pfsense.Groups) diag.Diagnostics {
	var diags diag.Diagnostics

	groupModels := []GroupModel{}
	for _, g := range groups {
		var groupModel GroupModel
		diags.Append(groupModel.Set(ctx, g)...)
		groupModels = append(groupModels, groupModel)
	}

	groupsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GroupModel{}.AttrTypes()}, groupModels)
	diags.Append(newDiags...)
	m.Groups = groupsValue

	return diags
}

func (d *GroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_groups", req.ProviderTypeName)
}

func (d *GroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := GroupModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves all pfSense local groups.",
		MarkdownDescription: "Retrieves all pfSense local [groups](https://docs.netgate.com/pfsense/en/latest/usermanager/groups.html).",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "List of groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: descriptions["name"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: descriptions["description"].Description,
							Computed:    true,
						},
						"scope": schema.StringAttribute{
							Description: descriptions["scope"].Description,
							Computed:    true,
						},
						"gid": schema.StringAttribute{
							Description: descriptions["gid"].Description,
							Computed:    true,
						},
						"members": schema.ListAttribute{
							Description: descriptions["members"].Description,
							Computed:    true,
							ElementType: types.StringType,
						},
						"privileges": schema.ListAttribute{
							Description: descriptions["privileges"].Description,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *GroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *GroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupsModel

	groups, err := d.client.GetGroups(ctx)
	if addError(&resp.Diagnostics, "Unable to get groups", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *groups)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
