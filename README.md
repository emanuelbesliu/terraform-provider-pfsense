# Terraform Provider for pfSense

Configure [pfSense](https://www.pfsense.org/) firewalls with Terraform. Manages firewall rules, aliases, gateways, DNS resolver, DHCP, IPsec VPN, interfaces, VLANs, system settings, users, and more through the pfSense web UI.

> [!WARNING]
> All versions released prior to `v1.0.0` are to be considered [breaking changes](https://semver.org/#how-do-i-know-when-to-release-100).

## How It Works

This provider interacts with pfSense by **scraping the web UI** -- it authenticates via the login page, maintains a session with CSRF tokens, and submits forms or executes PHP commands just like a user would through the browser. No REST API package or additional software is required on the pfSense box.

## Compatibility

| Provider | pfSense            | Terraform |
| :------: | :----------------: | :-------: |
| < v1.0.0 | >= 2.6.0, <= 2.7.2 | >= 1.7.0  |

Validated with pfSense CE. Compatibility with pfSense Plus is not guaranteed.

## Quick Start

```hcl
terraform {
  required_providers {
    pfsense = {
      source  = "emanuelbesliu/pfsense"
      version = "~> 0.43"
    }
  }
}

provider "pfsense" {
  url             = "https://192.168.1.1"
  username        = "admin"
  password        = var.pfsense_password
  tls_skip_verify = true
}

# Create a firewall IP alias
resource "pfsense_firewall_ip_alias" "trusted_hosts" {
  name        = "trusted_hosts"
  description = "Trusted management hosts"
  entries = [
    { ip = "10.0.0.10", description = "Admin workstation" },
    { ip = "10.0.0.11", description = "Monitoring server" },
  ]
  apply = true
}

# Create a firewall rule
resource "pfsense_firewall_rule" "allow_trusted" {
  interface   = "lan"
  type        = "pass"
  ip_protocol = "inet"
  protocol    = "tcp"
  source      = "trusted_hosts"
  destination = "any"
  description = "Allow trusted hosts outbound"
  apply       = true
}
```

## Provider Configuration

| Attribute | Type | Default | Env Var | Description |
|---|---|---|---|---|
| `url` | `string` | `https://192.168.1.1` | `TF_PFSENSE_URL` | pfSense web UI URL |
| `username` | `string` | `admin` | `TF_PFSENSE_USERNAME` | Login username |
| `password` | `string` | -- | `TF_PFSENSE_PASSWORD` | Login password (sensitive) |
| `tls_skip_verify` | `bool` | `false` | -- | Skip TLS certificate verification |
| `max_attempts` | `int` | `3` | -- | Max retry attempts for transient errors |
| `concurrent_writes` | `bool` | `false` | -- | Allow concurrent config writes (see [Concurrency](#concurrency)) |

### Environment Variables

All connection settings can be provided via environment variables, which is useful for CI/CD:

```hcl
# Credentials from TF_PFSENSE_URL, TF_PFSENSE_USERNAME, TF_PFSENSE_PASSWORD
provider "pfsense" {}
```

## Resources

### Firewall

| Resource | Description |
|---|---|
| `pfsense_firewall_rule` | Firewall rules (pass/block/reject) per interface |
| `pfsense_firewall_ip_alias` | IP-based firewall aliases |
| `pfsense_firewall_port_alias` | Port-based firewall aliases |
| `pfsense_firewall_url_alias` | URL table aliases (external blocklists, IP feeds) |

### Networking

| Resource | Description |
|---|---|
| `pfsense_gateway` | Network gateways |
| `pfsense_gateway_group` | Gateway groups for failover/load balancing |
| `pfsense_route` | Static routes |
| `pfsense_interface` | Network interface configuration |
| `pfsense_interface_group` | Interface groups |
| `pfsense_vlan` | 802.1Q VLANs |

### DNS Resolver (Unbound)

| Resource | Description |
|---|---|
| `pfsense_dnsresolver_general` | General DNS resolver settings (singleton) |
| `pfsense_dnsresolver_advanced` | Advanced DNS resolver settings (singleton) |
| `pfsense_dnsresolver_host_override` | DNS host overrides (A/AAAA records) |
| `pfsense_dnsresolver_domain_override` | DNS domain overrides (forwarding zones) |
| `pfsense_dnsresolver_config_file` | Custom Unbound configuration files |

### DHCP

| Resource | Description |
|---|---|
| `pfsense_dhcpv4_server` | DHCPv4 server per interface (singleton per interface) |
| `pfsense_dhcpv4_static_mapping` | DHCP static mappings (MAC-to-IP reservations) |

### VPN

| Resource | Description |
|---|---|
| `pfsense_ipsec_phase1` | IPsec Phase 1 (IKE) tunnels |
| `pfsense_ipsec_phase2` | IPsec Phase 2 (SA) tunnels |

### System

| Resource | Description |
|---|---|
| `pfsense_system_general` | Hostname, domain, DNS servers, timezone (singleton) |
| `pfsense_system_advanced_admin` | Web UI, SSH, console settings (singleton) |
| `pfsense_system_advanced_firewall` | Firewall optimization, state table (singleton) |
| `pfsense_system_advanced_misc` | Power savings, crypto, thermal (singleton) |
| `pfsense_system_advanced_networking` | Hardware offloading, ARP, IPv6 (singleton) |
| `pfsense_system_advanced_notifications` | SMTP, Growl, Telegram notifications (singleton) |
| `pfsense_system_tunable` | System tunables (sysctl) |
| `pfsense_cron_job` | Cron job management |

### Users & Groups

| Resource | Description |
|---|---|
| `pfsense_user` | Local user accounts |
| `pfsense_group` | Local user groups |
| `pfsense_auth_server` | LDAP/RADIUS authentication servers |

### Advanced

| Resource | Description |
|---|---|
| `pfsense_execute_php_command` | Execute arbitrary PHP on the pfSense box (escape hatch) |

## Data Sources

Every resource has a corresponding data source for reading existing configuration. Collection resources also have plural data sources that return all items:

```hcl
# Read a single gateway
data "pfsense_gateway" "wan" {
  name = "WAN_DHCP"
}

# List all gateways
data "pfsense_gateways" "all" {}

# Read system version
data "pfsense_system_version" "this" {}
```

**52 data sources** are available. See the [provider documentation](https://registry.terraform.io/providers/emanuelbesliu/pfsense/latest/docs) for the complete list.

## Key Concepts

### Singleton vs Collection Resources

**Singleton resources** represent configuration that always exists on pfSense (e.g., system general settings, DNS resolver config). Creating the resource adopts the existing config into Terraform state. Deleting it resets the settings to defaults -- it does not remove the configuration section from pfSense.

**Collection resources** represent items that can be individually created and destroyed (e.g., firewall rules, host overrides, static routes).

### Apply Attribute

Most resources support an `apply` attribute (defaults to `true`) that triggers pfSense to reload the relevant subsystem after changes. For example, firewall rules reload the filter, DNS changes restart Unbound, and DHCP changes restart the DHCP daemon. Set `apply = false` if you are making many changes and want to batch the reload.

### Import

All resources support `terraform import`. Import IDs vary by resource type:

| Resource | Import ID Format | Example |
|---|---|---|
| Gateway | `name` | `WAN_DHCP` |
| Static route | `network/gateway` | `10.0.0.0/24/WAN_GW` |
| DHCP static mapping | `interface,mac_address` | `lan,aa:bb:cc:dd:ee:ff` |
| DNS host override | `host,domain` | `myhost,example.com` |
| DNS domain override | `domain` | `corp.example.com` |
| IPsec Phase 1 | `ikeid` | `1` |
| IPsec Phase 2 | `uniqid` | `67e3a8b04e127` |
| Cron job | `command` | `/root/renew-cert.sh` |
| Firewall rule | `tracker` | `1234567890` |
| Singleton resources | any string (e.g., `pfsense`) | `pfsense` |

### Concurrency

pfSense stores all configuration in a single XML file (`/cf/conf/config.xml`). Concurrent writes can cause one change to overwrite another. By default, the provider serializes all write operations (`concurrent_writes = false`). This is slower but safe.

If you enable `concurrent_writes = true`, the provider uses per-subsystem locks but allows writes to different subsystems in parallel. This is faster but carries a small risk of config conflicts if pfSense itself or another session writes at the same time.

### Execute PHP Command

The `pfsense_execute_php_command` resource and data source provide an escape hatch for operations not yet covered by dedicated resources. The PHP code is executed on pfSense via `diag_command.php` and must output exactly one JSON value:

```hcl
# Read-only: query current config
data "pfsense_execute_php_command" "wan_ip" {
  command = "print(json_encode(get_interface_ip('wan')));"
}

# Write: modify config with optional cleanup
resource "pfsense_execute_php_command" "custom_setting" {
  command         = "config_set_path('custom/key', 'value'); write_config('Terraform'); print(json_encode(true));"
  destroy_command = "config_del_path('custom/key'); write_config('Terraform'); print(json_encode(true));"
}
```

## Building from Source

Requirements: [Go](https://golang.org/doc/install) 1.24+, [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.7+

```shell
git clone https://github.com/emanuelbesliu/terraform-provider-pfsense.git
cd terraform-provider-pfsense
go install .
```

To use a local build, add a [dev override](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) to your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "emanuelbesliu/pfsense" = "/path/to/go/bin"
  }
  direct {}
}
```

## Contributing

```shell
# Build
go install .

# Lint
go vet ./...

# Generate documentation
make docs

# Run acceptance tests (requires a pfSense instance)
make test/acc
```

## License

See [LICENSE](LICENSE) for details.
