package pfsense

import (
	"testing"
)

func TestCertificate_SetDescr(t *testing.T) {
	tests := []struct {
		name    string
		descr   string
		wantErr bool
	}{
		{"valid", "My Certificate", false},
		{"short", "a", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Certificate

			err := c.SetDescr(tt.descr)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDescr(%q) error = %v, wantErr %v", tt.descr, err, tt.wantErr)
			}

			if err == nil && c.Descr != tt.descr {
				t.Errorf("SetDescr(%q) got %q", tt.descr, c.Descr)
			}
		})
	}
}

func TestCertificate_SetCertType(t *testing.T) {
	tests := []struct {
		name    string
		ct      string
		wantErr bool
	}{
		{"server", "server", false},
		{"user", "user", false},
		{"invalid", "ca", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Certificate

			err := c.SetCertType(tt.ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCertType(%q) error = %v, wantErr %v", tt.ct, err, tt.wantErr)
			}

			if err == nil && c.CertType != tt.ct {
				t.Errorf("SetCertType(%q) got %q", tt.ct, c.CertType)
			}
		})
	}
}

func TestCertificate_SetCertificate(t *testing.T) {
	validPEM := "-----BEGIN CERTIFICATE-----\nMIIBxTCCAWugAwIBAgI...\n-----END CERTIFICATE-----"

	tests := []struct {
		name    string
		cert    string
		wantErr bool
	}{
		{"valid_pem", validPEM, false},
		{"not_pem", "not-a-certificate", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Certificate

			err := c.SetCertificate(tt.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCertificate(%q) error = %v, wantErr %v", tt.cert, err, tt.wantErr)
			}

			if err == nil && c.Certificate != tt.cert {
				t.Errorf("SetCertificate() got %q", c.Certificate)
			}
		})
	}
}

func TestCertificate_SetPrivateKey(t *testing.T) {
	var c Certificate

	err := c.SetPrivateKey("-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----")
	if err != nil {
		t.Errorf("SetPrivateKey() unexpected error = %v", err)
	}

	err = c.SetPrivateKey("")
	if err != nil {
		t.Errorf("SetPrivateKey(\"\") unexpected error = %v", err)
	}
}

func TestCertificate_SetCARef(t *testing.T) {
	var c Certificate

	err := c.SetCARef("abc123")
	if err != nil {
		t.Errorf("SetCARef() unexpected error = %v", err)
	}

	if c.CARef != "abc123" {
		t.Errorf("SetCARef() got %q, want %q", c.CARef, "abc123")
	}

	err = c.SetCARef("")
	if err != nil {
		t.Errorf("SetCARef(\"\") unexpected error = %v", err)
	}
}

func TestCertificate_CertTypes(t *testing.T) {
	c := Certificate{}
	types := c.CertTypes()

	if len(types) != 2 {
		t.Fatalf("CertTypes() len = %d, want 2", len(types))
	}

	if types[0] != "server" {
		t.Errorf("CertTypes()[0] = %q, want %q", types[0], "server")
	}

	if types[1] != "user" {
		t.Errorf("CertTypes()[1] = %q, want %q", types[1], "user")
	}
}

func TestParseCertificateResponse(t *testing.T) {
	resp := certificateResponse{
		RefID:         "abc123",
		Descr:         "Test Cert",
		CertType:      "server",
		CARef:         "ca456",
		Subject:       "CN=test.example.com",
		Issuer:        "CN=Example CA",
		Serial:        "1234",
		HasPrivateKey: true,
		IsSelfSigned:  false,
		ValidFrom:     "Jan  1 00:00:00 2024 GMT",
		ValidTo:       "Dec 31 23:59:59 2025 GMT",
		InUse:         true,
		Certificate:   "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
	}

	cert := parseCertificateResponse(resp)

	if cert.RefID != "abc123" {
		t.Errorf("RefID = %q, want %q", cert.RefID, "abc123")
	}

	if cert.Descr != "Test Cert" {
		t.Errorf("Descr = %q, want %q", cert.Descr, "Test Cert")
	}

	if cert.CertType != "server" {
		t.Errorf("CertType = %q, want %q", cert.CertType, "server")
	}

	if cert.CARef != "ca456" {
		t.Errorf("CARef = %q, want %q", cert.CARef, "ca456")
	}

	if cert.Subject != "CN=test.example.com" {
		t.Errorf("Subject = %q, want %q", cert.Subject, "CN=test.example.com")
	}

	if cert.Issuer != "CN=Example CA" {
		t.Errorf("Issuer = %q, want %q", cert.Issuer, "CN=Example CA")
	}

	if cert.Serial != "1234" {
		t.Errorf("Serial = %q, want %q", cert.Serial, "1234")
	}

	if !cert.HasPrivateKey {
		t.Error("HasPrivateKey = false, want true")
	}

	if cert.IsSelfSigned {
		t.Error("IsSelfSigned = true, want false")
	}

	if !cert.InUse {
		t.Error("InUse = false, want true")
	}

	if cert.Certificate != "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----" {
		t.Errorf("Certificate = %q", cert.Certificate)
	}
}

func TestCertificates_GetByRefID(t *testing.T) {
	certs := Certificates{
		{RefID: "aaa", Descr: "Cert A"},
		{RefID: "bbb", Descr: "Cert B"},
	}

	cert, err := certs.GetByRefID("bbb")
	if err != nil {
		t.Fatalf("GetByRefID() error = %v", err)
	}

	if cert.Descr != "Cert B" {
		t.Errorf("Descr = %q, want %q", cert.Descr, "Cert B")
	}

	_, err = certs.GetByRefID("nonexistent")
	if err == nil {
		t.Fatal("GetByRefID() expected error for nonexistent ID")
	}
}

func TestCertificates_GetByDescr(t *testing.T) {
	certs := Certificates{
		{RefID: "aaa", Descr: "Cert A"},
		{RefID: "bbb", Descr: "Cert B"},
	}

	cert, err := certs.GetByDescr("Cert A")
	if err != nil {
		t.Fatalf("GetByDescr() error = %v", err)
	}

	if cert.RefID != "aaa" {
		t.Errorf("RefID = %q, want %q", cert.RefID, "aaa")
	}

	_, err = certs.GetByDescr("nonexistent")
	if err == nil {
		t.Fatal("GetByDescr() expected error for nonexistent description")
	}
}
