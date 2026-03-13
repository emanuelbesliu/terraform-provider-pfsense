resource "pfsense_interface_group" "example" {
  name        = "SERVERS"
  members     = ["lan", "opt1"]
  description = "Server interfaces"
}
