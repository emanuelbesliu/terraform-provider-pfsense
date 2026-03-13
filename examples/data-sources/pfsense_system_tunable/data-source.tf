# Look up a single system tunable by name
data "pfsense_system_tunable" "nmbclusters" {
  tunable = "kern.ipc.nmbclusters"
}

output "nmbclusters_value" {
  value = data.pfsense_system_tunable.nmbclusters.value
}
