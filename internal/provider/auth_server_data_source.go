package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*AuthServerDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*AuthServerDataSource)(nil)
)

type AuthServerDataSourceModel struct {
	AuthServerModel
}

func NewAuthServerDataSource() datasource.DataSource { //nolint:ireturn
	return &AuthServerDataSource{}
}

type AuthServerDataSource struct {
	client *pfsense.Client
}

func (d *AuthServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_auth_server", req.ProviderTypeName)
}

func (d *AuthServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := AuthServerModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single pfSense authentication server by name.",
		MarkdownDescription: "Retrieves a single pfSense [authentication server](https://docs.netgate.com/pfsense/en/latest/usermanager/authservers.html) by name.",
		Attributes:          authServerDataSourceAttributes(descriptions, false),
	}
}

func (d *AuthServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *AuthServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AuthServerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := d.client.GetAuthServer(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get auth server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// authServerDataSourceAttributes returns the common attribute schema for auth server data sources.
// When nested is true, all fields are Computed (for the list data source).
// When nested is false, "name" is Required (for the single data source lookup).
func authServerDataSourceAttributes(descriptions map[string]attrDescription, nested bool) map[string]schema.Attribute {
	nameAttr := schema.StringAttribute{
		Description: descriptions["name"].Description,
		Computed:    true,
	}
	if !nested {
		nameAttr = schema.StringAttribute{
			Description: descriptions["name"].Description,
			Required:    true,
		}
	}

	return map[string]schema.Attribute{
		"name": nameAttr,
		"type": schema.StringAttribute{
			Description: descriptions["type"].Description,
			Computed:    true,
		},
		"host": schema.StringAttribute{
			Description: descriptions["host"].Description,
			Computed:    true,
		},
		"refid": schema.StringAttribute{
			Description: descriptions["refid"].Description,
			Computed:    true,
		},
		"ldap_port": schema.StringAttribute{
			Description: descriptions["ldap_port"].Description,
			Computed:    true,
		},
		"ldap_urltype": schema.StringAttribute{
			Description: descriptions["ldap_urltype"].Description,
			Computed:    true,
		},
		"ldap_protver": schema.StringAttribute{
			Description: descriptions["ldap_protver"].Description,
			Computed:    true,
		},
		"ldap_scope": schema.StringAttribute{
			Description: descriptions["ldap_scope"].Description,
			Computed:    true,
		},
		"ldap_basedn": schema.StringAttribute{
			Description: descriptions["ldap_basedn"].Description,
			Computed:    true,
		},
		"ldap_authcn": schema.StringAttribute{
			Description: descriptions["ldap_authcn"].Description,
			Computed:    true,
		},
		"ldap_binddn": schema.StringAttribute{
			Description: descriptions["ldap_binddn"].Description,
			Computed:    true,
		},
		"ldap_bindpw": schema.StringAttribute{
			Description: descriptions["ldap_bindpw"].Description,
			Computed:    true,
			Sensitive:   true,
		},
		"ldap_caref": schema.StringAttribute{
			Description: descriptions["ldap_caref"].Description,
			Computed:    true,
		},
		"ldap_timeout": schema.StringAttribute{
			Description: descriptions["ldap_timeout"].Description,
			Computed:    true,
		},
		"ldap_extended_enabled": schema.BoolAttribute{
			Description: descriptions["ldap_extended_enabled"].Description,
			Computed:    true,
		},
		"ldap_extended_query": schema.StringAttribute{
			Description: descriptions["ldap_extended_query"].Description,
			Computed:    true,
		},
		"ldap_attr_user": schema.StringAttribute{
			Description: descriptions["ldap_attr_user"].Description,
			Computed:    true,
		},
		"ldap_attr_group": schema.StringAttribute{
			Description: descriptions["ldap_attr_group"].Description,
			Computed:    true,
		},
		"ldap_attr_member": schema.StringAttribute{
			Description: descriptions["ldap_attr_member"].Description,
			Computed:    true,
		},
		"ldap_attr_groupobj": schema.StringAttribute{
			Description: descriptions["ldap_attr_groupobj"].Description,
			Computed:    true,
		},
		"ldap_pam_groupdn": schema.StringAttribute{
			Description: descriptions["ldap_pam_groupdn"].Description,
			Computed:    true,
		},
		"ldap_utf8": schema.BoolAttribute{
			Description: descriptions["ldap_utf8"].Description,
			Computed:    true,
		},
		"ldap_nostrip_at": schema.BoolAttribute{
			Description: descriptions["ldap_nostrip_at"].Description,
			Computed:    true,
		},
		"ldap_allow_unauthenticated": schema.BoolAttribute{
			Description: descriptions["ldap_allow_unauthenticated"].Description,
			Computed:    true,
		},
		"ldap_rfc2307": schema.BoolAttribute{
			Description: descriptions["ldap_rfc2307"].Description,
			Computed:    true,
		},
		"ldap_rfc2307_userdn": schema.BoolAttribute{
			Description: descriptions["ldap_rfc2307_userdn"].Description,
			Computed:    true,
		},
		"radius_protocol": schema.StringAttribute{
			Description: descriptions["radius_protocol"].Description,
			Computed:    true,
		},
		"radius_auth_port": schema.StringAttribute{
			Description: descriptions["radius_auth_port"].Description,
			Computed:    true,
		},
		"radius_acct_port": schema.StringAttribute{
			Description: descriptions["radius_acct_port"].Description,
			Computed:    true,
		},
		"radius_secret": schema.StringAttribute{
			Description: descriptions["radius_secret"].Description,
			Computed:    true,
			Sensitive:   true,
		},
		"radius_timeout": schema.StringAttribute{
			Description: descriptions["radius_timeout"].Description,
			Computed:    true,
		},
		"radius_nasip_attribute": schema.StringAttribute{
			Description: descriptions["radius_nasip_attribute"].Description,
			Computed:    true,
		},
		"radius_srvcs": schema.StringAttribute{
			Description: descriptions["radius_srvcs"].Description,
			Computed:    true,
		},
	}
}
