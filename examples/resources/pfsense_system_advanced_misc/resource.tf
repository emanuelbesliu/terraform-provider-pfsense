resource "pfsense_system_advanced_misc" "example" {
  # Power Savings
  powerd_enable       = true
  powerd_ac_mode      = "hadp"
  powerd_battery_mode = "adp"
  powerd_normal_mode  = "hadp"

  # Cryptographic & Thermal Hardware
  crypto_hardware  = "aesni"
  thermal_hardware = "coretemp"

  # Gateway Monitoring
  gw_down_kill_states = "down"

  # RAM Disk Settings
  rrd_backup_interval  = 8
  dhcp_backup_interval = 8
  logs_backup_interval = 8
}
