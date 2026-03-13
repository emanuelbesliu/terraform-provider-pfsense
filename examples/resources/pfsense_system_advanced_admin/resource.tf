# Manage system advanced admin access settings
resource "pfsense_system_advanced_admin" "this" {
  webgui_protocol = "https"
  webgui_port     = 443
  max_processes   = 2

  roaming            = true
  login_autocomplete = true
  page_name_first    = false

  ssh_enabled           = true
  sshd_key_only         = "disabled"
  sshd_agent_forwarding = false
  ssh_port              = 22

  serial_speed    = 115200
  primary_console = "video"
}
