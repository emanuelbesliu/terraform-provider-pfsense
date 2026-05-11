resource "pfsense_wake_on_lan" "server" {
  interface   = "lan"
  mac         = "00:11:22:33:44:55"
  description = "File server"
}
