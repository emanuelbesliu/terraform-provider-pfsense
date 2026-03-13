# Manage a RADIUS authentication server
resource "pfsense_auth_server" "radius_example" {
  name             = "Corporate RADIUS"
  type             = "radius"
  host             = "203.0.113.10"
  radius_protocol  = "MSCHAPv2"
  radius_secret    = "shared-secret-here"
  radius_auth_port = "1812"
  radius_acct_port = "1813"
  radius_timeout   = "5"
  radius_srvcs     = "both"
}

# Manage an LDAP authentication server (Active Directory)
resource "pfsense_auth_server" "ldap_example" {
  name             = "Corporate AD"
  type             = "ldap"
  host             = "198.51.100.10"
  ldap_port        = "389"
  ldap_urltype     = "Standard TCP"
  ldap_protver     = "3"
  ldap_scope       = "subtree"
  ldap_basedn      = "DC=example,DC=com"
  ldap_authcn      = "CN=Users;OU=Staff"
  ldap_binddn      = "CN=svc-pfsense,OU=ServiceAccounts,DC=example,DC=com"
  ldap_bindpw      = "bind-password-here"
  ldap_timeout     = "25"
  ldap_attr_user   = "samAccountName"
  ldap_attr_group  = "cn"
  ldap_attr_member = "memberOf"
}
