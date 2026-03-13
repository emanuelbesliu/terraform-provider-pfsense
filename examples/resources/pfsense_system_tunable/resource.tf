# Manage a system tunable to increase the number of mbuf clusters
resource "pfsense_system_tunable" "nmbclusters" {
  tunable     = "kern.ipc.nmbclusters"
  value       = "131072"
  description = "Increase mbuf clusters for high-traffic environments"
}

# Manage a system tunable with the default value
resource "pfsense_system_tunable" "syncookies" {
  tunable = "net.inet.tcp.syncookies"
  value   = "default"
}
