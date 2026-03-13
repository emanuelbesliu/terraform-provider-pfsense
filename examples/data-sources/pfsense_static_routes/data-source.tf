data "pfsense_static_routes" "all" {}

output "route_networks" {
  value = [for r in data.pfsense_static_routes.all.routes : r.network]
}
