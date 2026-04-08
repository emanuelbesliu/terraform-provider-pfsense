package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type CertificateAuthorityModel struct {
	RefID         types.String `tfsdk:"refid"`
	Descr         types.String `tfsdk:"descr"`
	Certificate   types.String `tfsdk:"certificate"`
	PrivateKey    types.String `tfsdk:"private_key"`
	Trust         types.Bool   `tfsdk:"trust"`
	RandomSerial  types.Bool   `tfsdk:"random_serial"`
	NextSerial    types.Int64  `tfsdk:"next_serial"`
	Subject       types.String `tfsdk:"subject"`
	Issuer        types.String `tfsdk:"issuer"`
	Serial        types.String `tfsdk:"serial"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
	IsSelfSigned  types.Bool   `tfsdk:"is_self_signed"`
	ValidFrom     types.String `tfsdk:"valid_from"`
	ValidTo       types.String `tfsdk:"valid_to"`
	InUse         types.Bool   `tfsdk:"in_use"`
}

func (CertificateAuthorityModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"refid": {
			Description: "Unique reference ID assigned by pfSense.",
		},
		"descr": {
			Description: "Descriptive name for the certificate authority.",
		},
		"certificate": {
			Description: "PEM-encoded CA certificate.",
		},
		"private_key": {
			Description: "PEM-encoded private key for the CA (optional, only needed for signing).",
		},
		"trust": {
			Description: "Whether the CA is added to the operating system trust store.",
		},
		"random_serial": {
			Description: "Whether to use random serial numbers for issued certificates.",
		},
		"next_serial": {
			Description: "Next serial number to use when issuing certificates (sequential mode).",
		},
		"subject": {
			Description: "Certificate subject (distinguished name).",
		},
		"issuer": {
			Description: "Certificate issuer (distinguished name).",
		},
		"serial": {
			Description: "Certificate serial number.",
		},
		"has_private_key": {
			Description: "Whether the CA has a private key stored.",
		},
		"is_self_signed": {
			Description: "Whether the CA certificate is self-signed.",
		},
		"valid_from": {
			Description: "Certificate validity start date.",
		},
		"valid_to": {
			Description: "Certificate validity end date.",
		},
		"in_use": {
			Description: "Whether the CA is referenced by certificates, CRLs, or other configuration.",
		},
	}
}

func (CertificateAuthorityModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"refid":           types.StringType,
		"descr":           types.StringType,
		"certificate":     types.StringType,
		"private_key":     types.StringType,
		"trust":           types.BoolType,
		"random_serial":   types.BoolType,
		"next_serial":     types.Int64Type,
		"subject":         types.StringType,
		"issuer":          types.StringType,
		"serial":          types.StringType,
		"has_private_key": types.BoolType,
		"is_self_signed":  types.BoolType,
		"valid_from":      types.StringType,
		"valid_to":        types.StringType,
		"in_use":          types.BoolType,
	}
}

func (m *CertificateAuthorityModel) Set(_ context.Context, ca pfsense.CertificateAuthority) diag.Diagnostics {
	var diags diag.Diagnostics

	m.RefID = types.StringValue(ca.RefID)
	m.Descr = types.StringValue(ca.Descr)
	m.Certificate = types.StringValue(ca.Certificate)

	if ca.PrivateKey != "" {
		m.PrivateKey = types.StringValue(ca.PrivateKey)
	} else {
		m.PrivateKey = types.StringNull()
	}

	m.Trust = types.BoolValue(ca.Trust)
	m.RandomSerial = types.BoolValue(ca.RandomSerial)
	m.NextSerial = types.Int64Value(int64(ca.NextSerial))
	m.Subject = types.StringValue(ca.Subject)
	m.Issuer = types.StringValue(ca.Issuer)
	m.Serial = types.StringValue(ca.Serial)
	m.HasPrivateKey = types.BoolValue(ca.HasPrivateKey)
	m.IsSelfSigned = types.BoolValue(ca.IsSelfSigned)

	if ca.ValidFrom != "" {
		m.ValidFrom = types.StringValue(ca.ValidFrom)
	} else {
		m.ValidFrom = types.StringNull()
	}

	if ca.ValidTo != "" {
		m.ValidTo = types.StringValue(ca.ValidTo)
	} else {
		m.ValidTo = types.StringNull()
	}

	m.InUse = types.BoolValue(ca.InUse)

	return diags
}

func (m CertificateAuthorityModel) Value(_ context.Context, ca *pfsense.CertificateAuthority) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("descr"),
		"CA description cannot be parsed",
		ca.SetDescr(m.Descr.ValueString()),
	)

	if !m.Certificate.IsNull() && !m.Certificate.IsUnknown() {
		addPathError(
			&diags,
			path.Root("certificate"),
			"CA certificate cannot be parsed",
			ca.SetCertificate(m.Certificate.ValueString()),
		)
	}

	if !m.PrivateKey.IsNull() {
		addPathError(
			&diags,
			path.Root("private_key"),
			"CA private key cannot be parsed",
			ca.SetPrivateKey(m.PrivateKey.ValueString()),
		)
	}

	if !m.Trust.IsNull() && !m.Trust.IsUnknown() {
		addPathError(
			&diags,
			path.Root("trust"),
			"CA trust cannot be parsed",
			ca.SetTrust(m.Trust.ValueBool()),
		)
	}

	if !m.NextSerial.IsNull() && !m.NextSerial.IsUnknown() {
		addPathError(
			&diags,
			path.Root("next_serial"),
			"CA next serial cannot be parsed",
			ca.SetNextSerial(int(m.NextSerial.ValueInt64())),
		)
	}

	if !m.RandomSerial.IsNull() && !m.RandomSerial.IsUnknown() {
		ca.RandomSerial = m.RandomSerial.ValueBool()
	}

	return diags
}
