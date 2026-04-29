resource "pfsense_rest_api_settings" "example" {
  auth_methods = ["BasicAuth", "KeyAuth"]
}
