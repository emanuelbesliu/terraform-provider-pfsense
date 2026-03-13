# Look up a single gateway group by name
data "pfsense_gateway_group" "failover" {
  name = "WAN_FAILOVER"
}

output "failover_trigger" {
  value = data.pfsense_gateway_group.failover.trigger
}

output "failover_members" {
  value = data.pfsense_gateway_group.failover.members
}
