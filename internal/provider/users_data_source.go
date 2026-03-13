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
	_ datasource.DataSource              = (*UsersDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*UsersDataSource)(nil)
)

type UsersModel struct {
	Users types.List `tfsdk:"users"`
}

func NewUsersDataSource() datasource.DataSource { //nolint:ireturn
	return &UsersDataSource{}
}

type UsersDataSource struct {
	client *pfsense.Client
}

func (m *UsersModel) Set(ctx context.Context, users pfsense.Users) diag.Diagnostics {
	var diags diag.Diagnostics

	userModels := []UserModel{}
	for _, u := range users {
		var userModel UserModel
		diags.Append(userModel.Set(ctx, u)...)
		userModels = append(userModels, userModel)
	}

	usersValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: UserModel{}.AttrTypes()}, userModels)
	diags.Append(newDiags...)
	m.Users = usersValue

	return diags
}

func (d *UsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_users", req.ProviderTypeName)
}

func (d *UsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := UserModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves all pfSense local user accounts.",
		MarkdownDescription: "Retrieves all pfSense local [user](https://docs.netgate.com/pfsense/en/latest/usermanager/index.html) accounts.",
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Description: "List of user accounts.",
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
						"uid": schema.StringAttribute{
							Description: descriptions["uid"].Description,
							Computed:    true,
						},
						"disabled": schema.BoolAttribute{
							Description: descriptions["disabled"].Description,
							Computed:    true,
						},
						"expires": schema.StringAttribute{
							Description: descriptions["expires"].Description,
							Computed:    true,
						},
						"authorized_keys": schema.StringAttribute{
							Description: descriptions["authorized_keys"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"ipsec_psk": schema.StringAttribute{
							Description: descriptions["ipsec_psk"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"privileges": schema.ListAttribute{
							Description: descriptions["privileges"].Description,
							Computed:    true,
							ElementType: types.StringType,
						},
						"groups": schema.ListAttribute{
							Description: descriptions["groups"].Description,
							Computed:    true,
							ElementType: types.StringType,
						},
						"custom_settings": schema.BoolAttribute{
							Description: descriptions["custom_settings"].Description,
							Computed:    true,
						},
						"webgui_css": schema.StringAttribute{
							Description: descriptions["webgui_css"].Description,
							Computed:    true,
						},
						"dashboard_columns": schema.StringAttribute{
							Description: descriptions["dashboard_columns"].Description,
							Computed:    true,
						},
						"keep_history": schema.BoolAttribute{
							Description: descriptions["keep_history"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *UsersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersModel

	users, err := d.client.GetUsers(ctx)
	if addError(&resp.Diagnostics, "Unable to get users", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *users)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
