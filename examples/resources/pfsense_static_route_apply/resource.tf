# Create multiple static routes with apply deferred to a single operation
resource "pfsense_static_route" "example" {
  for_each = toset(["10.10.0.0/24", "10.10.4.0/24", "10.10.8.0/24"])

  network     = each.key
  gateway     = "WAN_GW"
  description = "Route for ${each.key}"
  apply       = false
}

# Apply all route changes at once
resource "pfsense_static_route_apply" "example" {
  lifecycle {
    replace_triggered_by = [
      pfsense_static_route.example,
    ]
  }
}
