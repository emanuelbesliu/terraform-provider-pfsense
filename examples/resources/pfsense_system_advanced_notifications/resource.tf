resource "pfsense_system_advanced_notifications" "example" {
  # Certificate Expiration
  cert_enable_notify         = true
  revoked_cert_ignore_notify = false
  cert_expire_days           = 27

  # SMTP (Email) Notifications
  disable_smtp        = false
  smtp_ipaddress      = "mail.example.com"
  smtp_port           = 587
  smtp_timeout        = 20
  smtp_ssl            = true
  smtp_ssl_validate   = true
  smtp_from_address   = "pfsense@example.com"
  smtp_notify_address = "admin@example.com"
  smtp_username       = "pfsense@example.com"
  smtp_password       = "smtp-password"
  smtp_auth_mechanism = "PLAIN"

  # Sound
  console_bell = true
  disable_beep = false

  # Telegram
  telegram_enabled = false

  # Pushover
  pushover_enabled  = false
  pushover_sound    = "devicedefault"
  pushover_priority = "0"
  pushover_retry    = 60
  pushover_expire   = 300

  # Slack
  slack_enabled = false
}
