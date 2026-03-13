# Create multiple gateways with apply deferred to a single operation
resource "pfsense_gateway" "example" {
  for_each = {
    primary = "203.0.113.1"
    backup  = "198.51.100.1"
  }

  name       = "${upper(each.key)}_GW"
  interface  = "wan"
  ipprotocol = "inet"
  gateway    = each.value
  apply      = false
}

# Apply all gateway changes at once
resource "pfsense_gateway_apply" "example" {
  lifecycle {
    replace_triggered_by = [
      pfsense_gateway.example,
    ]
  }
}
