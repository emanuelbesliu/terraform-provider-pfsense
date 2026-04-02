package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ============================================================================
// Response types (JSON from PHP config read)
// ============================================================================

type ipsecPhase2EncryptionAlgorithmOptionResponse struct {
	Name   string `json:"name"`
	KeyLen string `json:"keylen"`
}

// ipsecPhase2EncryptionAlgorithmOptions handles pfSense returning either a
// single object or an array for the encryption-algorithm-option field.
type ipsecPhase2EncryptionAlgorithmOptions []ipsecPhase2EncryptionAlgorithmOptionResponse

func (p *ipsecPhase2EncryptionAlgorithmOptions) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []ipsecPhase2EncryptionAlgorithmOptionResponse
	if err := json.Unmarshal(data, &arr); err == nil {
		*p = arr
		return nil
	}

	// Try single object
	var single ipsecPhase2EncryptionAlgorithmOptionResponse
	if err := json.Unmarshal(data, &single); err == nil {
		*p = []ipsecPhase2EncryptionAlgorithmOptionResponse{single}
		return nil
	}

	return fmt.Errorf("unable to unmarshal phase2 encryption algorithm options")
}

// ipsecPhase2HashAlgorithms handles pfSense returning either a single string
// or an array of strings for the hash-algorithm-option field.
type ipsecPhase2HashAlgorithms []string

func (p *ipsecPhase2HashAlgorithms) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*p = arr
		return nil
	}

	// Try single string
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*p = []string{single}
		return nil
	}

	return fmt.Errorf("unable to unmarshal phase2 hash algorithms")
}

type ipsecPhase2IDResponse struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	NetBits string `json:"netbits"`
}

type ipsecPhase2Response struct {
	UniqID                    string                                `json:"uniqid"`
	IKEId                     string                                `json:"ikeid"`
	Mode                      string                                `json:"mode"`
	ReqID                     string                                `json:"reqid"`
	LocalID                   ipsecPhase2IDResponse                 `json:"localid"`
	RemoteID                  ipsecPhase2IDResponse                 `json:"remoteid"`
	NATLocalID                *ipsecPhase2IDResponse                `json:"natlocalid"`
	Protocol                  string                                `json:"protocol"`
	EncryptionAlgorithmOption ipsecPhase2EncryptionAlgorithmOptions `json:"encryption-algorithm-option"`
	HashAlgorithmOption       ipsecPhase2HashAlgorithms             `json:"hash-algorithm-option"`
	PFSGroup                  string                                `json:"pfsgroup"`
	Lifetime                  string                                `json:"lifetime"`
	RekeyTime                 string                                `json:"rekey_time"`
	RandTime                  string                                `json:"rand_time"`
	PingHost                  string                                `json:"pinghost"`
	Keepalive                 string                                `json:"keepalive"`
	Description               string                                `json:"descr"`
	Disabled                  *string                               `json:"disabled"`
	Mobile                    *string                               `json:"mobile"`
}

// ============================================================================
// Domain types
// ============================================================================

type IPsecPhase2EncryptionAlgorithm struct {
	Name   string
	KeyLen string
}

type IPsecPhase2ID struct {
	Type    string
	Address string
	NetBits string
}

type IPsecPhase2 struct {
	UniqID                    string
	IKEId                     string
	Mode                      string
	ReqID                     string
	LocalID                   IPsecPhase2ID
	RemoteID                  IPsecPhase2ID
	NATLocalID                *IPsecPhase2ID
	Protocol                  string
	EncryptionAlgorithmOption []IPsecPhase2EncryptionAlgorithm
	HashAlgorithmOption       []string
	PFSGroup                  string
	Lifetime                  string
	RekeyTime                 string
	RandTime                  string
	PingHost                  string
	Keepalive                 string
	Description               string
	Disabled                  bool
	Mobile                    bool
}

type IPsecPhase2s []IPsecPhase2

func (ps IPsecPhase2s) GetByUniqID(uniqID string) (*IPsecPhase2, error) {
	for _, p := range ps {
		if p.UniqID == uniqID {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("ipsec phase2 %w with uniqid '%s'", ErrNotFound, uniqID)
}

func (ps IPsecPhase2s) GetIndexByUniqID(uniqID string) (*int, error) {
	for index, p := range ps {
		if p.UniqID == uniqID {
			return &index, nil
		}
	}

	return nil, fmt.Errorf("ipsec phase2 %w with uniqid '%s'", ErrNotFound, uniqID)
}

func (ps IPsecPhase2s) GetByIKEId(ikeId string) IPsecPhase2s {
	var result IPsecPhase2s
	for _, p := range ps {
		if p.IKEId == ikeId {
			result = append(result, p)
		}
	}

	return result
}

// ============================================================================
// Parsing
// ============================================================================

func parseIPsecPhase2IDResponse(resp ipsecPhase2IDResponse) IPsecPhase2ID {
	return IPsecPhase2ID{
		Type:    resp.Type,
		Address: resp.Address,
		NetBits: resp.NetBits,
	}
}

func parseIPsecPhase2Response(resp ipsecPhase2Response) IPsecPhase2 {
	p2 := IPsecPhase2{
		UniqID:      resp.UniqID,
		IKEId:       resp.IKEId,
		Mode:        resp.Mode,
		ReqID:       resp.ReqID,
		LocalID:     parseIPsecPhase2IDResponse(resp.LocalID),
		RemoteID:    parseIPsecPhase2IDResponse(resp.RemoteID),
		Protocol:    resp.Protocol,
		PFSGroup:    resp.PFSGroup,
		Lifetime:    resp.Lifetime,
		RekeyTime:   resp.RekeyTime,
		RandTime:    resp.RandTime,
		PingHost:    resp.PingHost,
		Keepalive:   resp.Keepalive,
		Description: resp.Description,
		Disabled:    resp.Disabled != nil,
		Mobile:      resp.Mobile != nil,
	}

	// NAT local ID (optional)
	if resp.NATLocalID != nil && resp.NATLocalID.Type != "" {
		natLocalID := parseIPsecPhase2IDResponse(*resp.NATLocalID)
		p2.NATLocalID = &natLocalID
	}

	// Encryption algorithms
	for _, alg := range resp.EncryptionAlgorithmOption {
		p2.EncryptionAlgorithmOption = append(p2.EncryptionAlgorithmOption, IPsecPhase2EncryptionAlgorithm{
			Name:   alg.Name,
			KeyLen: alg.KeyLen,
		})
	}

	// Hash algorithms
	p2.HashAlgorithmOption = append(p2.HashAlgorithmOption, resp.HashAlgorithmOption...)

	return p2
}

// generatePHPUniqID generates a unique identifier matching PHP's uniqid()
// format: 8 hex chars of seconds since epoch + 5 hex chars of microseconds.
func generatePHPUniqID() string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixMicro() % 1_000_000

	return fmt.Sprintf("%08x%05x", sec, usec)
}

// ============================================================================
// Form values for POST
// ============================================================================

func ipsecPhase2IDFormValues(prefix string, id IPsecPhase2ID, values url.Values) {
	values.Set(prefix+"id_type", id.Type)
	values.Set(prefix+"id_address", id.Address)
	values.Set(prefix+"id_netbits", id.NetBits)
}

func ipsecPhase2FormValues(p2 IPsecPhase2) url.Values {
	values := url.Values{
		"save":     {"Save"},
		"ikeid":    {p2.IKEId},
		"mode":     {p2.Mode},
		"proto":    {p2.Protocol},
		"pfsgroup": {p2.PFSGroup},
		"lifetime": {p2.Lifetime},
		"descr":    {p2.Description},
	}

	// Unique ID — required by pfSense for both create and update.
	// For creates, generate a PHP-style uniqid (hex timestamp).
	if p2.UniqID != "" {
		values.Set("uniqid", p2.UniqID)
	} else {
		values.Set("uniqid", generatePHPUniqID())
	}

	// ReqID
	if p2.ReqID != "" {
		values.Set("reqid", p2.ReqID)
	}

	// Local/Remote ID
	ipsecPhase2IDFormValues("local", p2.LocalID, values)
	ipsecPhase2IDFormValues("remote", p2.RemoteID, values)

	// NAT local ID (optional)
	if p2.NATLocalID != nil {
		ipsecPhase2IDFormValues("natlocal", *p2.NATLocalID, values)
	}

	// Timers
	if p2.RekeyTime != "" {
		values.Set("rekey_time", p2.RekeyTime)
	}
	if p2.RandTime != "" {
		values.Set("rand_time", p2.RandTime)
	}

	// Keepalive / ping host
	if p2.PingHost != "" {
		values.Set("pinghost", p2.PingHost)
	}
	if p2.Keepalive == "enabled" {
		values.Set("keepalive", "yes")
	}

	// Boolean options (presence-based)
	if p2.Disabled {
		values.Set("disabled", "yes")
	}
	if p2.Mobile {
		values.Set("mobile", "true")
	}

	// Encryption algorithms — Phase 2 uses checkbox array: ealgos[]=name + keylen_name=bits
	for _, alg := range p2.EncryptionAlgorithmOption {
		values.Add("ealgos[]", alg.Name)
		if alg.KeyLen != "" {
			values.Set("keylen_"+alg.Name, alg.KeyLen)
		}
	}

	// Hash algorithms — Phase 2 uses checkbox array: halgos[]=name
	for _, halg := range p2.HashAlgorithmOption {
		values.Add("halgos[]", halg)
	}

	return values
}

// ============================================================================
// Client methods
// ============================================================================

func (pf *Client) getIPsecPhase2s(ctx context.Context) (*IPsecPhase2s, error) {
	command := `
$phase2 = config_get_path('ipsec/phase2', array());
if (!is_array($phase2)) { $phase2 = array(); }
// Handle single entry (not wrapped in array)
if (isset($phase2['uniqid'])) { $phase2 = array($phase2); }
print(json_encode($phase2));
`
	var resp []ipsecPhase2Response
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	phase2s := make(IPsecPhase2s, 0, len(resp))
	for _, r := range resp {
		phase2s = append(phase2s, parseIPsecPhase2Response(r))
	}

	return &phase2s, nil
}

func (pf *Client) GetIPsecPhase2s(ctx context.Context) (*IPsecPhase2s, error) {
	defer pf.read(&pf.mutexes.IPsecPhase2)()

	phase2s, err := pf.getIPsecPhase2s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2 entries, %w", ErrGetOperationFailed, err)
	}

	return phase2s, nil
}

func (pf *Client) GetIPsecPhase2(ctx context.Context, uniqID string) (*IPsecPhase2, error) {
	defer pf.read(&pf.mutexes.IPsecPhase2)()

	phase2s, err := pf.getIPsecPhase2s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2 entries, %w", ErrGetOperationFailed, err)
	}

	p2, err := phase2s.GetByUniqID(uniqID)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2, %w", ErrGetOperationFailed, err)
	}

	return p2, nil
}

func (pf *Client) createOrUpdateIPsecPhase2(ctx context.Context, p2 IPsecPhase2, uniqID *string) error {
	relativeURL := url.URL{Path: "vpn_ipsec_phase2.php"}

	if uniqID != nil {
		q := relativeURL.Query()
		q.Set("uniqid", *uniqID)
		relativeURL.RawQuery = q.Encode()
	}

	// For create, pass ikeid as query parameter too
	if uniqID == nil {
		q := relativeURL.Query()
		q.Set("ikeid", p2.IKEId)
		relativeURL.RawQuery = q.Encode()
	}

	values := ipsecPhase2FormValues(p2)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateIPsecPhase2(ctx context.Context, p2Req IPsecPhase2) (*IPsecPhase2, error) {
	defer pf.write(&pf.mutexes.IPsecPhase2)()

	if err := pf.createOrUpdateIPsecPhase2(ctx, p2Req, nil); err != nil {
		return nil, fmt.Errorf("%w ipsec phase2, %w", ErrCreateOperationFailed, err)
	}

	// Wait for pfSense to finish processing the IPsec config change before
	// attempting to read back via diag_command.php. The create POST triggers
	// an internal config write that can temporarily make PHP command execution
	// unavailable.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(5 * time.Second):
	}

	phase2s, err := pf.getIPsecPhase2s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2 entries after creating, %w", ErrGetOperationFailed, err)
	}

	// Find the newly created entry — last one matching ikeid and description
	var found *IPsecPhase2
	for i := len(*phase2s) - 1; i >= 0; i-- {
		p := (*phase2s)[i]
		if p.IKEId == p2Req.IKEId && p.Description == p2Req.Description {
			found = &p

			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("%w ipsec phase2 after creating, could not find newly created entry", ErrGetOperationFailed)
	}

	return found, nil
}

func (pf *Client) UpdateIPsecPhase2(ctx context.Context, p2Req IPsecPhase2) (*IPsecPhase2, error) {
	defer pf.write(&pf.mutexes.IPsecPhase2)()

	if err := pf.createOrUpdateIPsecPhase2(ctx, p2Req, &p2Req.UniqID); err != nil {
		return nil, fmt.Errorf("%w ipsec phase2, %w", ErrUpdateOperationFailed, err)
	}

	phase2s, err := pf.getIPsecPhase2s(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2 entries after updating, %w", ErrGetOperationFailed, err)
	}

	p2, err := phase2s.GetByUniqID(p2Req.UniqID)
	if err != nil {
		return nil, fmt.Errorf("%w ipsec phase2 after updating, %w", ErrGetOperationFailed, err)
	}

	return p2, nil
}

func (pf *Client) deleteIPsecPhase2(ctx context.Context, index int) error {
	relativeURL := url.URL{Path: "vpn_ipsec.php"}
	values := url.Values{
		"del_p2":    {"Delete selected Phase 2 entries"},
		"p2entry[]": {strconv.Itoa(index)},
	}

	_, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)

	return err
}

func (pf *Client) DeleteIPsecPhase2(ctx context.Context, uniqID string) error {
	defer pf.write(&pf.mutexes.IPsecPhase2)()

	phase2s, err := pf.getIPsecPhase2s(ctx)
	if err != nil {
		return fmt.Errorf("%w ipsec phase2 entries, %w", ErrGetOperationFailed, err)
	}

	index, err := phase2s.GetIndexByUniqID(uniqID)
	if err != nil {
		return fmt.Errorf("%w ipsec phase2, %w", ErrGetOperationFailed, err)
	}

	if err := pf.deleteIPsecPhase2(ctx, *index); err != nil {
		return fmt.Errorf("%w ipsec phase2, %w", ErrDeleteOperationFailed, err)
	}

	phase2s, err = pf.getIPsecPhase2s(ctx)
	if err != nil {
		return fmt.Errorf("%w ipsec phase2 entries after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := phase2s.GetByUniqID(uniqID); err == nil {
		return fmt.Errorf("%w ipsec phase2, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
