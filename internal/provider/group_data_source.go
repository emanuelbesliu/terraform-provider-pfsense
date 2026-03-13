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
	_ datasource.DataSource              = (*GroupDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*GroupDataSource)(nil)
)

type GroupDataSourceModel struct {
	GroupModel
}

func NewGroupDataSource() datasource.DataSource { //nolint:ireturn
	return &GroupDataSource{}
}

type GroupDataSource struct {
	client *pfsense.Client
}

func (d *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_group", req.ProviderTypeName)
}

func (d *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := GroupModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single pfSense local group by name.",
		MarkdownDescription: "Retrieves a single pfSense local [group](https://docs.netgate.com/pfsense/en/latest/usermanager/groups.html) by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: descriptions["name"].Description,
				Required:    true,
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
	}
}

func (d *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
