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
	_ datasource.DataSource              = (*UserDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*UserDataSource)(nil)
)

type UserDataSourceModel struct {
	UserModel
}

func NewUserDataSource() datasource.DataSource { //nolint:ireturn
	return &UserDataSource{}
}

type UserDataSource struct {
	client *pfsense.Client
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := UserModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single pfSense local user account by username.",
		MarkdownDescription: "Retrieves a single pfSense local [user](https://docs.netgate.com/pfsense/en/latest/usermanager/index.html) account by username.",
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
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get user", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *user)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
