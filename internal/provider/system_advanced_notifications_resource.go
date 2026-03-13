package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemAdvancedNotificationsResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemAdvancedNotificationsResource)(nil)
	_ resource.ResourceWithImportState = (*SystemAdvancedNotificationsResource)(nil)
)

type SystemAdvancedNotificationsResourceModel struct {
	SystemAdvancedNotificationsModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemAdvancedNotificationsResource() resource.Resource { //nolint:ireturn
	return &SystemAdvancedNotificationsResource{}
}

type SystemAdvancedNotificationsResource struct {
	client *pfsense.Client
}

func (r *SystemAdvancedNotificationsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_notifications", req.ProviderTypeName)
}

func (r *SystemAdvancedNotificationsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := SystemAdvancedNotificationsModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "System advanced notifications configuration including SMTP, Telegram, Pushover, Slack, certificate expiration, and sound settings. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System advanced notifications](https://docs.netgate.com/pfsense/en/latest/config/advanced/notifications.html) configuration including SMTP, Telegram, Pushover, Slack, certificate expiration, and sound settings. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			// General Settings - Certificate Expiration
			"cert_enable_notify": schema.BoolAttribute{
				Description: descriptions["cert_enable_notify"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"revoked_cert_ignore_notify": schema.BoolAttribute{
				Description: descriptions["revoked_cert_ignore_notify"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"cert_expire_days": schema.Int64Attribute{
				Description: descriptions["cert_expire_days"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedNotificationsCertExpireDays)),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},

			// SMTP
			"disable_smtp": schema.BoolAttribute{
				Description: descriptions["disable_smtp"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"smtp_ip_address": schema.StringAttribute{
				Description: descriptions["smtp_ip_address"].Description,
				Optional:    true,
			},
			"smtp_port": schema.Int64Attribute{
				Description: descriptions["smtp_port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
			"smtp_timeout": schema.Int64Attribute{
				Description: descriptions["smtp_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"smtp_ssl": schema.BoolAttribute{
				Description: descriptions["smtp_ssl"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ssl_validate": schema.BoolAttribute{
				Description: descriptions["ssl_validate"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"smtp_from_address": schema.StringAttribute{
				Description: descriptions["smtp_from_address"].Description,
				Optional:    true,
			},
			"smtp_notify_email_address": schema.StringAttribute{
				Description: descriptions["smtp_notify_email_address"].Description,
				Optional:    true,
			},
			"smtp_username": schema.StringAttribute{
				Description: descriptions["smtp_username"].Description,
				Optional:    true,
			},
			"smtp_password": schema.StringAttribute{
				Description: descriptions["smtp_password"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"smtp_auth_mechanism": schema.StringAttribute{
				Description:         descriptions["smtp_auth_mechanism"].Description,
				MarkdownDescription: descriptions["smtp_auth_mechanism"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedNotificationsSMTPAuthMech),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedNotifications{}.SMTPAuthMechOptions()...),
				},
			},

			// Sounds
			"console_bell": schema.BoolAttribute{
				Description: descriptions["console_bell"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"disable_beep": schema.BoolAttribute{
				Description: descriptions["disable_beep"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Telegram
			"telegram_enable": schema.BoolAttribute{
				Description: descriptions["telegram_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"telegram_api": schema.StringAttribute{
				Description: descriptions["telegram_api"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"telegram_chat_id": schema.StringAttribute{
				Description: descriptions["telegram_chat_id"].Description,
				Optional:    true,
			},

			// Pushover
			"pushover_enable": schema.BoolAttribute{
				Description: descriptions["pushover_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"pushover_api_key": schema.StringAttribute{
				Description: descriptions["pushover_api_key"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"pushover_user_key": schema.StringAttribute{
				Description: descriptions["pushover_user_key"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"pushover_sound": schema.StringAttribute{
				Description:         descriptions["pushover_sound"].Description,
				MarkdownDescription: descriptions["pushover_sound"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedNotificationsPushoverSound),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedNotifications{}.PushoverSoundOptions()...),
				},
			},
			"pushover_priority": schema.StringAttribute{
				Description:         descriptions["pushover_priority"].Description,
				MarkdownDescription: descriptions["pushover_priority"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedNotificationsPushoverPriority),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedNotifications{}.PushoverPriorityOptions()...),
				},
			},
			"pushover_retry": schema.Int64Attribute{
				Description: descriptions["pushover_retry"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedNotificationsPushoverRetry)),
				Validators: []validator.Int64{
					int64validator.AtLeast(30),
				},
			},
			"pushover_expire": schema.Int64Attribute{
				Description: descriptions["pushover_expire"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedNotificationsPushoverExpire)),
				Validators: []validator.Int64{
					int64validator.AtMost(10800),
				},
			},

			// Slack
			"slack_enable": schema.BoolAttribute{
				Description: descriptions["slack_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"slack_api": schema.StringAttribute{
				Description: descriptions["slack_api"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"slack_channel": schema.StringAttribute{
				Description: descriptions["slack_channel"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(80),
				},
			},

			// Apply
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *SystemAdvancedNotificationsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemAdvancedNotificationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemAdvancedNotificationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedNotifications
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedNotifications(ctx, aReq)
	if addError(&resp.Diagnostics, "Error creating system advanced notifications settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNotificationsChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced notifications changes", err)
	}
}

func (r *SystemAdvancedNotificationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemAdvancedNotificationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetAdvancedNotifications(ctx)
	if addError(&resp.Diagnostics, "Error reading system advanced notifications settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemAdvancedNotificationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemAdvancedNotificationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedNotifications
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedNotifications(ctx, aReq)
	if addError(&resp.Diagnostics, "Error updating system advanced notifications settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNotificationsChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced notifications changes", err)
	}
}

func (r *SystemAdvancedNotificationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemAdvancedNotificationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultNotif := pfsense.AdvancedNotifications{
		CertEnableNotify: true,
		CertExpireDays:   pfsense.DefaultAdvancedNotificationsCertExpireDays,
		SSLValidate:      true,
		SMTPAuthMech:     pfsense.DefaultAdvancedNotificationsSMTPAuthMech,
		ConsoleBell:      true,
		PushoverSound:    pfsense.DefaultAdvancedNotificationsPushoverSound,
		PushoverPriority: pfsense.DefaultAdvancedNotificationsPushoverPriority,
		PushoverRetry:    pfsense.DefaultAdvancedNotificationsPushoverRetry,
		PushoverExpire:   pfsense.DefaultAdvancedNotificationsPushoverExpire,
	}

	_, err := r.client.UpdateAdvancedNotifications(ctx, defaultNotif)
	if addError(&resp.Diagnostics, "Error resetting system advanced notifications settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedNotificationsChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced notifications changes", err)
	}
}

func (r *SystemAdvancedNotificationsResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	a, err := r.client.GetAdvancedNotifications(ctx)
	if addError(&resp.Diagnostics, "Error importing system advanced notifications settings", err) {
		return
	}

	var data SystemAdvancedNotificationsResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
