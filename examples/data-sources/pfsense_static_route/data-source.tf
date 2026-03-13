# Look up a single static route by network
data "pfsense_static_route" "internal" {
  network = "10.10.0.0/24"
}

output "internal_route_gateway" {
  value = data.pfsense_static_route.internal.gateway
}
