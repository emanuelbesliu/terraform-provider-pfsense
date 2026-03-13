package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultAdvancedNotificationsCertExpireDays   = 27
	DefaultAdvancedNotificationsPushoverRetry    = 60
	DefaultAdvancedNotificationsPushoverExpire   = 300
	DefaultAdvancedNotificationsSMTPAuthMech     = "PLAIN"
	DefaultAdvancedNotificationsPushoverPriority = "0"
	DefaultAdvancedNotificationsPushoverSound    = "devicedefault"
)

// AdvancedNotifications represents the System > Advanced > Notifications configuration.
type AdvancedNotifications struct {
	// General Settings - Certificate Expiration
	CertEnableNotify        bool // enable daily cert expiration notifications
	RevokedCertIgnoreNotify bool // ignore notifications for revoked certificates
	CertExpireDays          int  // days threshold; 0 = use default (27)

	// E-Mail (SMTP)
	DisableSMTP         bool
	SMTPIPAddress       string // FQDN or IP
	SMTPPort            int    // 0 = not set
	SMTPTimeout         int    // 0 = default (20s)
	SMTPSSL             bool   // enable SSL/TLS
	SSLValidate         bool   // validate SSL/TLS certificate
	SMTPFromAddress     string
	SMTPNotifyEmailAddr string
	SMTPUsername        string
	SMTPPassword        string // sensitive
	SMTPAuthMech        string // PLAIN or LOGIN

	// Sounds
	ConsoleBell bool // enable console bell
	DisableBeep bool // disable startup/shutdown beep

	// Telegram
	TelegramEnable bool
	TelegramAPI    string // sensitive
	TelegramChatID string

	// Pushover
	PushoverEnable   bool
	PushoverAPIKey   string // sensitive
	PushoverUserKey  string // sensitive
	PushoverSound    string
	PushoverPriority string // "-2" to "2"
	PushoverRetry    int    // seconds, min 30, default 60
	PushoverExpire   int    // seconds, max 10800, default 300

	// Slack
	SlackEnable  bool
	SlackAPI     string // sensitive
	SlackChannel string
}

func (a *AdvancedNotifications) SetSMTPAuthMech(mech string) error {
	valid := AdvancedNotifications{}.SMTPAuthMechOptions()
	for _, v := range valid {
		if mech == v {
			a.SMTPAuthMech = mech

			return nil
		}
	}

	return fmt.Errorf("%w, SMTP auth mechanism must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedNotifications) SMTPAuthMechOptions() []string {
	return []string{"PLAIN", "LOGIN"}
}

func (a *AdvancedNotifications) SetPushoverSound(sound string) error {
	valid := AdvancedNotifications{}.PushoverSoundOptions()
	for _, v := range valid {
		if sound == v {
			a.PushoverSound = sound

			return nil
		}
	}

	return fmt.Errorf("%w, Pushover sound must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedNotifications) PushoverSoundOptions() []string {
	return []string{
		"devicedefault", "pushover", "bike", "bugle", "cashregister",
		"classical", "cosmic", "falling", "gamelan", "incoming",
		"intermission", "magic", "mechanical", "pianobar", "siren",
		"spacealarm", "tugboat", "alien", "climb", "persistent",
		"echo", "updown", "vibrate", "none",
	}
}

func (a *AdvancedNotifications) SetPushoverPriority(priority string) error {
	valid := AdvancedNotifications{}.PushoverPriorityOptions()
	for _, v := range valid {
		if priority == v {
			a.PushoverPriority = priority

			return nil
		}
	}

	return fmt.Errorf("%w, Pushover priority must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedNotifications) PushoverPriorityOptions() []string {
	return []string{"-2", "-1", "0", "1", "2"}
}

// advancedNotificationsResponse is the JSON shape returned by the PHP read command.
type advancedNotificationsResponse struct {
	// General Settings
	CertExpireEnable        json.RawMessage `json:"cert_expire_enable"`
	CertExpireIgnoreRevoked json.RawMessage `json:"cert_expire_ignore_revoked"`
	CertExpireDays          json.RawMessage `json:"cert_expire_days"`

	// SMTP
	SMTPDisable         json.RawMessage `json:"smtp_disable"`
	SMTPIPAddress       json.RawMessage `json:"smtp_ipaddress"`
	SMTPPort            json.RawMessage `json:"smtp_port"`
	SMTPTimeout         json.RawMessage `json:"smtp_timeout"`
	SMTPSSL             json.RawMessage `json:"smtp_ssl"`
	SMTPSSLValidate     json.RawMessage `json:"smtp_sslvalidate"`
	SMTPFromAddress     json.RawMessage `json:"smtp_fromaddress"`
	SMTPNotifyEmailAddr json.RawMessage `json:"smtp_notifyemailaddress"`
	SMTPUsername        json.RawMessage `json:"smtp_username"`
	SMTPPassword        json.RawMessage `json:"smtp_password"`
	SMTPAuthMech        json.RawMessage `json:"smtp_authentication_mechanism"`

	// Sounds
	ConsoleBell json.RawMessage `json:"consolebell"`
	DisableBeep json.RawMessage `json:"disablebeep"`

	// Telegram
	TelegramEnabled json.RawMessage `json:"telegram_enabled"`
	TelegramAPI     json.RawMessage `json:"telegram_api"`
	TelegramChatID  json.RawMessage `json:"telegram_chatid"`

	// Pushover
	PushoverEnabled  json.RawMessage `json:"pushover_enabled"`
	PushoverAPIKey   json.RawMessage `json:"pushover_apikey"`
	PushoverUserKey  json.RawMessage `json:"pushover_userkey"`
	PushoverSound    json.RawMessage `json:"pushover_sound"`
	PushoverPriority json.RawMessage `json:"pushover_priority"`
	PushoverRetry    json.RawMessage `json:"pushover_retry"`
	PushoverExpire   json.RawMessage `json:"pushover_expire"`

	// Slack
	SlackEnabled json.RawMessage `json:"slack_enabled"`
	SlackAPI     json.RawMessage `json:"slack_api"`
	SlackChannel json.RawMessage `json:"slack_channel"`
}

func parseAdvancedNotificationsResponse(resp advancedNotificationsResponse) (AdvancedNotifications, error) {
	var a AdvancedNotifications

	// General Settings - Certificate Expiration
	// cert_expire_enable: config stores "enabled"/"disabled", default is "enabled" (not disabled)
	certExpireEnable := rawToString(resp.CertExpireEnable)
	a.CertEnableNotify = (certExpireEnable != "disabled")

	certExpireIgnoreRevoked := rawToString(resp.CertExpireIgnoreRevoked)
	a.RevokedCertIgnoreNotify = (certExpireIgnoreRevoked == "enabled")

	a.CertExpireDays = rawToInt(resp.CertExpireDays)
	if a.CertExpireDays == 0 {
		a.CertExpireDays = DefaultAdvancedNotificationsCertExpireDays
	}

	// SMTP
	a.DisableSMTP = rawIsPresent(resp.SMTPDisable)
	a.SMTPIPAddress = rawToString(resp.SMTPIPAddress)
	a.SMTPPort = rawToInt(resp.SMTPPort)
	a.SMTPTimeout = rawToInt(resp.SMTPTimeout)
	a.SMTPSSL = rawIsPresent(resp.SMTPSSL)

	// sslvalidate: config stores "enabled"/"disabled", default is "enabled" (not disabled)
	sslValidate := rawToString(resp.SMTPSSLValidate)
	a.SSLValidate = (sslValidate != "disabled")

	a.SMTPFromAddress = rawToString(resp.SMTPFromAddress)
	a.SMTPNotifyEmailAddr = rawToString(resp.SMTPNotifyEmailAddr)
	a.SMTPUsername = rawToString(resp.SMTPUsername)
	a.SMTPPassword = rawToString(resp.SMTPPassword)

	authMech := rawToString(resp.SMTPAuthMech)
	if authMech == "" {
		authMech = DefaultAdvancedNotificationsSMTPAuthMech
	}

	if err := a.SetSMTPAuthMech(authMech); err != nil {
		return a, err
	}

	// Sounds
	// consolebell: config stores "enabled"/"disabled", default is "enabled"
	consoleBell := rawToString(resp.ConsoleBell)
	if consoleBell == "" {
		consoleBell = "enabled"
	}

	a.ConsoleBell = (consoleBell == "enabled")
	a.DisableBeep = rawIsPresent(resp.DisableBeep)

	// Telegram
	telegramEnabled := rawToString(resp.TelegramEnabled)
	a.TelegramEnable = (telegramEnabled == "1" || telegramEnabled == "true")
	a.TelegramAPI = rawToString(resp.TelegramAPI)
	a.TelegramChatID = rawToString(resp.TelegramChatID)

	// Pushover
	pushoverEnabled := rawToString(resp.PushoverEnabled)
	a.PushoverEnable = (pushoverEnabled == "1" || pushoverEnabled == "true")
	a.PushoverAPIKey = rawToString(resp.PushoverAPIKey)
	a.PushoverUserKey = rawToString(resp.PushoverUserKey)

	pushoverSound := rawToString(resp.PushoverSound)
	if pushoverSound == "" {
		pushoverSound = DefaultAdvancedNotificationsPushoverSound
	}

	if err := a.SetPushoverSound(pushoverSound); err != nil {
		return a, err
	}

	pushoverPriority := rawToString(resp.PushoverPriority)
	if pushoverPriority == "" {
		pushoverPriority = DefaultAdvancedNotificationsPushoverPriority
	}

	if err := a.SetPushoverPriority(pushoverPriority); err != nil {
		return a, err
	}

	a.PushoverRetry = rawToInt(resp.PushoverRetry)
	if a.PushoverRetry == 0 {
		a.PushoverRetry = DefaultAdvancedNotificationsPushoverRetry
	}

	a.PushoverExpire = rawToInt(resp.PushoverExpire)
	if a.PushoverExpire == 0 {
		a.PushoverExpire = DefaultAdvancedNotificationsPushoverExpire
	}

	// Slack
	slackEnabled := rawToString(resp.SlackEnabled)
	a.SlackEnable = (slackEnabled == "1" || slackEnabled == "true")
	a.SlackAPI = rawToString(resp.SlackAPI)
	a.SlackChannel = rawToString(resp.SlackChannel)

	return a, nil
}

func (pf *Client) getAdvancedNotifications(ctx context.Context) (*AdvancedNotifications, error) {
	command := "$notif = config_get_path('notifications', array());" +
		"$sys = config_get_path('system', array());" +
		"$out = array(" +
		// General Settings
		"'cert_expire_enable' => isset($notif['certexpire']['enable']) ? $notif['certexpire']['enable'] : null," +
		"'cert_expire_ignore_revoked' => isset($notif['certexpire']['ignore_revoked']) ? $notif['certexpire']['ignore_revoked'] : null," +
		"'cert_expire_days' => isset($notif['certexpire']['expiredays']) ? $notif['certexpire']['expiredays'] : null," +
		// SMTP
		"'smtp_disable' => isset($notif['smtp']['disable']) ? $notif['smtp']['disable'] : null," +
		"'smtp_ipaddress' => isset($notif['smtp']['ipaddress']) ? $notif['smtp']['ipaddress'] : null," +
		"'smtp_port' => isset($notif['smtp']['port']) ? $notif['smtp']['port'] : null," +
		"'smtp_timeout' => isset($notif['smtp']['timeout']) ? $notif['smtp']['timeout'] : null," +
		"'smtp_ssl' => isset($notif['smtp']['ssl']) ? $notif['smtp']['ssl'] : null," +
		"'smtp_sslvalidate' => isset($notif['smtp']['sslvalidate']) ? $notif['smtp']['sslvalidate'] : null," +
		"'smtp_fromaddress' => isset($notif['smtp']['fromaddress']) ? $notif['smtp']['fromaddress'] : null," +
		"'smtp_notifyemailaddress' => isset($notif['smtp']['notifyemailaddress']) ? $notif['smtp']['notifyemailaddress'] : null," +
		"'smtp_username' => isset($notif['smtp']['username']) ? $notif['smtp']['username'] : null," +
		"'smtp_password' => isset($notif['smtp']['password']) ? $notif['smtp']['password'] : null," +
		"'smtp_authentication_mechanism' => isset($notif['smtp']['authentication_mechanism']) ? $notif['smtp']['authentication_mechanism'] : null," +
		// Sounds
		"'consolebell' => isset($sys['consolebell']) ? $sys['consolebell'] : null," +
		"'disablebeep' => isset($sys['disablebeep']) ? $sys['disablebeep'] : null," +
		// Telegram
		"'telegram_enabled' => isset($notif['telegram']['enabled']) ? $notif['telegram']['enabled'] : null," +
		"'telegram_api' => isset($notif['telegram']['api']) ? $notif['telegram']['api'] : null," +
		"'telegram_chatid' => isset($notif['telegram']['chatid']) ? $notif['telegram']['chatid'] : null," +
		// Pushover
		"'pushover_enabled' => isset($notif['pushover']['enabled']) ? $notif['pushover']['enabled'] : null," +
		"'pushover_apikey' => isset($notif['pushover']['apikey']) ? $notif['pushover']['apikey'] : null," +
		"'pushover_userkey' => isset($notif['pushover']['userkey']) ? $notif['pushover']['userkey'] : null," +
		"'pushover_sound' => isset($notif['pushover']['sound']) ? $notif['pushover']['sound'] : null," +
		"'pushover_priority' => isset($notif['pushover']['priority']) ? $notif['pushover']['priority'] : null," +
		"'pushover_retry' => isset($notif['pushover']['retry']) ? $notif['pushover']['retry'] : null," +
		"'pushover_expire' => isset($notif['pushover']['expire']) ? $notif['pushover']['expire'] : null," +
		// Slack
		"'slack_enabled' => isset($notif['slack']['enabled']) ? $notif['slack']['enabled'] : null," +
		"'slack_api' => isset($notif['slack']['api']) ? $notif['slack']['api'] : null," +
		"'slack_channel' => isset($notif['slack']['channel']) ? $notif['slack']['channel'] : null" +
		");" +
		"print(json_encode($out));"

	var resp advancedNotificationsResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	a, err := parseAdvancedNotificationsResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w advanced notifications response, %w", ErrUnableToParse, err)
	}

	return &a, nil
}

func (pf *Client) GetAdvancedNotifications(ctx context.Context) (*AdvancedNotifications, error) {
	defer pf.read(&pf.mutexes.AdvancedNotifications)()

	a, err := pf.getAdvancedNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced notifications, %w", ErrGetOperationFailed, err)
	}

	return a, nil
}

func advancedNotificationsFormValues(a AdvancedNotifications) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// General Settings - Certificate Expiration
	if a.CertEnableNotify {
		values.Set("cert_enable_notify", "yes")
	}

	if a.RevokedCertIgnoreNotify {
		values.Set("revoked_cert_ignore_notify", "yes")
	}

	if a.CertExpireDays > 0 {
		values.Set("certexpiredays", strconv.Itoa(a.CertExpireDays))
	}

	// SMTP
	if a.DisableSMTP {
		values.Set("disable_smtp", "yes")
	}

	values.Set("smtpipaddress", a.SMTPIPAddress)

	if a.SMTPPort > 0 {
		values.Set("smtpport", strconv.Itoa(a.SMTPPort))
	}

	if a.SMTPTimeout > 0 {
		values.Set("smtptimeout", strconv.Itoa(a.SMTPTimeout))
	}

	if a.SMTPSSL {
		values.Set("smtpssl", "yes")
	}

	if a.SSLValidate {
		values.Set("sslvalidate", "yes")
	}

	values.Set("smtpfromaddress", a.SMTPFromAddress)
	values.Set("smtpnotifyemailaddress", a.SMTPNotifyEmailAddr)
	values.Set("smtpusername", a.SMTPUsername)
	values.Set("smtppassword", a.SMTPPassword)
	values.Set("smtppassword_confirm", a.SMTPPassword)
	values.Set("smtpauthmech", a.SMTPAuthMech)

	// Sounds
	if a.ConsoleBell {
		values.Set("consolebell", "yes")
	}

	if a.DisableBeep {
		values.Set("disablebeep", "yes")
	}

	// Telegram
	if a.TelegramEnable {
		values.Set("enable_telegram", "yes")
	}

	values.Set("api", a.TelegramAPI)
	values.Set("chatid", a.TelegramChatID)

	// Pushover
	if a.PushoverEnable {
		values.Set("enable_pushover", "yes")
	}

	values.Set("pushoverapikey", a.PushoverAPIKey)
	values.Set("pushoveruserkey", a.PushoverUserKey)
	values.Set("pushoversound", a.PushoverSound)
	values.Set("pushoverpriority", a.PushoverPriority)
	values.Set("pushoverretry", strconv.Itoa(a.PushoverRetry))
	values.Set("pushoverexpire", strconv.Itoa(a.PushoverExpire))

	// Slack
	if a.SlackEnable {
		values.Set("enable_slack", "yes")
	}

	values.Set("slack_api", a.SlackAPI)
	values.Set("slack_channel", a.SlackChannel)

	return values
}

func (pf *Client) UpdateAdvancedNotifications(ctx context.Context, a AdvancedNotifications) (*AdvancedNotifications, error) {
	defer pf.write(&pf.mutexes.AdvancedNotifications)()

	relativeURL := url.URL{Path: "system_advanced_notifications.php"}
	values := advancedNotificationsFormValues(a)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w advanced notifications, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w advanced notifications, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getAdvancedNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced notifications after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

// ApplyAdvancedNotificationsChanges is a no-op since notification settings
// take effect immediately after save (no service reload needed).
// We provide this method to maintain the consistent singleton pattern.
func (pf *Client) ApplyAdvancedNotificationsChanges(ctx context.Context) error {
	pf.mutexes.AdvancedNotificationsApply.Lock()
	defer pf.mutexes.AdvancedNotificationsApply.Unlock()

	// Notifications config takes effect immediately; no apply step needed.
	// The saveAdvancedNotifications() PHP function just calls write_config().
	_ = ctx

	return nil
}
