package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type certificateAuthorityResponse struct {
	RefID         string `json:"refid"`
	Descr         string `json:"descr"`
	Subject       string `json:"subject"`
	Issuer        string `json:"issuer"`
	Serial        string `json:"serial"`
	HasPrivateKey bool   `json:"has_private_key"`
	IsSelfSigned  bool   `json:"is_self_signed"`
	TrustEnabled  bool   `json:"trust_enabled"`
	RandomSerial  bool   `json:"randomserial"`
	NextSerial    int    `json:"next_serial"`
	ValidFrom     string `json:"valid_from"`
	ValidTo       string `json:"valid_to"`
	InUse         bool   `json:"in_use"`
	Certificate   string `json:"certificate"`
}

type certificateAuthorityMutationResponse struct {
	Success bool    `json:"success"`
	Error   *string `json:"error"`
	RefID   *string `json:"refid"`
}

type CertificateAuthority struct {
	RefID         string
	Descr         string
	Certificate   string
	PrivateKey    string
	Subject       string
	Issuer        string
	Serial        string
	HasPrivateKey bool
	IsSelfSigned  bool
	Trust         bool
	RandomSerial  bool
	NextSerial    int
	ValidFrom     string
	ValidTo       string
	InUse         bool
}

func (ca *CertificateAuthority) SetDescr(descr string) error {
	if descr == "" {
		return fmt.Errorf("%w, CA description is required", ErrClientValidation)
	}

	ca.Descr = descr

	return nil
}

func (ca *CertificateAuthority) SetCertificate(cert string) error {
	if cert == "" {
		return fmt.Errorf("%w, CA certificate is required", ErrClientValidation)
	}

	if !strings.Contains(cert, "BEGIN CERTIFICATE") {
		return fmt.Errorf("%w, certificate does not appear to be in PEM format", ErrClientValidation)
	}

	ca.Certificate = cert

	return nil
}

func (ca *CertificateAuthority) SetPrivateKey(key string) error {
	ca.PrivateKey = key

	return nil
}

func (ca *CertificateAuthority) SetTrust(trust bool) error {
	ca.Trust = trust

	return nil
}

func (ca *CertificateAuthority) SetNextSerial(serial int) error {
	ca.NextSerial = serial

	return nil
}

type CertificateAuthorities []CertificateAuthority

func (cas CertificateAuthorities) GetByRefID(refid string) (*CertificateAuthority, error) {
	for _, ca := range cas {
		if ca.RefID == refid {
			return &ca, nil
		}
	}

	return nil, fmt.Errorf("certificate authority %w with refid '%s'", ErrNotFound, refid)
}

func (cas CertificateAuthorities) GetByDescr(descr string) (*CertificateAuthority, error) {
	for _, ca := range cas {
		if ca.Descr == descr {
			return &ca, nil
		}
	}

	return nil, fmt.Errorf("certificate authority %w with description '%s'", ErrNotFound, descr)
}

func parseCertificateAuthorityResponse(resp certificateAuthorityResponse) CertificateAuthority {
	return CertificateAuthority{
		RefID:         resp.RefID,
		Descr:         resp.Descr,
		Certificate:   resp.Certificate,
		Subject:       resp.Subject,
		Issuer:        resp.Issuer,
		Serial:        resp.Serial,
		HasPrivateKey: resp.HasPrivateKey,
		IsSelfSigned:  resp.IsSelfSigned,
		Trust:         resp.TrustEnabled,
		RandomSerial:  resp.RandomSerial,
		NextSerial:    resp.NextSerial,
		ValidFrom:     resp.ValidFrom,
		ValidTo:       resp.ValidTo,
		InUse:         resp.InUse,
	}
}

func (pf *Client) getCertificateAuthorities(ctx context.Context) (*CertificateAuthorities, error) {
	command := "require_once('guiconfig.inc');" +
		"require_once('certs.inc');" +
		"$cas = array();" +
		"foreach (config_get_path('ca', array()) as $idx => $ca) {" +
		"if (!is_array($ca) || empty($ca)) { continue; }" +
		"$crt_details = openssl_x509_parse(base64_decode($ca['crt']));" +
		"$subject = cert_get_subject($ca['crt'], true);" +
		"$issuer = cert_get_issuer($ca['crt'], true);" +
		"$item = array();" +
		"$item['refid'] = $ca['refid'];" +
		"$item['descr'] = $ca['descr'];" +
		"$item['subject'] = $subject;" +
		"$item['issuer'] = $issuer;" +
		"$item['serial'] = (string)cert_get_serial($ca['crt'], true);" +
		"$item['has_private_key'] = !empty($ca['prv']);" +
		"$item['is_self_signed'] = ($subject === $issuer);" +
		"$item['trust_enabled'] = ($ca['trust'] === 'enabled');" +
		"$item['randomserial'] = ($ca['randomserial'] === 'enabled');" +
		"$item['next_serial'] = (int)($ca['serial'] ?? 0);" +
		"$item['valid_from'] = $crt_details['validFrom'] ?? '';" +
		"$item['valid_to'] = $crt_details['validTo'] ?? '';" +
		"$item['in_use'] = ca_in_use($ca['refid']);" +
		"$item['certificate'] = base64_decode($ca['crt']);" +
		"array_push($cas, $item);" +
		"};" +
		"print(json_encode($cas));"

	var caResp []certificateAuthorityResponse
	if err := pf.executePHPCommand(ctx, command, &caResp); err != nil {
		return nil, err
	}

	cas := make(CertificateAuthorities, 0, len(caResp))
	for _, resp := range caResp {
		cas = append(cas, parseCertificateAuthorityResponse(resp))
	}

	return &cas, nil
}

func (pf *Client) GetCertificateAuthorities(ctx context.Context) (*CertificateAuthorities, error) {
	defer pf.read(&pf.mutexes.CertificateAuthority)()

	cas, err := pf.getCertificateAuthorities(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authorities, %w", ErrGetOperationFailed, err)
	}

	return cas, nil
}

func (pf *Client) GetCertificateAuthority(ctx context.Context, refid string) (*CertificateAuthority, error) {
	defer pf.read(&pf.mutexes.CertificateAuthority)()

	cas, err := pf.getCertificateAuthorities(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authorities, %w", ErrGetOperationFailed, err)
	}

	ca, err := cas.GetByRefID(refid)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authority, %w", ErrGetOperationFailed, err)
	}

	return ca, nil
}

func (pf *Client) GetCertificateAuthorityByDescr(ctx context.Context, descr string) (*CertificateAuthority, error) {
	defer pf.read(&pf.mutexes.CertificateAuthority)()

	cas, err := pf.getCertificateAuthorities(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authorities, %w", ErrGetOperationFailed, err)
	}

	ca, err := cas.GetByDescr(descr)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authority, %w", ErrGetOperationFailed, err)
	}

	return ca, nil
}

func (pf *Client) ImportCertificateAuthority(ctx context.Context, req CertificateAuthority) (*CertificateAuthority, error) {
	defer pf.write(&pf.mutexes.CertificateAuthority)()

	trustValue := "disabled"
	if req.Trust {
		trustValue = "enabled"
	}

	command := fmt.Sprintf(
		"require_once('guiconfig.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'refid' => null, 'error' => null);"+
			"$descr = '%s';"+
			"$cert_pem = '%s';"+
			"$key_pem = '%s';"+
			"$next_serial = %d;"+
			"if (!strstr($cert_pem, 'BEGIN CERTIFICATE') || !strstr($cert_pem, 'END CERTIFICATE')) {"+
			"$result['error'] = 'Certificate is not in PEM format';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$purpose = cert_get_purpose($cert_pem, false);"+
			"if ($purpose['ca'] !== 'Yes') {"+
			"$result['error'] = 'Certificate does not have CA constraints';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$serial_number = cert_get_serial($cert_pem, false);"+
			"foreach (config_get_path('ca', array()) as $existing_ca) {"+
			"if (!is_array($existing_ca) || empty($existing_ca)) { continue; }"+
			"$existing_serial = cert_get_serial($existing_ca['crt'], true);"+
			"if ($serial_number === $existing_serial) {"+
			"$result['error'] = 'CA with this serial already exists (refid: ' . $existing_ca['refid'] . ')';"+
			"print(json_encode($result)); return;"+
			"}"+
			"}"+
			"foreach (config_get_path('ca', array()) as $existing_ca) {"+
			"if (!is_array($existing_ca) || empty($existing_ca)) { continue; }"+
			"if ($existing_ca['descr'] === $descr) {"+
			"$result['error'] = 'CA with name \"' . $descr . '\" already exists (refid: ' . $existing_ca['refid'] . ')';"+
			"print(json_encode($result)); return;"+
			"}"+
			"}"+
			"$ca = array('refid' => uniqid(), 'descr' => $descr);"+
			"$old_err = error_reporting(0);"+
			"$ok = ca_import($ca, $cert_pem, $key_pem, $next_serial);"+
			"error_reporting($old_err);"+
			"if (!$ok) {"+
			"$result['error'] = 'ca_import() failed';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$ca['trust'] = '%s';"+
			"config_set_path('ca/', $ca);"+
			"write_config('Terraform: imported CA: ' . $descr);"+
			"ca_setup_trust_store();"+
			"$result['success'] = true;"+
			"$result['refid'] = $ca['refid'];"+
			"print(json_encode($result));",
		phpEscape(req.Descr),
		phpEscape(req.Certificate),
		phpEscape(req.PrivateKey),
		req.NextSerial,
		phpEscape(trustValue),
	)

	var result certificateAuthorityMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w certificate authority, %w", ErrCreateOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return nil, fmt.Errorf("%w certificate authority '%s', %s", ErrCreateOperationFailed, req.Descr, errMsg)
	}

	if result.RefID == nil {
		return nil, fmt.Errorf("%w certificate authority '%s', no refid returned", ErrCreateOperationFailed, req.Descr)
	}

	cas, err := pf.getCertificateAuthorities(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authorities after importing, %w", ErrGetOperationFailed, err)
	}

	ca, err := cas.GetByRefID(*result.RefID)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authority after importing '%s', %w", ErrGetOperationFailed, req.Descr, err)
	}

	return ca, nil
}

func (pf *Client) UpdateCertificateAuthority(ctx context.Context, refid string, req CertificateAuthority) (*CertificateAuthority, error) {
	defer pf.write(&pf.mutexes.CertificateAuthority)()

	trustValue := "disabled"
	if req.Trust {
		trustValue = "enabled"
	}

	randomSerialValue := "disabled"
	if req.RandomSerial {
		randomSerialValue = "enabled"
	}

	command := fmt.Sprintf(
		"require_once('guiconfig.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'error' => null);"+
			"$refid = '%s';"+
			"$ca_lookup = lookup_ca($refid);"+
			"if (!$ca_lookup['item']) {"+
			"$result['error'] = 'CA not found';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$ca_idx = $ca_lookup['idx'];"+
			"$ca = $ca_lookup['item'];"+
			"$ca['descr'] = '%s';"+
			"$ca['trust'] = '%s';"+
			"$ca['randomserial'] = '%s';"+
			"config_set_path('ca/' . $ca_idx, $ca);"+
			"write_config('Terraform: updated CA: ' . $ca['descr']);"+
			"ca_setup_trust_store();"+
			"$result['success'] = true;"+
			"print(json_encode($result));",
		phpEscape(refid),
		phpEscape(req.Descr),
		phpEscape(trustValue),
		phpEscape(randomSerialValue),
	)

	var result certificateAuthorityMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w certificate authority, %w", ErrUpdateOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return nil, fmt.Errorf("%w certificate authority '%s', %s", ErrUpdateOperationFailed, refid, errMsg)
	}

	cas, err := pf.getCertificateAuthorities(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authorities after updating, %w", ErrGetOperationFailed, err)
	}

	ca, err := cas.GetByRefID(refid)
	if err != nil {
		return nil, fmt.Errorf("%w certificate authority after updating, %w", ErrGetOperationFailed, err)
	}

	return ca, nil
}

func (pf *Client) DeleteCertificateAuthority(ctx context.Context, refid string) error {
	defer pf.write(&pf.mutexes.CertificateAuthority)()

	command := fmt.Sprintf(
		"require_once('guiconfig.inc');"+
			"require_once('certs.inc');"+
			"$result = array('success' => false, 'error' => null);"+
			"$refid = '%s';"+
			"$ca_lookup = lookup_ca($refid);"+
			"if (!$ca_lookup['item']) {"+
			"$result['error'] = 'CA not found';"+
			"print(json_encode($result)); return;"+
			"}"+
			"$ca_idx = $ca_lookup['idx'];"+
			"$ca = $ca_lookup['item'];"+
			"if (ca_in_use($refid)) {"+
			"$result['error'] = 'CA is in use and cannot be deleted';"+
			"print(json_encode($result)); return;"+
			"}"+
			"foreach (config_get_path('cert', array()) as $cid => $cert) {"+
			"if ($cert['caref'] === $refid) { config_del_path('cert/' . $cid . '/caref'); }"+
			"}"+
			"foreach (config_get_path('ca', array()) as $cid => $ca_child) {"+
			"if ($ca_child['caref'] === $refid) { config_del_path('ca/' . $cid . '/caref'); }"+
			"}"+
			"foreach (config_get_path('crl', array()) as $cid => $crl) {"+
			"if ($crl['caref'] === $refid) { config_del_path('crl/' . $cid); }"+
			"}"+
			"config_del_path('ca/' . $ca_idx);"+
			"write_config('Terraform: deleted CA: ' . $ca['descr']);"+
			"ca_setup_trust_store();"+
			"$result['success'] = true;"+
			"print(json_encode($result));",
		phpEscape(refid),
	)

	var result certificateAuthorityMutationResponse
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w certificate authority, %w", ErrDeleteOperationFailed, err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}

		return fmt.Errorf("%w certificate authority '%s', %s", ErrDeleteOperationFailed, refid, errMsg)
	}

	return nil
}
