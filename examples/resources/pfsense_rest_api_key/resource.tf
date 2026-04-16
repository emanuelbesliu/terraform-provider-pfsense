resource "pfsense_rest_api_key" "terraform" {
  description  = "Terraform automation"
  hash_algo    = "sha256"
  length_bytes = 24
}

resource "pfsense_rest_api_key" "monitoring" {
  description = "Holmes monitoring agent"
}
