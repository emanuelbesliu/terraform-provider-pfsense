package pfsense

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (pf *Client) ApplyIPsecChanges(ctx context.Context) error {
	pf.mutexes.IPsecApply.Lock()
	defer pf.mutexes.IPsecApply.Unlock()

	relativeURL := url.URL{Path: "vpn_ipsec.php"}
	values := url.Values{
		"apply": {"Apply Changes"},
	}

	resp, err := pf.call(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return fmt.Errorf("%w ipsec changes, %w", ErrApplyOperationFailed, err)
	}

	defer resp.Body.Close() //nolint:errcheck
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}
