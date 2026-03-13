package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemAdvancedNotificationsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemAdvancedNotificationsDataSource)(nil)
)

type SystemAdvancedNotificationsDataSourceModel struct {
	SystemAdvancedNotificationsModel
}

func NewSystemAdvancedNotificationsDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemAdvancedNotificationsDataSource{}
}

type SystemAdvancedNotificationsDataSource struct {
	client *pfsense.Client
}

func (d *SystemAdvancedNotificationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_notifications", req.ProviderTypeName)
}

func (d *SystemAdvancedNotificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := SystemAdvancedNotificationsModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Retrieves the system advanced notifications configuration including SMTP, Telegram, Pushover, Slack, certificate expiration, and sound settings.",
		MarkdownDescription: "Retrieves the [system advanced notifications](https://docs.netgate.com/pfsense/en/latest/config/advanced/notifications.html) configuration including SMTP, Telegram, Pushover, Slack, certificate expiration, and sound settings.",
		Attributes: map[string]schema.Attribute{
			// General Settings - Certificate Expiration
			"cert_enable_notify": schema.BoolAttribute{
				Description: descriptions["cert_enable_notify"].Description,
				Computed:    true,
			},
			"revoked_cert_ignore_notify": schema.BoolAttribute{
				Description: descriptions["revoked_cert_ignore_notify"].Description,
				Computed:    true,
			},
			"cert_expire_days": schema.Int64Attribute{
				Description: descriptions["cert_expire_days"].Description,
				Computed:    true,
			},

			// SMTP
			"disable_smtp": schema.BoolAttribute{
				Description: descriptions["disable_smtp"].Description,
				Computed:    true,
			},
			"smtp_ip_address": schema.StringAttribute{
				Description: descriptions["smtp_ip_address"].Description,
				Computed:    true,
			},
			"smtp_port": schema.Int64Attribute{
				Description: descriptions["smtp_port"].Description,
				Computed:    true,
			},
			"smtp_timeout": schema.Int64Attribute{
				Description: descriptions["smtp_timeout"].Description,
				Computed:    true,
			},
			"smtp_ssl": schema.BoolAttribute{
				Description: descriptions["smtp_ssl"].Description,
				Computed:    true,
			},
			"ssl_validate": schema.BoolAttribute{
				Description: descriptions["ssl_validate"].Description,
				Computed:    true,
			},
			"smtp_from_address": schema.StringAttribute{
				Description: descriptions["smtp_from_address"].Description,
				Computed:    true,
			},
			"smtp_notify_email_address": schema.StringAttribute{
				Description: descriptions["smtp_notify_email_address"].Description,
				Computed:    true,
			},
			"smtp_username": schema.StringAttribute{
				Description: descriptions["smtp_username"].Description,
				Computed:    true,
			},
			"smtp_password": schema.StringAttribute{
				Description: descriptions["smtp_password"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"smtp_auth_mechanism": schema.StringAttribute{
				Description:         descriptions["smtp_auth_mechanism"].Description,
				MarkdownDescription: descriptions["smtp_auth_mechanism"].MarkdownDescription,
				Computed:            true,
			},

			// Sounds
			"console_bell": schema.BoolAttribute{
				Description: descriptions["console_bell"].Description,
				Computed:    true,
			},
			"disable_beep": schema.BoolAttribute{
				Description: descriptions["disable_beep"].Description,
				Computed:    true,
			},

			// Telegram
			"telegram_enable": schema.BoolAttribute{
				Description: descriptions["telegram_enable"].Description,
				Computed:    true,
			},
			"telegram_api": schema.StringAttribute{
				Description: descriptions["telegram_api"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"telegram_chat_id": schema.StringAttribute{
				Description: descriptions["telegram_chat_id"].Description,
				Computed:    true,
			},

			// Pushover
			"pushover_enable": schema.BoolAttribute{
				Description: descriptions["pushover_enable"].Description,
				Computed:    true,
			},
			"pushover_api_key": schema.StringAttribute{
				Description: descriptions["pushover_api_key"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"pushover_user_key": schema.StringAttribute{
				Description: descriptions["pushover_user_key"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"pushover_sound": schema.StringAttribute{
				Description:         descriptions["pushover_sound"].Description,
				MarkdownDescription: descriptions["pushover_sound"].MarkdownDescription,
				Computed:            true,
			},
			"pushover_priority": schema.StringAttribute{
				Description:         descriptions["pushover_priority"].Description,
				MarkdownDescription: descriptions["pushover_priority"].MarkdownDescription,
				Computed:            true,
			},
			"pushover_retry": schema.Int64Attribute{
				Description: descriptions["pushover_retry"].Description,
				Computed:    true,
			},
			"pushover_expire": schema.Int64Attribute{
				Description: descriptions["pushover_expire"].Description,
				Computed:    true,
			},

			// Slack
			"slack_enable": schema.BoolAttribute{
				Description: descriptions["slack_enable"].Description,
				Computed:    true,
			},
			"slack_api": schema.StringAttribute{
				Description: descriptions["slack_api"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"slack_channel": schema.StringAttribute{
				Description: descriptions["slack_channel"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemAdvancedNotificationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemAdvancedNotificationsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemAdvancedNotificationsDataSourceModel

	a, err := d.client.GetAdvancedNotifications(ctx)
	if addError(&resp.Diagnostics, "Unable to get system advanced notifications settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
