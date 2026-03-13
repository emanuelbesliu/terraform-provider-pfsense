# Look up a single gateway by name
data "pfsense_gateway" "wan" {
  name = "WAN_GW"
}

output "wan_gateway_ip" {
  value = data.pfsense_gateway.wan.gateway
}

output "wan_is_default" {
  value = data.pfsense_gateway.wan.default_gw
}
