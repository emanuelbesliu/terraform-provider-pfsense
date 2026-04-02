package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ============================================================================
// Response types (JSON from PHP config read)
// ============================================================================

type ipsecPhase1EncryptionItemResponse struct {
	EncryptionAlgorithm ipsecEncryptionAlgorithmResponse `json:"encryption-algorithm"`
	HashAlgorithm       string                           `json:"hash-algorithm"`
	PRFAlgorithm        string                           `json:"prf-algorithm"`
	DHGroup             string                           `json:"dhgroup"`
}

type ipsecEncryptionAlgorithmResponse struct {
	Name   string `json:"name"`
	KeyLen string `json:"keylen"`
}

type ipsecPhase1EncryptionResponse struct {
	Item ipsecPhase1EncryptionItems `json:"item"`
}

// ipsecPhase1EncryptionItems handles pfSense returning either a single object
// or an array for the encryption items.
type ipsecPhase1EncryptionItems []ipsecPhase1EncryptionItemResponse

func (p *ipsecPhase1EncryptionItems) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []ipsecPhase1EncryptionItemResponse
	if err := json.Unmarshal(data, &arr); err == nil {
		*p = arr
		return nil
	}

	// Try single object
	var single ipsecPhase1EncryptionItemResponse
	if err := json.Unmarshal(data, &single); err == nil {
		*p = []ipsecPhase1EncryptionItemResponse{single}
		return nil
	}

	return fmt.Errorf("unable to unmarshal encryption items")
}

type ipsecPhase1Response struct {
	IKEId                string                        `json:"ikeid"`
	IKEType              string                        `json:"iketype"`
	Interface            string                        `json:"interface"`
	Protocol             string                        `json:"protocol"`
	RemoteGateway        string                        `json:"remote-gateway"`
	AuthenticationMethod string                        `json:"authentication_method"`
	PreSharedKey         string                        `json:"pre-shared-key"`
	MyIDType             string                        `json:"myid_type"`
	MyIDData             string                        `json:"myid_data"`
	PeerIDType           string                        `json:"peerid_type"`
	PeerIDData           string                        `json:"peerid_data"`
	Description          string                        `json:"descr"`
	NATTraversal         string                        `json:"nat_traversal"`
	Mobike               string                        `json:"mobike"`
	DPDDelay             string                        `json:"dpd_delay"`
	DPDMaxFail           string                        `json:"dpd_maxfail"`
	Lifetime             string                        `json:"lifetime"`
	RekeyTime            string                        `json:"rekey_time"`
	ReauthTime           string                        `json:"reauth_time"`
	RandTime             string                        `json:"rand_time"`
	StartAction          string                        `json:"startaction"`
	CloseAction          string                        `json:"closeaction"`
	Encryption           ipsecPhase1EncryptionResponse `json:"encryption"`
	Disabled             *string                       `json:"disabled"`
	CertRef              string                        `json:"certref"`
	CARef                string                        `json:"caref"`
	PKCS11CertRef        string                        `json:"pkcs11certref"`
	PKCS11Pin            string                        `json:"pkcs11pin"`
	PrivateKey           string                        `json:"private-key"`
	Mobile               *string                       `json:"mobile"`
	IKEPort              string                        `json:"ikeport"`
	NATTPort             string                        `json:"nattport"`
	GWDuplicates         *string                       `json:"gw_duplicates"`
	PRFSelectEnable      *string                       `json:"prfselect_enable"`
	SplitConn            *string                       `json:"splitconn"`
	TFCEnable            *string                       `json:"tfc_enable"`
	TFCBytes             string                        `json:"tfc_bytes"`
}

// ============================================================================
// Domain types
// ============================================================================

type IPsecPhase1EncryptionItem struct {
	Algorithm    string
	KeyLen       string
	HashAlgo     string
	PRFAlgorithm string
	DHGroup      string
}

type IPsecPhase1 struct {
	IKEId                string
	IKEType              string
	Interface            string
	Protocol             string
	RemoteGateway        string
	AuthenticationMethod string
	PreSharedKey         string
	MyIDType             string
	MyIDData             string
	PeerIDType           string
	PeerIDData           string
	Description          string
	NATTraversal         string
	Mobike               string
	DPDDelay             string
	DPDMaxFail           string
	Lifetime             string
	RekeyTime            string
	ReauthTime           string
	RandTime             string
	StartAction          string
	CloseAction          string
	Encryption           []IPsecPhase1EncryptionItem
	Disabled             bool
	CertRef              string
	CARef                string
	PKCS11CertRef        string
	PKCS11Pin            string
	Mobile               bool
	IKEPort              string
	NATTPort             string
	GWDuplicates         bool
	PRFSelectEnable      bool
	SplitConn            bool
	TFCEnable            bool
	TFCBytes             string
}

type IPsecPhase1s []IPsecPhase1

func (ps IPsecPhase1s) GetByIKEId(ikeId string) (*IPsecPhase1, error) {
	for _, p := range ps {
		if p.IKEId == ikeId {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("ipsec phase1 %w with ikeid '%s'", ErrNotFound, ikeId)
}

func (ps IPsecPhase1s) GetIndexByIKEId(ikeId string) (*int, error) {
	for index, p := range ps {
		if p.IKEId == ikeId {
			return &index, nil
		}
	}

	return nil, fmt.Errorf("ipsec phase1 %w with ikeid '%s'", ErrNotFound, ikeId)
}

// ============================================================================
// Parsing
// ============================================================================

func parseIPsecPhase1Response(resp ipsecPhase1Response) IPsecPhase1 {
	p1 := IPsecPhase1{
		IKEId:                resp.IKEId,
		IKEType:              resp.IKEType,
		Interface:            resp.Interface,
		Protocol:             resp.Protocol,
		RemoteGateway:        resp.RemoteGateway,
		AuthenticationMethod: resp.AuthenticationMethod,
		PreSharedKey:         resp.PreSharedKey,
		MyIDType:             resp.MyIDType,
		MyIDData:             resp.MyIDData,
		PeerIDType:           resp.PeerIDType,
		PeerIDData:           resp.PeerIDData,
		Description:          resp.Description,
		NATTraversal:         resp.NATTraversal,
		Mobike:               resp.Mobike,
		DPDDelay:             resp.DPDDelay,
		DPDMaxFail:           resp.DPDMaxFail,
		Lifetime:             resp.Lifetime,
		RekeyTime:            resp.RekeyTime,
		ReauthTime:           resp.ReauthTime,
		RandTime:             resp.RandTime,
		StartAction:          resp.StartAction,
		CloseAction:          resp.CloseAction,
		Disabled:             resp.Disabled != nil,
		CertRef:              resp.CertRef,
		CARef:                resp.CARef,
		PKCS11CertRef:        resp.PKCS11CertRef,
		PKCS11Pin:            resp.PKCS11Pin,
		Mobile:               resp.Mobile != nil,
		IKEPort:              resp.IKEPort,
		NATTPort:             resp.NATTPort,
		GWDuplicates:         resp.GWDuplicates != nil,
		PRFSelectEnable:      resp.PRFSelectEnable != nil,
		SplitConn:            resp.SplitConn != nil,
		TFCEnable:            resp.TFCEnable != nil,
		TFCBytes:             resp.TFCBytes,
	}

	for _, item := range resp.Encryption.Item {
		p1.Encryption = append(p1.Encryption, IPsecPhase1EncryptionItem{
			Algorithm:    item.EncryptionAlgorithm.Name,
			KeyLen:       item.EncryptionAlgorithm.KeyLen,
			HashAlgo:     item.HashAlgorithm,
			PRFAlgorithm: item.PRFAlgorithm,
			DHGroup:      item.DHGroup,
		})
	}

	return p1
}

// ============================================================================
// Form values for POST
// ============================================================================

func ipsecPhase1FormValues(p1 IPsecPhase1) url.Values {
	values := url.Values{
		"save":                  {"Save"},
		"iketype":               {p1.IKEType},
		"interface":             {p1.Interface},
		"protocol":              {p1.Protocol},
		"authentication_method": {p1.AuthenticationMethod},
		"descr":                 {p1.Description},
		"nat_traversal":         {p1.NATTraversal},
		"mobike":                {p1.Mobike},
		"lifetime":              {p1.Lifetime},
		"myid_type":             {p1.MyIDType},
		"myid_data":             {p1.MyIDData},
		"peerid_type":           {p1.PeerIDType},
		"peerid_data":           {p1.PeerIDData},
	}

	// IKE ID (for updates)
	if p1.IKEId != "" {
		values.Set("ikeid", p1.IKEId)
	}

	// Remote gateway (only for non-mobile)
	if !p1.Mobile {
		values.Set("remotegw", p1.RemoteGateway)
	} else {
		values.Set("mobile", "true")
	}

	// Pre-shared key
	if p1.AuthenticationMethod == "pre_shared_key" || p1.AuthenticationMethod == "xauth_psk_server" {
		values.Set("pskey", p1.PreSharedKey)
	}

	// Certificate-based authentication
	if p1.CertRef != "" {
		values.Set("certref", p1.CertRef)
	}
	if p1.CARef != "" {
		values.Set("caref", p1.CARef)
	}
	if p1.PKCS11CertRef != "" {
		values.Set("pkcs11certref", p1.PKCS11CertRef)
	}
	if p1.PKCS11Pin != "" {
		values.Set("pkcs11pin", p1.PKCS11Pin)
	}

	// Timers
	if p1.RekeyTime != "" {
		values.Set("rekey_time", p1.RekeyTime)
	}
	if p1.ReauthTime != "" {
		values.Set("reauth_time", p1.ReauthTime)
	}
	if p1.RandTime != "" {
		values.Set("rand_time", p1.RandTime)
	}

	// Start/Close actions
	if p1.StartAction != "" {
		values.Set("startaction", p1.StartAction)
	}
	if p1.CloseAction != "" {
		values.Set("closeaction", p1.CloseAction)
	}

	// DPD
	if p1.DPDDelay != "" && p1.DPDMaxFail != "" {
		values.Set("dpd_enable", "yes")
		values.Set("dpd_delay", p1.DPDDelay)
		values.Set("dpd_maxfail", p1.DPDMaxFail)
	}

	// IKE/NAT-T ports
	if p1.IKEPort != "" {
		values.Set("ikeport", p1.IKEPort)
	}
	if p1.NATTPort != "" {
		values.Set("nattport", p1.NATTPort)
	}

	// Boolean options (presence-based)
	if p1.Disabled {
		values.Set("disabled", "yes")
	}
	if p1.GWDuplicates {
		values.Set("gw_duplicates", "yes")
	}
	if p1.PRFSelectEnable {
		values.Set("prfselect_enable", "yes")
	}
	if p1.SplitConn {
		values.Set("splitconn", "yes")
	}
	if p1.TFCEnable {
		values.Set("tfc_enable", "yes")
	}
	if p1.TFCBytes != "" {
		values.Set("tfc_bytes", p1.TFCBytes)
	}

	// Encryption algorithms (indexed fields)
	for i, enc := range p1.Encryption {
		values.Set(fmt.Sprintf("ealgo_algo%d", i), enc.Algorithm)
		values.Set(fmt.Sprintf("ealgo_keylen%d", i), enc.KeyLen)
		values.Set(fmt.Sprintf("halgo%d", i), enc.HashAlgo)
		values.Set(fmt.Sprintf("prfalgo%d", i), enc.PRFAlgorithm)
		values.Set(fmt.Sprintf("dhgroup%d", i), enc.DHGroup)
	}

	return values
}

// ============================================================================
// Client methods
// ============================================================================

func (pf *Client) getIPsecPhase1s(ctx context.Context) (*IPsecPhase1s, error) {
	command := `
$phase1 = config_get_path('ipsec/phase1', array());
if (!is_array($phase1)) { $phase1 = array(); }
// Handle single entry (not wrapped in array)
if (isset($phase1['ikeid'])) { $phase1 = array($phase1); }
print(json_encode($phase1));
`
	var resp []ipsecPhase1Response
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	phase1s := make(IPsecPhase1s, 0, len(resp))
	for _, r := range resp {
		phase1s = append(phase1s, parseIPsecPhase1Response(r))
	}

	return &phase1s, nil
}

func (pf *Client) GetIPsecPhase1s(ctx context.Context) (*IPsecPhase1s, error) {
	defer pf.read(&pf.mutexes.IPsecPhase1)()

	phase1s, err := pf.getIPsecPhase1s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1 entries, %w", ErrGetOperationFailed, err)
	}

	return phase1s, nil
}

func (pf *Client) GetIPsecPhase1(ctx context.Context, ikeId string) (*IPsecPhase1, error) {
	defer pf.read(&pf.mutexes.IPsecPhase1)()

	phase1s, err := pf.getIPsecPhase1s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1 entries, %w", ErrGetOperationFailed, err)
	}

	p1, err := phase1s.GetByIKEId(ikeId)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1, %w", ErrGetOperationFailed, err)
	}

	return p1, nil
}

func (pf *Client) createOrUpdateIPsecPhase1(ctx context.Context, p1 IPsecPhase1, ikeId *string) error {
	relativeURL := url.URL{Path: "vpn_ipsec_phase1.php"}

	if ikeId != nil {
		q := relativeURL.Query()
		q.Set("ikeid", *ikeId)
		relativeURL.RawQuery = q.Encode()
	}

	values := ipsecPhase1FormValues(p1)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateIPsecPhase1(ctx context.Context, p1Req IPsecPhase1) (*IPsecPhase1, error) {
	defer pf.write(&pf.mutexes.IPsecPhase1)()

	if err := pf.createOrUpdateIPsecPhase1(ctx, p1Req, nil); err != nil {
		return nil, fmt.Errorf("%w ipsec phase1, %w", ErrCreateOperationFailed, err)
	}

	phase1s, err := pf.getIPsecPhase1s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1 entries after creating, %w", ErrGetOperationFailed, err)
	}

	// Find the newly created entry (last one with matching description)
	var found *IPsecPhase1
	for i := len(*phase1s) - 1; i >= 0; i-- {
		p := (*phase1s)[i]
		if p.Description == p1Req.Description &&
			p.RemoteGateway == p1Req.RemoteGateway {
			found = &p

			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("%w ipsec phase1 after creating, could not find newly created entry", ErrGetOperationFailed)
	}

	return found, nil
}

func (pf *Client) UpdateIPsecPhase1(ctx context.Context, p1Req IPsecPhase1) (*IPsecPhase1, error) {
	defer pf.write(&pf.mutexes.IPsecPhase1)()

	if err := pf.createOrUpdateIPsecPhase1(ctx, p1Req, &p1Req.IKEId); err != nil {
		return nil, fmt.Errorf("%w ipsec phase1, %w", ErrUpdateOperationFailed, err)
	}

	phase1s, err := pf.getIPsecPhase1s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1 entries after updating, %w", ErrGetOperationFailed, err)
	}

	p1, err := phase1s.GetByIKEId(p1Req.IKEId)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase1 after updating, %w", ErrGetOperationFailed, err)
	}

	return p1, nil
}

func (pf *Client) deleteIPsecPhase1(ctx context.Context, index int) error {
	relativeURL := url.URL{Path: "vpn_ipsec.php"}
	values := url.Values{
		"del":       {"Delete selected Phase 1 entries"},
		"p1entry[]": {strconv.Itoa(index)},
	}

	_, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)

	return err
}

func (pf *Client) DeleteIPsecPhase1(ctx context.Context, ikeId string) error {
	defer pf.write(&pf.mutexes.IPsecPhase1)()

	phase1s, err := pf.getIPsecPhase1s(ctx)
	if err != nil {
		return fmt.Errorf("%w ipsec phase1 entries, %w", ErrGetOperationFailed, err)
	}

	index, err := phase1s.GetIndexByIKEId(ikeId)
	if err != nil {
		return fmt.Errorf("%w ipsec phase1, %w", ErrGetOperationFailed, err)
	}

	if err := pf.deleteIPsecPhase1(ctx, *index); err != nil {
		return fmt.Errorf("%w ipsec phase1, %w", ErrDeleteOperationFailed, err)
	}

	phase1s, err = pf.getIPsecPhase1s(ctx)
	if err != nil {
		return fmt.Errorf("%w ipsec phase1 entries after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := phase1s.GetByIKEId(ikeId); err == nil {
		return fmt.Errorf("%w ipsec phase1, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
