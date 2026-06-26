# Simple bridge joining two interfaces at layer 2
resource "pfsense_interface_bridge" "lan_bridge" {
  members     = ["lan", "opt1"]
  description = "LAN bridge"
}

# Bridge with Rapid Spanning Tree enabled
resource "pfsense_interface_bridge" "stp_bridge" {
  members        = ["opt1", "opt2"]
  description    = "Bridge with RSTP"
  enable_stp     = true
  protocol       = "rstp"
  priority       = 32768
  stp_interfaces = ["opt1", "opt2"]
}
