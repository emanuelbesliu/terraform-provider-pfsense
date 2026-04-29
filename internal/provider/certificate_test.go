package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCertificateResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig("tf-acc-test-cert", "server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_certificate.test", "descr", "tf-acc-test-cert"),
					resource.TestCheckResourceAttr("pfsense_certificate.test", "type", "server"),
					resource.TestCheckResourceAttrSet("pfsense_certificate.test", "refid"),
					resource.TestCheckResourceAttrSet("pfsense_certificate.test", "subject"),
					resource.TestCheckResourceAttrSet("pfsense_certificate.test", "issuer"),
					resource.TestCheckResourceAttrSet("pfsense_certificate.test", "serial"),
					resource.TestCheckResourceAttr("pfsense_certificate.test", "has_private_key", "true"),
					resource.TestCheckResourceAttr("pfsense_certificate.test", "is_self_signed", "true"),
				),
			},
			{
				ResourceName:            "pfsense_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "private_key"},
			},
			{
				Config: testAccCertificateConfig("tf-acc-test-cert-updated", "server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_certificate.test", "descr", "tf-acc-test-cert-updated"),
				),
			},
		},
	})
}

func TestAccCertificateResource_user(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig("tf-acc-test-user-cert", "user"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_certificate.test", "descr", "tf-acc-test-user-cert"),
					resource.TestCheckResourceAttr("pfsense_certificate.test", "type", "user"),
				),
			},
		},
	})
}

func TestAccCertificatesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig("tf-acc-test-ds-cert", "server") + `
data "pfsense_certificates" "all" {
  depends_on = [pfsense_certificate.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_certificates.all", "certificates.#"),
				),
			},
		},
	})
}

func TestAccCertificateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig("tf-acc-test-single-ds", "server") + `
data "pfsense_certificate" "test" {
  descr      = "tf-acc-test-single-ds"
  depends_on = [pfsense_certificate.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_certificate.test", "descr", "tf-acc-test-single-ds"),
					resource.TestCheckResourceAttrSet("data.pfsense_certificate.test", "refid"),
					resource.TestCheckResourceAttrSet("data.pfsense_certificate.test", "subject"),
				),
			},
		},
	})
}

// testAccCertificateConfig generates a self-signed certificate inline for testing.
// Uses a pre-generated self-signed cert/key pair to avoid external dependencies.
func testAccCertificateConfig(descr string, certType string) string {
	return fmt.Sprintf(`
resource "pfsense_certificate" "test" {
  descr       = %q
  type        = %q
  certificate = <<-EOT
-----BEGIN CERTIFICATE-----
MIICpDCCAYwCCQDU+pQ4pHgSpDANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAls
b2NhbGhvc3QwHhcNMjQwMTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjAUMRIwEAYD
VQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7
o4qne60TB3pOYaBy/YPaS3JTSqMkfr5fPb3S+Sw1x6MrFCKJfBwi2DmgDqIhQmg
lI47JMBmSE2DNOEM4wM0KKMRLaxnPH0MYBwVRMkdCCnnuBCFVFxqTYNFP2Jl8yLj
bGIQB0BJMWkrnPFHMtG7cRNPRrFIwLBWF2ECwADhw/qGGCcHBMZJJnAOfRwSNxsW
1PeDLPV0VDBzxQBQjzBEOHPaVKnov9bfUbFiDhWHkp9bsFMEnCLnPTMJKKPB5mRg
O+C+sSBFGPJAM6mVfYOl6R6mTyhFnM9JzSeBVbFmXVmJIEKhiHRKYmNpjUqcMY7y
S1+MF9GxCXJETreyLuupAgMBAAEwDQYJKoZIhvcNAQELBQADggEBABkADzxUKy0a
mVnMmC2oAKPJbATdOfdMdCjMpGMAAB4SiO4AtRqwVyOvj+kJFs/PjlUX2hLkfMSi
qVTqyDwDbB4MKaR5CVDvZRfdVDiUyWbHmFiCG0nlwFnGjYUuaXF0YV5TQ+8txHLn
v1DkGOJGjVR/08UVlGCxq3JCRZ/YcdNpgUn3LT26ickBhWRfFwGPLSPFaPWNO2sS
QeYaccaRLM3o6LMiFJGOqJz8gBASnoDJqoO8RF/vBiFE3VbMzSSMXQa7C1E8r1SS
FKRTJvSj9n5tP9UDr9k8SLoMaB5b0cOf0lqB7TpD6nFq/PZCFpjJ7HQfNhH8pYsN
dCmU7CKUJXE=
-----END CERTIFICATE-----
  EOT
  private_key = <<-EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7o4qne60TB3pO
YaBy/YPaS3JTSqMkfr5fPb3S+Sw1x6MrFCKJfBwi2DmgDqIhQmglI47JMBmSE2D
NOEM4wM0KKMRLaxnPH0MYBwVRMkdCCnnuBCFVFxqTYNFP2Jl8yLjbGIQB0BJMWkr
nPFHMtG7cRNPRrFIwLBWF2ECwADhw/qGGCcHBMZJJnAOfRwSNxsW1PeDLPV0VDBz
xQBQjzBEOHPaVKnov9bfUbFiDhWHkp9bsFMEnCLnPTMJKKPB5mRgO+C+sSBFGPJA
M6mVfYOl6R6mTyhFnM9JzSeBVbFmXVmJIEKhiHRKYmNpjUqcMY7yS1+MF9GxCXJE
TreyLuupAgMBAAECggEAQj2mlBxR8ckDOBBFJXbDMwaHGsTKKM6QPMMI5oYVU04e0
m7GFwJb6mUsNBqGIZUaNx0/TKPjqoEi7Bd7EuUYDPzCF9mRk7fIdEKAsjBuWbG/F
FTJECHkqFCgVmQzMlr/D/D+29JGt7DamkYQKaUV1CjHX7qLVCDlJu+MxV8k/N7Qh
vRKhWNMk7fJD/lLPWjBaJlMPkMYs/F/bYwUKJ6ksWaP0jlXS3wVDlBQ6/HQXEAaO
b4NhD/cmgmyjqUDKbJHqxFzOkh6P3W4mV1lwNXJMy5F3VGbI3qViLqkiAQwqBniH
1uMdXhpcROzrQ+OWX8EhkJqH/CTJKHnlsJDKZQKBgQDsJNnBOqiglhah/lrzkQN7
T+G2n1jYQVdFikUsiUCFLNypT+1eTjBz0tr7ydFvExGIljJtLzsPZfIAMg7xSD7c
F3M8HKOGMnW0p5fDaGEMnUBN+DP/UuSbhJzfoIvfjPp0PsOuhdC2PjEr5gDPVAV+
cWjG7F94L2L/FbHl0RaOdwKBgQDMHN1FLR3jHdjkCiES4T+3mVAjhBaI2VmxcQ5m
c1K0pcudBGU7OThpNzPbKzQJFOHDlzTp/qi9PYXdZp9bJfBW0snM2DvK0NUdSbs3
w3fNf0FBgXZNOb+9E0qf+g58mUBe3/rMCkx2i1A0pB4w2w5ynRm7i5I0myREXp0R
PwjHnwKBgBg7RFw8mLAEYVIc9kV2sCLqxB3JrmGUhS1Bej51wBJdlPf0o0BLBGrn
VVZahePjSJHZdOa1r3II8rC3uDdRNrwp+mYN32qJymMHnVEPMtymJbJOeSPFCDsJ
UePHpAdN6MR7DnXdNaMk4dJkpPrKaLFRoWM9d7ILjV6IDirDYZ0RAoGBALRo7B7v
V+ShqFaIekuMQw7m0YPFaESthobxHrjCbIuFPVgBDy4xdhDVGpTVMtFUfGLaMQFH
K+rcKz2BFJjZ2L9m0s1Wai+kJqsl7GVPKNoU+Gp+uAH+7R2EgRiMnJ/9YibW5yJ
o1Kx+H9aMpzIR3A1RCkVNN1M8dNN8BZrOiGNAoGBAMV3gHVE1XiR0OeSPxM4BkmY
IYF3lQdzP4Hg2K8yXgFzCjE5cJQH1yMM3sNb/OOgxSr9b4yl+QEnIrmJFl9oGVhv
dGy6JcJkJYH2z/n1FY0h7N9ADriSfUPFBwrJ7vqBa5qhTO9ixBrnC8zM9bfz0Lmt
GHVNAZklnPDrZzfq72gb
-----END PRIVATE KEY-----
  EOT
}
`, descr, certType)
}
