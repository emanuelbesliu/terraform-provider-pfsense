# Failover gateway group - traffic goes to WAN_GW, falls back to BACKUP_GW
resource "pfsense_gateway_group" "failover" {
  name    = "WAN_FAILOVER"
  trigger = "down"

  members = [
    {
      gateway = "WAN_GW"
      tier    = 1
    },
    {
      gateway = "BACKUP_GW"
      tier    = 2
    }
  ]
}

# Load balancing gateway group - traffic is distributed across both gateways
resource "pfsense_gateway_group" "loadbalance" {
  name        = "WAN_LOADBALANCE"
  trigger     = "downloss"
  description = "Load balance across WAN links"

  members = [
    {
      gateway = "WAN_GW"
      tier    = 1
    },
    {
      gateway = "BACKUP_GW"
      tier    = 1
    }
  ]
}

# Gateway group with failover state handling and virtual IPs
resource "pfsense_gateway_group" "advanced" {
  name                 = "WAN_ADVANCED"
  trigger              = "downlosslatency"
  description          = "Advanced failover with state handling"
  keep_failover_states = "kill"

  members = [
    {
      gateway    = "WAN_GW"
      tier       = 1
      virtual_ip = "203.0.113.10"
    },
    {
      gateway    = "BACKUP_GW"
      tier       = 2
      virtual_ip = "198.51.100.10"
    }
  ]
}
