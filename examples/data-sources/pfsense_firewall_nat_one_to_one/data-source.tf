data "pfsense_firewall_nat_one_to_one" "webserver" {
  description = "Web server 1:1 NAT"
}

output "webserver_external" {
  value = data.pfsense_firewall_nat_one_to_one.webserver.external
}
