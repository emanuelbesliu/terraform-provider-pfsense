package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemAdvancedAdminDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemAdvancedAdminDataSource)(nil)
)

type SystemAdvancedAdminDataSourceModel struct {
	SystemAdvancedAdminModel
}

func NewSystemAdvancedAdminDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemAdvancedAdminDataSource{}
}

type SystemAdvancedAdminDataSource struct {
	client *pfsense.Client
}

func (d *SystemAdvancedAdminDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_admin", req.ProviderTypeName)
}

func (d *SystemAdvancedAdminDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := SystemAdvancedAdminModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Retrieves the system advanced admin access configuration including webConfigurator, SSH, login protection, serial console, and console settings.",
		MarkdownDescription: "Retrieves the [system advanced admin access](https://docs.netgate.com/pfsense/en/latest/config/advanced/admin.html) configuration including webConfigurator, SSH, login protection, serial console, and console settings.",
		Attributes: map[string]schema.Attribute{
			// webConfigurator
			"webgui_protocol": schema.StringAttribute{
				Description:         descriptions["webgui_protocol"].Description,
				MarkdownDescription: descriptions["webgui_protocol"].MarkdownDescription,
				Computed:            true,
			},
			"ssl_certificate": schema.StringAttribute{
				Description: descriptions["ssl_certificate"].Description,
				Computed:    true,
			},
			"webgui_port": schema.Int64Attribute{
				Description: descriptions["webgui_port"].Description,
				Computed:    true,
			},
			"max_processes": schema.Int64Attribute{
				Description: descriptions["max_processes"].Description,
				Computed:    true,
			},
			"disable_http_redirect": schema.BoolAttribute{
				Description: descriptions["disable_http_redirect"].Description,
				Computed:    true,
			},
			"disable_hsts": schema.BoolAttribute{
				Description: descriptions["disable_hsts"].Description,
				Computed:    true,
			},
			"ocsp_staple": schema.BoolAttribute{
				Description: descriptions["ocsp_staple"].Description,
				Computed:    true,
			},
			"login_autocomplete": schema.BoolAttribute{
				Description: descriptions["login_autocomplete"].Description,
				Computed:    true,
			},
			"quiet_login": schema.BoolAttribute{
				Description: descriptions["quiet_login"].Description,
				Computed:    true,
			},
			"roaming": schema.BoolAttribute{
				Description: descriptions["roaming"].Description,
				Computed:    true,
			},
			"disable_anti_lockout": schema.BoolAttribute{
				Description: descriptions["disable_anti_lockout"].Description,
				Computed:    true,
			},
			"disable_dns_rebind_check": schema.BoolAttribute{
				Description: descriptions["disable_dns_rebind_check"].Description,
				Computed:    true,
			},
			"disable_http_referer_check": schema.BoolAttribute{
				Description: descriptions["disable_http_referer_check"].Description,
				Computed:    true,
			},
			"alternate_hostnames": schema.StringAttribute{
				Description: descriptions["alternate_hostnames"].Description,
				Computed:    true,
			},
			"page_name_first": schema.BoolAttribute{
				Description: descriptions["page_name_first"].Description,
				Computed:    true,
			},

			// SSH
			"ssh_enabled": schema.BoolAttribute{
				Description: descriptions["ssh_enabled"].Description,
				Computed:    true,
			},
			"sshd_key_only": schema.StringAttribute{
				Description:         descriptions["sshd_key_only"].Description,
				MarkdownDescription: descriptions["sshd_key_only"].MarkdownDescription,
				Computed:            true,
			},
			"sshd_agent_forwarding": schema.BoolAttribute{
				Description: descriptions["sshd_agent_forwarding"].Description,
				Computed:    true,
			},
			"ssh_port": schema.Int64Attribute{
				Description: descriptions["ssh_port"].Description,
				Computed:    true,
			},

			// Login Protection
			"login_protection_threshold": schema.Int64Attribute{
				Description: descriptions["login_protection_threshold"].Description,
				Computed:    true,
			},
			"login_protection_blocktime": schema.Int64Attribute{
				Description: descriptions["login_protection_blocktime"].Description,
				Computed:    true,
			},
			"login_protection_detection_time": schema.Int64Attribute{
				Description: descriptions["login_protection_detection_time"].Description,
				Computed:    true,
			},
			"login_protection_pass_list": schema.StringAttribute{
				Description: descriptions["login_protection_pass_list"].Description,
				Computed:    true,
			},

			// Serial
			"serial_terminal": schema.BoolAttribute{
				Description: descriptions["serial_terminal"].Description,
				Computed:    true,
			},
			"serial_speed": schema.Int64Attribute{
				Description:         descriptions["serial_speed"].Description,
				MarkdownDescription: descriptions["serial_speed"].MarkdownDescription,
				Computed:            true,
			},
			"primary_console": schema.StringAttribute{
				Description:         descriptions["primary_console"].Description,
				MarkdownDescription: descriptions["primary_console"].MarkdownDescription,
				Computed:            true,
			},

			// Console
			"disable_console_menu": schema.BoolAttribute{
				Description: descriptions["disable_console_menu"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemAdvancedAdminDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemAdvancedAdminDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemAdvancedAdminDataSourceModel

	a, err := d.client.GetAdvancedAdmin(ctx)
	if addError(&resp.Diagnostics, "Unable to get system advanced admin settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
