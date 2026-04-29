package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type CertificateModel struct {
	RefID         types.String `tfsdk:"refid"`
	Descr         types.String `tfsdk:"descr"`
	CertType      types.String `tfsdk:"type"`
	CARef         types.String `tfsdk:"caref"`
	Certificate   types.String `tfsdk:"certificate"`
	PrivateKey    types.String `tfsdk:"private_key"`
	Subject       types.String `tfsdk:"subject"`
	Issuer        types.String `tfsdk:"issuer"`
	Serial        types.String `tfsdk:"serial"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
	IsSelfSigned  types.Bool   `tfsdk:"is_self_signed"`
	ValidFrom     types.String `tfsdk:"valid_from"`
	ValidTo       types.String `tfsdk:"valid_to"`
	InUse         types.Bool   `tfsdk:"in_use"`
}

func (CertificateModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"refid": {
			Description: "Unique reference ID assigned by pfSense.",
		},
		"descr": {
			Description: "Descriptive name for the certificate.",
		},
		"type": {
			Description:         fmt.Sprintf("Certificate type. Options: %s.", wrapElementsJoin(pfsense.Certificate{}.CertTypes(), "'")),
			MarkdownDescription: fmt.Sprintf("Certificate type. Options: %s.", wrapElementsJoin(pfsense.Certificate{}.CertTypes(), "`")),
		},
		"caref": {
			Description: "Reference ID of the signing certificate authority. Empty for self-signed or externally signed certificates.",
		},
		"certificate": {
			Description: "PEM-encoded certificate.",
		},
		"private_key": {
			Description: "PEM-encoded private key for the certificate.",
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
			Description: "Whether the certificate has a private key stored.",
		},
		"is_self_signed": {
			Description: "Whether the certificate is self-signed.",
		},
		"valid_from": {
			Description: "Certificate validity start date.",
		},
		"valid_to": {
			Description: "Certificate validity end date.",
		},
		"in_use": {
			Description: "Whether the certificate is referenced by services or other configuration.",
		},
	}
}

func (CertificateModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"refid":           types.StringType,
		"descr":           types.StringType,
		"type":            types.StringType,
		"caref":           types.StringType,
		"certificate":     types.StringType,
		"private_key":     types.StringType,
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

func (m *CertificateModel) Set(_ context.Context, c pfsense.Certificate) diag.Diagnostics {
	var diags diag.Diagnostics

	m.RefID = types.StringValue(c.RefID)
	m.Descr = types.StringValue(c.Descr)
	m.CertType = types.StringValue(c.CertType)

	if c.CARef != "" {
		m.CARef = types.StringValue(c.CARef)
	} else {
		m.CARef = types.StringNull()
	}

	m.Certificate = types.StringValue(c.Certificate)

	if c.PrivateKey != "" {
		m.PrivateKey = types.StringValue(c.PrivateKey)
	} else {
		m.PrivateKey = types.StringNull()
	}

	m.Subject = types.StringValue(c.Subject)
	m.Issuer = types.StringValue(c.Issuer)
	m.Serial = types.StringValue(c.Serial)
	m.HasPrivateKey = types.BoolValue(c.HasPrivateKey)
	m.IsSelfSigned = types.BoolValue(c.IsSelfSigned)

	if c.ValidFrom != "" {
		m.ValidFrom = types.StringValue(c.ValidFrom)
	} else {
		m.ValidFrom = types.StringNull()
	}

	if c.ValidTo != "" {
		m.ValidTo = types.StringValue(c.ValidTo)
	} else {
		m.ValidTo = types.StringNull()
	}

	m.InUse = types.BoolValue(c.InUse)

	return diags
}

func (m CertificateModel) Value(_ context.Context, c *pfsense.Certificate) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("descr"),
		"Certificate description cannot be parsed",
		c.SetDescr(m.Descr.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("type"),
		"Certificate type cannot be parsed",
		c.SetCertType(m.CertType.ValueString()),
	)

	if !m.Certificate.IsNull() && !m.Certificate.IsUnknown() {
		addPathError(
			&diags,
			path.Root("certificate"),
			"Certificate cannot be parsed",
			c.SetCertificate(m.Certificate.ValueString()),
		)
	}

	if !m.PrivateKey.IsNull() {
		addPathError(
			&diags,
			path.Root("private_key"),
			"Private key cannot be parsed",
			c.SetPrivateKey(m.PrivateKey.ValueString()),
		)
	}

	if !m.CARef.IsNull() {
		addPathError(
			&diags,
			path.Root("caref"),
			"CA reference cannot be parsed",
			c.SetCARef(m.CARef.ValueString()),
		)
	}

	return diags
}
