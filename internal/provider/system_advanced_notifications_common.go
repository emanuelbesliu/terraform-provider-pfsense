package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// SystemAdvancedNotificationsModel represents the Terraform model for system advanced notifications settings.
type SystemAdvancedNotificationsModel struct {
	// General Settings - Certificate Expiration
	CertEnableNotify        types.Bool  `tfsdk:"cert_enable_notify"`
	RevokedCertIgnoreNotify types.Bool  `tfsdk:"revoked_cert_ignore_notify"`
	CertExpireDays          types.Int64 `tfsdk:"cert_expire_days"`

	// E-Mail (SMTP)
	DisableSMTP         types.Bool   `tfsdk:"disable_smtp"`
	SMTPIPAddress       types.String `tfsdk:"smtp_ip_address"`
	SMTPPort            types.Int64  `tfsdk:"smtp_port"`
	SMTPTimeout         types.Int64  `tfsdk:"smtp_timeout"`
	SMTPSSL             types.Bool   `tfsdk:"smtp_ssl"`
	SSLValidate         types.Bool   `tfsdk:"ssl_validate"`
	SMTPFromAddress     types.String `tfsdk:"smtp_from_address"`
	SMTPNotifyEmailAddr types.String `tfsdk:"smtp_notify_email_address"`
	SMTPUsername        types.String `tfsdk:"smtp_username"`
	SMTPPassword        types.String `tfsdk:"smtp_password"`
	SMTPAuthMech        types.String `tfsdk:"smtp_auth_mechanism"`

	// Sounds
	ConsoleBell types.Bool `tfsdk:"console_bell"`
	DisableBeep types.Bool `tfsdk:"disable_beep"`

	// Telegram
	TelegramEnable types.Bool   `tfsdk:"telegram_enable"`
	TelegramAPI    types.String `tfsdk:"telegram_api"`
	TelegramChatID types.String `tfsdk:"telegram_chat_id"`

	// Pushover
	PushoverEnable   types.Bool   `tfsdk:"pushover_enable"`
	PushoverAPIKey   types.String `tfsdk:"pushover_api_key"`
	PushoverUserKey  types.String `tfsdk:"pushover_user_key"`
	PushoverSound    types.String `tfsdk:"pushover_sound"`
	PushoverPriority types.String `tfsdk:"pushover_priority"`
	PushoverRetry    types.Int64  `tfsdk:"pushover_retry"`
	PushoverExpire   types.Int64  `tfsdk:"pushover_expire"`

	// Slack
	SlackEnable  types.Bool   `tfsdk:"slack_enable"`
	SlackAPI     types.String `tfsdk:"slack_api"`
	SlackChannel types.String `tfsdk:"slack_channel"`
}

func (SystemAdvancedNotificationsModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"cert_enable_notify": {
			Description: "Enable daily notifications of expired and soon-to-expire certificates. When enabled, the firewall checks CA and certificate expiration times daily.",
		},
		"revoked_cert_ignore_notify": {
			Description: "Ignore notifications for revoked certificates. When enabled, the firewall will not check expiration for certificates that have been revoked.",
		},
		"cert_expire_days": {
			Description: fmt.Sprintf("Number of days before certificate expiration to start generating notifications. Default is %d days.", pfsense.DefaultAdvancedNotificationsCertExpireDays),
		},
		"disable_smtp": {
			Description: "Disable SMTP notifications but preserve the settings. Some other mechanisms, such as packages, may need the SMTP settings in place to function.",
		},
		"smtp_ip_address": {
			Description: "FQDN or IP address of the SMTP E-Mail server to which notifications will be sent.",
		},
		"smtp_port": {
			Description: "Port of the SMTP E-Mail server. Typically 25, 587 (submission), or 465 (smtps).",
		},
		"smtp_timeout": {
			Description: "Connection timeout in seconds for the SMTP server. Default is 20 seconds.",
		},
		"smtp_ssl": {
			Description: "Enable SMTP over SSL/TLS for secure email delivery.",
		},
		"ssl_validate": {
			Description: "Validate the SSL/TLS certificate presented by the SMTP server. When disabled, the server certificate will not be validated but encryption will still be used if available.",
		},
		"smtp_from_address": {
			Description: "E-mail address that will appear in the From field of notification emails.",
		},
		"smtp_notify_email_address": {
			Description: "E-mail address to send notification emails to.",
		},
		"smtp_username": {
			Description: "Username for SMTP authentication.",
		},
		"smtp_password": {
			Description: "Password for SMTP authentication.",
		},
		"smtp_auth_mechanism": {
			Description:         fmt.Sprintf("Authentication mechanism used by the SMTP server. Options: %s. Defaults to '%s'. Most servers work with PLAIN; some servers like Exchange or Office365 might require LOGIN.", wrapElementsJoin(pfsense.AdvancedNotifications{}.SMTPAuthMechOptions(), "'"), pfsense.DefaultAdvancedNotificationsSMTPAuthMech),
			MarkdownDescription: fmt.Sprintf("Authentication mechanism used by the SMTP server. Options: %s. Defaults to `%s`. Most servers work with PLAIN; some servers like Exchange or Office365 might require LOGIN.", wrapElementsJoin(pfsense.AdvancedNotifications{}.SMTPAuthMechOptions(), "`"), pfsense.DefaultAdvancedNotificationsSMTPAuthMech),
		},
		"console_bell": {
			Description: "Enable the console bell. When enabled, emergency console messages will trigger a bell in connected consoles including serial terminals.",
		},
		"disable_beep": {
			Description: "Disable the startup/shutdown beep sound through the built-in PC speaker.",
		},
		"telegram_enable": {
			Description: "Enable Telegram notifications. Requires a Telegram Bot API key and Chat ID.",
		},
		"telegram_api": {
			Description: "Telegram Bot API key required to authenticate with the Telegram API server.",
		},
		"telegram_chat_id": {
			Description: "Telegram chat ID (private) or channel @username (public) for sending notifications.",
		},
		"pushover_enable": {
			Description: "Enable Pushover notifications. Requires a Pushover API key and user key.",
		},
		"pushover_api_key": {
			Description: "Pushover application API key for authentication.",
		},
		"pushover_user_key": {
			Description: "Pushover account user key.",
		},
		"pushover_sound": {
			Description:         fmt.Sprintf("Pushover notification sound. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedNotifications{}.PushoverSoundOptions(), "'"), pfsense.DefaultAdvancedNotificationsPushoverSound),
			MarkdownDescription: fmt.Sprintf("Pushover notification sound. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedNotifications{}.PushoverSoundOptions(), "`"), pfsense.DefaultAdvancedNotificationsPushoverSound),
		},
		"pushover_priority": {
			Description:         fmt.Sprintf("Pushover message priority. Options: %s. Defaults to '%s' (Normal Priority).", wrapElementsJoin(pfsense.AdvancedNotifications{}.PushoverPriorityOptions(), "'"), pfsense.DefaultAdvancedNotificationsPushoverPriority),
			MarkdownDescription: fmt.Sprintf("Pushover message priority. Options: %s. Defaults to `%s` (Normal Priority).", wrapElementsJoin(pfsense.AdvancedNotifications{}.PushoverPriorityOptions(), "`"), pfsense.DefaultAdvancedNotificationsPushoverPriority),
		},
		"pushover_retry": {
			Description: fmt.Sprintf("Pushover emergency priority notification retry interval in seconds. Minimum 30 seconds. Default is %d seconds.", pfsense.DefaultAdvancedNotificationsPushoverRetry),
		},
		"pushover_expire": {
			Description: fmt.Sprintf("Pushover emergency priority notification expiration time in seconds. Maximum 10800 seconds (3 hours). Default is %d seconds.", pfsense.DefaultAdvancedNotificationsPushoverExpire),
		},
		"slack_enable": {
			Description: "Enable Slack notifications. Requires a Slack API key and channel name.",
		},
		"slack_api": {
			Description: "Slack API key for authentication.",
		},
		"slack_channel": {
			Description: "Slack channel name to send notifications to. May only contain lowercase letters, numbers, hyphens, and underscores. Maximum 80 characters.",
		},
	}
}

func (SystemAdvancedNotificationsModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cert_enable_notify":         types.BoolType,
		"revoked_cert_ignore_notify": types.BoolType,
		"cert_expire_days":           types.Int64Type,
		"disable_smtp":               types.BoolType,
		"smtp_ip_address":            types.StringType,
		"smtp_port":                  types.Int64Type,
		"smtp_timeout":               types.Int64Type,
		"smtp_ssl":                   types.BoolType,
		"ssl_validate":               types.BoolType,
		"smtp_from_address":          types.StringType,
		"smtp_notify_email_address":  types.StringType,
		"smtp_username":              types.StringType,
		"smtp_password":              types.StringType,
		"smtp_auth_mechanism":        types.StringType,
		"console_bell":               types.BoolType,
		"disable_beep":               types.BoolType,
		"telegram_enable":            types.BoolType,
		"telegram_api":               types.StringType,
		"telegram_chat_id":           types.StringType,
		"pushover_enable":            types.BoolType,
		"pushover_api_key":           types.StringType,
		"pushover_user_key":          types.StringType,
		"pushover_sound":             types.StringType,
		"pushover_priority":          types.StringType,
		"pushover_retry":             types.Int64Type,
		"pushover_expire":            types.Int64Type,
		"slack_enable":               types.BoolType,
		"slack_api":                  types.StringType,
		"slack_channel":              types.StringType,
	}
}

func (m *SystemAdvancedNotificationsModel) Set(_ context.Context, a pfsense.AdvancedNotifications) diag.Diagnostics {
	var diags diag.Diagnostics

	// General Settings - Certificate Expiration
	m.CertEnableNotify = types.BoolValue(a.CertEnableNotify)
	m.RevokedCertIgnoreNotify = types.BoolValue(a.RevokedCertIgnoreNotify)
	m.CertExpireDays = types.Int64Value(int64(a.CertExpireDays))

	// SMTP
	m.DisableSMTP = types.BoolValue(a.DisableSMTP)

	if a.SMTPIPAddress != "" {
		m.SMTPIPAddress = types.StringValue(a.SMTPIPAddress)
	} else {
		m.SMTPIPAddress = types.StringNull()
	}

	m.SMTPPort = types.Int64Value(int64(a.SMTPPort))
	m.SMTPTimeout = types.Int64Value(int64(a.SMTPTimeout))
	m.SMTPSSL = types.BoolValue(a.SMTPSSL)
	m.SSLValidate = types.BoolValue(a.SSLValidate)

	if a.SMTPFromAddress != "" {
		m.SMTPFromAddress = types.StringValue(a.SMTPFromAddress)
	} else {
		m.SMTPFromAddress = types.StringNull()
	}

	if a.SMTPNotifyEmailAddr != "" {
		m.SMTPNotifyEmailAddr = types.StringValue(a.SMTPNotifyEmailAddr)
	} else {
		m.SMTPNotifyEmailAddr = types.StringNull()
	}

	if a.SMTPUsername != "" {
		m.SMTPUsername = types.StringValue(a.SMTPUsername)
	} else {
		m.SMTPUsername = types.StringNull()
	}

	if a.SMTPPassword != "" {
		m.SMTPPassword = types.StringValue(a.SMTPPassword)
	} else {
		m.SMTPPassword = types.StringNull()
	}

	m.SMTPAuthMech = types.StringValue(a.SMTPAuthMech)

	// Sounds
	m.ConsoleBell = types.BoolValue(a.ConsoleBell)
	m.DisableBeep = types.BoolValue(a.DisableBeep)

	// Telegram
	m.TelegramEnable = types.BoolValue(a.TelegramEnable)

	if a.TelegramAPI != "" {
		m.TelegramAPI = types.StringValue(a.TelegramAPI)
	} else {
		m.TelegramAPI = types.StringNull()
	}

	if a.TelegramChatID != "" {
		m.TelegramChatID = types.StringValue(a.TelegramChatID)
	} else {
		m.TelegramChatID = types.StringNull()
	}

	// Pushover
	m.PushoverEnable = types.BoolValue(a.PushoverEnable)

	if a.PushoverAPIKey != "" {
		m.PushoverAPIKey = types.StringValue(a.PushoverAPIKey)
	} else {
		m.PushoverAPIKey = types.StringNull()
	}

	if a.PushoverUserKey != "" {
		m.PushoverUserKey = types.StringValue(a.PushoverUserKey)
	} else {
		m.PushoverUserKey = types.StringNull()
	}

	m.PushoverSound = types.StringValue(a.PushoverSound)
	m.PushoverPriority = types.StringValue(a.PushoverPriority)
	m.PushoverRetry = types.Int64Value(int64(a.PushoverRetry))
	m.PushoverExpire = types.Int64Value(int64(a.PushoverExpire))

	// Slack
	m.SlackEnable = types.BoolValue(a.SlackEnable)

	if a.SlackAPI != "" {
		m.SlackAPI = types.StringValue(a.SlackAPI)
	} else {
		m.SlackAPI = types.StringNull()
	}

	if a.SlackChannel != "" {
		m.SlackChannel = types.StringValue(a.SlackChannel)
	} else {
		m.SlackChannel = types.StringNull()
	}

	return diags
}

func (m SystemAdvancedNotificationsModel) Value(_ context.Context, a *pfsense.AdvancedNotifications) diag.Diagnostics {
	var diags diag.Diagnostics

	// General Settings - Certificate Expiration
	a.CertEnableNotify = m.CertEnableNotify.ValueBool()
	a.RevokedCertIgnoreNotify = m.RevokedCertIgnoreNotify.ValueBool()
	a.CertExpireDays = int(m.CertExpireDays.ValueInt64())

	// SMTP
	a.DisableSMTP = m.DisableSMTP.ValueBool()

	if !m.SMTPIPAddress.IsNull() {
		a.SMTPIPAddress = m.SMTPIPAddress.ValueString()
	}

	a.SMTPPort = int(m.SMTPPort.ValueInt64())
	a.SMTPTimeout = int(m.SMTPTimeout.ValueInt64())
	a.SMTPSSL = m.SMTPSSL.ValueBool()
	a.SSLValidate = m.SSLValidate.ValueBool()

	if !m.SMTPFromAddress.IsNull() {
		a.SMTPFromAddress = m.SMTPFromAddress.ValueString()
	}

	if !m.SMTPNotifyEmailAddr.IsNull() {
		a.SMTPNotifyEmailAddr = m.SMTPNotifyEmailAddr.ValueString()
	}

	if !m.SMTPUsername.IsNull() {
		a.SMTPUsername = m.SMTPUsername.ValueString()
	}

	if !m.SMTPPassword.IsNull() {
		a.SMTPPassword = m.SMTPPassword.ValueString()
	}

	addPathError(
		&diags,
		path.Root("smtp_auth_mechanism"),
		"SMTP auth mechanism cannot be parsed",
		a.SetSMTPAuthMech(m.SMTPAuthMech.ValueString()),
	)

	// Sounds
	a.ConsoleBell = m.ConsoleBell.ValueBool()
	a.DisableBeep = m.DisableBeep.ValueBool()

	// Telegram
	a.TelegramEnable = m.TelegramEnable.ValueBool()

	if !m.TelegramAPI.IsNull() {
		a.TelegramAPI = m.TelegramAPI.ValueString()
	}

	if !m.TelegramChatID.IsNull() {
		a.TelegramChatID = m.TelegramChatID.ValueString()
	}

	// Pushover
	a.PushoverEnable = m.PushoverEnable.ValueBool()

	if !m.PushoverAPIKey.IsNull() {
		a.PushoverAPIKey = m.PushoverAPIKey.ValueString()
	}

	if !m.PushoverUserKey.IsNull() {
		a.PushoverUserKey = m.PushoverUserKey.ValueString()
	}

	addPathError(
		&diags,
		path.Root("pushover_sound"),
		"Pushover sound cannot be parsed",
		a.SetPushoverSound(m.PushoverSound.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("pushover_priority"),
		"Pushover priority cannot be parsed",
		a.SetPushoverPriority(m.PushoverPriority.ValueString()),
	)

	a.PushoverRetry = int(m.PushoverRetry.ValueInt64())
	a.PushoverExpire = int(m.PushoverExpire.ValueInt64())

	// Slack
	a.SlackEnable = m.SlackEnable.ValueBool()

	if !m.SlackAPI.IsNull() {
		a.SlackAPI = m.SlackAPI.ValueString()
	}

	if !m.SlackChannel.IsNull() {
		a.SlackChannel = m.SlackChannel.ValueString()
	}

	return diags
}
