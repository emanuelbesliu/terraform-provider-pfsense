data "pfsense_system_tunables" "all" {}

output "tunable_names" {
  value = [for t in data.pfsense_system_tunables.all.system_tunables : t.tunable]
}
