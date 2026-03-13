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
	_ datasource.DataSource              = (*AuthServersDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*AuthServersDataSource)(nil)
)

type AuthServersModel struct {
	AuthServers types.List `tfsdk:"auth_servers"`
}

func NewAuthServersDataSource() datasource.DataSource { //nolint:ireturn
	return &AuthServersDataSource{}
}

type AuthServersDataSource struct {
	client *pfsense.Client
}

func (m *AuthServersModel) Set(ctx context.Context, servers pfsense.AuthServers) diag.Diagnostics {
	var diags diag.Diagnostics

	serverModels := []AuthServerModel{}
	for _, s := range servers {
		var serverModel AuthServerModel
		diags.Append(serverModel.Set(ctx, s)...)
		serverModels = append(serverModels, serverModel)
	}

	serversValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AuthServerModel{}.AttrTypes()}, serverModels)
	diags.Append(newDiags...)
	m.AuthServers = serversValue

	return diags
}

func (d *AuthServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_auth_servers", req.ProviderTypeName)
}

func (d *AuthServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := AuthServerModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves all pfSense authentication servers.",
		MarkdownDescription: "Retrieves all pfSense [authentication servers](https://docs.netgate.com/pfsense/en/latest/usermanager/authservers.html).",
		Attributes: map[string]schema.Attribute{
			"auth_servers": schema.ListNestedAttribute{
				Description: "List of authentication servers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: authServerDataSourceAttributes(descriptions, true),
				},
			},
		},
	}
}

func (d *AuthServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *AuthServersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AuthServersModel

	servers, err := d.client.GetAuthServers(ctx)
	if addError(&resp.Diagnostics, "Unable to get auth servers", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *servers)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
