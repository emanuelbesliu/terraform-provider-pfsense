package pfsense

import (
	"encoding/json"
	"testing"
)

func TestFlexBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{"bool true", `true`, true},
		{"bool false", `false`, false},
		{"string yes", `"yes"`, true},
		{"string empty", `""`, false},
		{"string zero", `"0"`, false},
		{"string one", `"1"`, true},
		{"json null", `null`, false},
		{"string no", `"no"`, false},
		{"string disabled", `"disabled"`, false},
		{"string enabled", `"enabled"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b flexBool
			if err := json.Unmarshal([]byte(tt.data), &b); err != nil {
				t.Fatalf("Unmarshal(%s) unexpected error: %v", tt.data, err)
			}

			if bool(b) != tt.want {
				t.Errorf("flexBool(%s) = %v, want %v", tt.data, bool(b), tt.want)
			}
		})
	}
}

func TestFlexString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string
	}{
		{"string", `"1194"`, "1194"},
		{"number", `1194`, "1194"},
		{"empty string", `""`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s flexString
			if err := json.Unmarshal([]byte(tt.data), &s); err != nil {
				t.Fatalf("Unmarshal(%s) unexpected error: %v", tt.data, err)
			}

			if string(s) != tt.want {
				t.Errorf("flexString(%s) = %q, want %q", tt.data, string(s), tt.want)
			}
		})
	}
}

func TestSplitCommaList(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"whitespace", "   ", nil},
		{"single", "AES-256-GCM", []string{"AES-256-GCM"}},
		{"multiple", "AES-256-GCM,AES-128-GCM, CHACHA20-POLY1305", []string{"AES-256-GCM", "AES-128-GCM", "CHACHA20-POLY1305"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCommaList(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("splitCommaList(%q) = %v, want %v", tt.in, got, tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitCommaList(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// openVPNServerSampleJSON mirrors the shape pfSense returns from
// config_get_path('openvpn/openvpn-server'): presence-based booleans stored as
// PHP true, numeric-ish fields as strings, and comma-joined list fields.
const openVPNServerSampleJSON = `[
  {
    "vpnid": "1",
    "mode": "server_tls_user",
    "authmode": "Local Database",
    "dev_mode": "tun",
    "protocol": "UDP4",
    "interface": "wan",
    "local_port": "1194",
    "description": "Remote access VPN",
    "caref": "abc123caref",
    "certref": "abc123certref",
    "dh_length": "2048",
    "data_ciphers": "AES-256-GCM,AES-128-GCM,CHACHA20-POLY1305",
    "data_ciphers_fallback": "AES-256-CBC",
    "digest": "SHA256",
    "tunnel_network": "10.8.0.0/24",
    "topology": "subnet",
    "dns_server_enable": true,
    "dns_server1": "10.0.0.1",
    "gwredir": true,
    "client2client": true,
    "max_clients": "50",
    "username_as_common_name": "enabled",
    "verbosity_level": "3"
  },
  {
    "vpnid": "2",
    "mode": "p2p_shared_key",
    "dev_mode": "tun",
    "protocol": "UDP4",
    "interface": "wan",
    "local_port": "1195",
    "description": "Site-to-site tunnel",
    "shared_key": "-----BEGIN OpenVPN Static key V1-----",
    "tunnel_network": "10.9.0.0/30",
    "remote_network": "192.168.20.0/24",
    "data_ciphers_fallback": "AES-256-CBC",
    "digest": "SHA256"
  }
]`

func parseOpenVPNServersFromJSON(t *testing.T, data string) OpenVPNServers {
	t.Helper()

	var resp []openVPNServerResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatalf("failed to unmarshal sample: %v", err)
	}

	servers := make(OpenVPNServers, 0, len(resp))
	for index, r := range resp {
		servers = append(servers, parseOpenVPNServerResponse(r, index))
	}

	return servers
}

func TestParseOpenVPNServerResponse(t *testing.T) {
	servers := parseOpenVPNServersFromJSON(t, openVPNServerSampleJSON)

	if len(servers) != 2 {
		t.Fatalf("parsed %d servers, want 2", len(servers))
	}

	s1 := servers[0]
	if s1.VPNID != "1" {
		t.Errorf("VPNID = %q, want 1", s1.VPNID)
	}
	if s1.controlID != 0 {
		t.Errorf("controlID = %d, want 0", s1.controlID)
	}
	if s1.Mode != "server_tls_user" {
		t.Errorf("Mode = %q, want server_tls_user", s1.Mode)
	}
	if len(s1.AuthMode) != 1 || s1.AuthMode[0] != "Local Database" {
		t.Errorf("AuthMode = %v, want [Local Database]", s1.AuthMode)
	}
	wantCiphers := []string{"AES-256-GCM", "AES-128-GCM", "CHACHA20-POLY1305"}
	if len(s1.DataCiphers) != len(wantCiphers) {
		t.Fatalf("DataCiphers = %v, want %v", s1.DataCiphers, wantCiphers)
	}
	for i := range wantCiphers {
		if s1.DataCiphers[i] != wantCiphers[i] {
			t.Errorf("DataCiphers[%d] = %q, want %q", i, s1.DataCiphers[i], wantCiphers[i])
		}
	}
	if !s1.DNSServerEnable {
		t.Error("DNSServerEnable = false, want true (bool true in JSON)")
	}
	if !s1.GWRedir {
		t.Error("GWRedir = false, want true")
	}
	if !s1.Client2Client {
		t.Error("Client2Client = false, want true")
	}
	if s1.LocalPort != "1194" {
		t.Errorf("LocalPort = %q, want 1194", s1.LocalPort)
	}
	if s1.UsernameAsCommonName != "enabled" {
		t.Errorf("UsernameAsCommonName = %q, want enabled", s1.UsernameAsCommonName)
	}
	// Absent boolean keys must default to false.
	if s1.Disable {
		t.Error("Disable = true, want false (key absent)")
	}
	if s1.PassTOS {
		t.Error("PassTOS = true, want false (key absent)")
	}

	s2 := servers[1]
	if s2.controlID != 1 {
		t.Errorf("controlID = %d, want 1", s2.controlID)
	}
	if s2.Mode != "p2p_shared_key" {
		t.Errorf("Mode = %q, want p2p_shared_key", s2.Mode)
	}
	if s2.RemoteNetwork != "192.168.20.0/24" {
		t.Errorf("RemoteNetwork = %q, want 192.168.20.0/24", s2.RemoteNetwork)
	}
	if len(s2.AuthMode) != 0 {
		t.Errorf("AuthMode = %v, want empty", s2.AuthMode)
	}
}

func TestOpenVPNServers_Lookup(t *testing.T) {
	servers := parseOpenVPNServersFromJSON(t, openVPNServerSampleJSON)

	server, err := servers.GetByVPNID("2")
	if err != nil {
		t.Fatalf("GetByVPNID(2) error: %v", err)
	}
	if server.Mode != "p2p_shared_key" {
		t.Errorf("GetByVPNID(2).Mode = %q, want p2p_shared_key", server.Mode)
	}

	controlID, err := servers.GetControlIDByVPNID("2")
	if err != nil {
		t.Fatalf("GetControlIDByVPNID(2) error: %v", err)
	}
	if *controlID != 1 {
		t.Errorf("GetControlIDByVPNID(2) = %d, want 1", *controlID)
	}

	if _, err := servers.GetByVPNID("999"); err == nil {
		t.Error("GetByVPNID(999) expected error, got nil")
	}
}

func TestOpenVPNServerFormValues(t *testing.T) {
	server := OpenVPNServer{
		Mode:                 "server_tls_user",
		DevMode:              "tun",
		Protocol:             "UDP4",
		Interface:            "wan",
		IPAddr:               "203.0.113.5",
		LocalPort:            "1194",
		Description:          "Remote access VPN",
		Disable:              true,
		AuthMode:             []string{"Local Database", "MyLDAP"},
		TLS:                  "-----BEGIN OpenVPN Static key V1-----",
		TLSType:              "auth",
		DataCiphers:          []string{"AES-256-GCM", "AES-128-GCM"},
		Client2Client:        true,
		UsernameAsCommonName: "enabled",
	}

	values := openVPNServerFormValues(server)

	if values.Get("save") != "Save" {
		t.Errorf("save = %q, want Save", values.Get("save"))
	}
	// Interface with virtual IP suffix.
	if got := values.Get("interface"); got != "wan|203.0.113.5" {
		t.Errorf("interface = %q, want wan|203.0.113.5", got)
	}
	// Presence-based booleans emit "yes".
	if values.Get("disable") != "yes" {
		t.Errorf("disable = %q, want yes", values.Get("disable"))
	}
	if values.Get("client2client") != "yes" {
		t.Errorf("client2client = %q, want yes", values.Get("client2client"))
	}
	// TLS presence enables tlsauth_enable.
	if values.Get("tlsauth_enable") != "yes" {
		t.Errorf("tlsauth_enable = %q, want yes", values.Get("tlsauth_enable"))
	}
	// username_as_common_name == "enabled" -> checkbox "yes".
	if values.Get("username_as_common_name") != "yes" {
		t.Errorf("username_as_common_name = %q, want yes", values.Get("username_as_common_name"))
	}
	// Array fields.
	authModes := values["authmode[]"]
	if len(authModes) != 2 || authModes[0] != "Local Database" || authModes[1] != "MyLDAP" {
		t.Errorf("authmode[] = %v, want [Local Database MyLDAP]", authModes)
	}
	ciphers := values["data_ciphers[]"]
	if len(ciphers) != 2 || ciphers[0] != "AES-256-GCM" || ciphers[1] != "AES-128-GCM" {
		t.Errorf("data_ciphers[] = %v, want [AES-256-GCM AES-128-GCM]", ciphers)
	}
	// Absent boolean must not be present.
	if _, ok := values["pass_tos"]; ok {
		t.Error("pass_tos should not be set when false")
	}
}

func TestOpenVPNServerFormValues_UsernameDisabled(t *testing.T) {
	server := OpenVPNServer{
		Mode:                 "server_tls",
		DevMode:              "tun",
		Protocol:             "UDP4",
		Interface:            "lan",
		UsernameAsCommonName: "disabled",
	}

	values := openVPNServerFormValues(server)

	if _, ok := values["username_as_common_name"]; ok {
		t.Error("username_as_common_name should not be set when disabled")
	}
	// No virtual IP suffix when IPAddr empty.
	if got := values.Get("interface"); got != "lan" {
		t.Errorf("interface = %q, want lan", got)
	}
}
