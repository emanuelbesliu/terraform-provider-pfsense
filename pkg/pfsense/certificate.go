package pfsense

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

type certificateResponse struct {
	RefID         string `json:"refid"`
	Descr         string `json:"descr"`
	CertType      string `json:"type"`
	CARef         string `json:"caref"`
	Subject       string `json:"subject"`
	Issuer        string `json:"issuer"`
	Serial        string `json:"serial"`
	HasPrivateKey bool   `json:"has_private_key"`
	IsSelfSigned  bool   `json:"is_self_signed"`
	ValidFrom     string `json:"valid_from"`
	ValidTo       string `json:"valid_to"`
	InUse         bool   `json:"in_use"`
	Certificate   string `json:"certificate"`
}

type certificateMutationResponse struct {
	Success bool    `json:"success"`
	Error   *string `json:"error"`
	RefID   *string `json:"refid"`
}

type Certificate struct {
	RefID         string
	Descr         string
	CertType      string
	CARef         string
	Certificate   string
	PrivateKey    string
	Subject       string
	Issuer        string
	Serial        string
	HasPrivateKey bool
	IsSelfSigned  bool
	ValidFrom     string
	ValidTo       string
	InUse         bool
}

func (Certificate) CertTypes() []string {
	return []string{"server", "user"}
}

func (c *Certificate) SetDescr(descr string) error {
	if descr == "" {
		return fmt.Errorf("%w, certificate description is required", ErrClientValidation)
	}

	c.Descr = descr

	return nil
}

func (c *Certificate) SetCertType(certType string) error {
	for _, ct := range c.CertTypes() {
		if ct == certType {
			c.CertType = certType

			return nil
		}
	}

	return fmt.Errorf("%w, invalid certificate type '%s'", ErrClientValidation, certType)
}

func (c *Certificate) SetCertificate(cert string) error {
	if cert == "" {
		return fmt.Errorf("%w, certificate is required", ErrClientValidation)
	}

	if !strings.Contains(cert, "BEGIN CERTIFICATE") {
		return fmt.Errorf("%w, certificate does not appear to be in PEM format", ErrClientValidation)
	}

	c.Certificate = cert

	return nil
}

func (c *Certificate) SetPrivateKey(key string) error {
	c.PrivateKey = key

	return nil
}

func (c *Certificate) SetCARef(caref string) error {
	c.CARef = caref

	return nil
}

type Certificates []Certificate

func (certs Certificates) GetByRefID(refid string) (*Certificate, error) {
	for _, c := range certs {
		if c.RefID == refid {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("certificate %w with refid '%s'", ErrNotFound, refid)
}

func (certs Certificates) GetByDescr(descr string) (*Certificate, error) {
	for _, c := range certs {
		if c.Descr == descr {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("certificate %w with description '%s'", ErrNotFound, descr)
}

func parseCertificateResponse(resp certificateResponse) Certificate {
	return Certificate{
		RefID:         resp.RefID,
		Descr:         resp.Descr,
		CertType:      resp.CertType,
		CARef:         resp.CARef,
		Certificate:   resp.Certificate,
		Subject:       resp.Subject,
		Issuer:        resp.Issuer,
		Serial:        resp.Serial,
		HasPrivateKey: resp.HasPrivateKey,
		IsSelfSigned:  resp.IsSelfSigned,
		ValidFrom:     resp.ValidFrom,
		ValidTo:       resp.ValidTo,
		InUse:         resp.InUse,
	}
}

func (pf *Client) getCertificates(ctx context.Context) (*Certificates, error) {
	command := "require_once('config.inc');" +
		"require_once('certs.inc');" +
		"$old_err = error_reporting(0);" +
		"$certs = array();" +
		"foreach (config_get_path('cert', array()) as $idx => $cert) {" +
		"if (!is_array($cert) || empty($cert)) { continue; }" +
		"$crt_details = openssl_x509_parse(base64_decode($cert['crt']));" +
		"$subject = cert_get_subject($cert['crt'], true);" +
		"$issuer = cert_get_issuer($cert['crt'], true);" +
		"$item = array();" +
		"$item['refid'] = $cert['refid'];" +
		"$item['descr'] = $cert['descr'];" +
		"$item['type'] = $cert['type'] ?? 'server';" +
		"$item['caref'] = $cert['caref'] ?? '';" +
		"$item['subject'] = $subject;" +
		"$item['issuer'] = $issuer;" +
		"$item['serial'] = (string)cert_get_serial($cert['crt'], true);" +
		"$item['has_private_key'] = !empty($cert['prv']);" +
		"$item['is_self_signed'] = ($subject === $issuer);" +
		"$item['valid_from'] = $crt_details['validFrom'] ?? '';" +
		"$item['valid_to'] = $crt_details['validTo'] ?? '';" +
		"$item['in_use'] = cert_in_use($cert['refid']);" +
		"$item['certificate'] = base64_decode($cert['crt']);" +
		"array_push($certs, $item);" +
		"};" +
		"error_reporting($old_err);" +
		"print(json_encode($certs));"

	var certResp []certificateResponse
	if err := pf.executePHPCommand(ctx, command, &certResp); err != nil {
		return nil, err
	}

	certs := make(Certificates, 0, len(certResp))
	for _, resp := range certResp {
		certs = append(certs, parseCertificateResponse(resp))
	}

	return &certs, nil
}

func (pf *Client) GetCertificates(ctx context.Context) (*Certificates, error) {
	defer pf.read(&pf.mutexes.Certificate)()

	certs, err := pf.getCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificates, %w", ErrGetOperationFailed, err)
	}

	return certs, nil
}

func (pf *Client) GetCertificate(ctx context.Context, refid string) (*Certificate, error) {
	defer pf.read(&pf.mutexes.Certificate)()

	certs, err := pf.getCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificates, %w", ErrGetOperationFailed, err)
	}

	cert, err := certs.GetByRefID(refid)
	if err != nil {
		return nil, fmt.Errorf("%w certificate, %w", ErrGetOperationFailed, err)
	}

	return cert, nil
}

func (pf *Client) GetCertificateByDescr(ctx context.Context, descr string) (*Certificate, error) {
	defer pf.read(&pf.mutexes.Certificate)()

	certs, err := pf.getCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificates, %w", ErrGetOperationFailed, err)
	}

	cert, err := certs.GetByDescr(descr)
	if err != nil {
		return nil, fmt.Errorf("%w certificate, %w", ErrGetOperationFailed, err)
	}

	return cert, nil
}

func (pf *Client) ImportCertificate(ctx context.Context, req Certificate) (*Certificate, error) {
	defer pf.write(&pf.mutexes.Certificate)()

	certB64 := base64.StdEncoding.EncodeToString([]byte(req.Certificate))
	keyB64 := base64.StdEncoding.EncodeToString([]byte(req.PrivateKey))

	command := fmt.Sprintf(
		"require_once('config.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'refid' => null, 'error' => null);"+
			"$descr = '%s';"+
			"$cert_pem = base64_decode('%s');"+
			"$key_pem = base64_decode('%s');"+
			"$cert_type = '%s';"+
			"$caref = '%s';"+
			"if (!strstr($cert_pem, 'BEGIN CERTIFICATE') || !strstr($cert_pem, 'END CERTIFICATE')) {"+
			"$result['error'] = 'Certificate is not in PEM format';"+
			"print(json_encode($result)); return;"+
			"}"+
			"foreach (config_get_path('cert', array()) as $existing_cert) {"+
			"if (!is_array($existing_cert) || empty($existing_cert)) { continue; }"+
			"if ($existing_cert['descr'] === $descr) {"+
			"$result['error'] = 'Certificate with name \"' . $descr . '\" already exists (refid: ' . $existing_cert['refid'] . ')';"+
			"print(json_encode($result)); return;"+
			"}"+
			"}"+
			"$cert = array('refid' => uniqid(), 'descr' => $descr, 'type' => $cert_type);"+
			"$old_err = error_reporting(0);"+
			"cert_import($cert, $cert_pem, $key_pem);"+
			"error_reporting($old_err);"+
			"if (!empty($caref)) { $cert['caref'] = $caref; }"+
			"config_set_path('cert/', $cert);"+
			"write_config('Terraform: imported certificate: ' . $descr);"+
			"$result['success'] = true;"+
			"$result['refid'] = $cert['refid'];"+
			"print(json_encode($result));",
		phpEscape(req.Descr),
		certB64,
		keyB64,
		phpEscape(req.CertType),
		phpEscape(req.CARef),
	)

	var result certificateMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w certificate, %w", ErrCreateOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return nil, fmt.Errorf("%w certificate '%s', %s", ErrCreateOperationFailed, req.Descr, errMsg)
	}

	if result.RefID == nil {
		return nil, fmt.Errorf("%w certificate '%s', no refid returned", ErrCreateOperationFailed, req.Descr)
	}

	certs, err := pf.getCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificates after importing, %w", ErrGetOperationFailed, err)
	}

	cert, err := certs.GetByRefID(*result.RefID)
	if err != nil {
		return nil, fmt.Errorf("%w certificate after importing '%s', %w", ErrGetOperationFailed, req.Descr, err)
	}

	return cert, nil
}

func (pf *Client) UpdateCertificate(ctx context.Context, refid string, req Certificate) (*Certificate, error) {
	defer pf.write(&pf.mutexes.Certificate)()

	command := fmt.Sprintf(
		"require_once('config.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'error' => null);"+
			"$refid = '%s';"+
			"$cert_lookup = lookup_cert($refid);"+
			"if (!$cert_lookup['item']) {"+
			"$result['error'] = 'Certificate not found';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$cert_idx = $cert_lookup['idx'];"+
			"$cert = $cert_lookup['item'];"+
			"$cert['descr'] = '%s';"+
			"$cert['type'] = '%s';"+
			"$caref = '%s';"+
			"if (!empty($caref)) { $cert['caref'] = $caref; } else { unset($cert['caref']); }"+
			"config_set_path('cert/' . $cert_idx, $cert);"+
			"write_config('Terraform: updated certificate: ' . $cert['descr']);"+
			"$result['success'] = true;"+
			"print(json_encode($result));",
		phpEscape(refid),
		phpEscape(req.Descr),
		phpEscape(req.CertType),
		phpEscape(req.CARef),
	)

	var result certificateMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w certificate, %w", ErrUpdateOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return nil, fmt.Errorf("%w certificate '%s', %s", ErrUpdateOperationFailed, refid, errMsg)
	}

	certs, err := pf.getCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificates after updating, %w", ErrGetOperationFailed, err)
	}

	cert, err := certs.GetByRefID(refid)
	if err != nil {
		return nil, fmt.Errorf("%w certificate after updating, %w", ErrGetOperationFailed, err)
	}

	return cert, nil
}

func (pf *Client) DeleteCertificate(ctx context.Context, refid string) error {
	defer pf.write(&pf.mutexes.Certificate)()

	command := fmt.Sprintf(
		"require_once('config.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'error' => null);"+
			"$refid = '%s';"+
			"$cert_lookup = lookup_cert($refid);"+
			"if (!$cert_lookup['item']) {"+
			"$result['error'] = 'Certificate not found';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$cert_idx = $cert_lookup['idx'];"+
			"$cert = $cert_lookup['item'];"+
			"if (cert_in_use($refid)) {"+
			"$result['error'] = 'Certificate is in use and cannot be deleted';"+
			"print(json_encode($result)); return;"+
			"}"+
			"config_del_path('cert/' . $cert_idx);"+
			"write_config('Terraform: deleted certificate: ' . $cert['descr']);"+
			"$result['success'] = true;"+
			"print(json_encode($result));",
		phpEscape(refid),
	)

	var result certificateMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w certificate, %w", ErrDeleteOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return fmt.Errorf("%w certificate '%s', %s", ErrDeleteOperationFailed, refid, errMsg)
	}

	return nil
}
