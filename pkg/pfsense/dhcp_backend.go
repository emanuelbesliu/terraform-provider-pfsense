package pfsense

import (
	"context"
	"fmt"
)

// DhcpBackend is the pfSense global DHCP server backend selector
// ($config['dhcpbackend']): either 'isc' (legacy ISC DHCP, end-of-life
// upstream and removed in future pfSense versions) or 'kea' (the current
// Kea DHCP server). Both backends read the same $config['dhcpd'] section,
// which is what the DHCPv4 server resources manage; the backend decides
// which engine parses it.
type DhcpBackend = string

const (
	DhcpBackendISC DhcpBackend = "isc"
	DhcpBackendKea DhcpBackend = "kea"
)

func (pf *Client) getDhcpBackend(ctx context.Context) (DhcpBackend, error) {
	command := "print(json_encode(array('backend' => (config_get_path('dhcpbackend') ?: 'isc'))));"

	var resp struct {
		Backend string `json:"backend"`
	}
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return "", err
	}

	return resp.Backend, nil
}

// GetDhcpBackend returns the currently configured DHCP server backend.
func (pf *Client) GetDhcpBackend(ctx context.Context) (DhcpBackend, error) {
	defer pf.read(&pf.mutexes.DhcpBackend)()

	backend, err := pf.getDhcpBackend(ctx)
	if err != nil {
		return "", fmt.Errorf("%w DHCP backend, %w", ErrGetOperationFailed, err)
	}

	return backend, nil
}

func (pf *Client) setDhcpBackend(ctx context.Context, backend DhcpBackend) error {
	if backend != DhcpBackendISC && backend != DhcpBackendKea {
		return fmt.Errorf("%w: invalid DHCP backend %q", ErrClientValidation, backend)
	}

	command := fmt.Sprintf(
		"$b = '%s';"+
			"if (config_get_path('dhcpbackend') !== $b) {"+
			"config_set_path('dhcpbackend', $b);"+
			"write_config('Terraform: DHCP server backend changed');"+
			"}"+
			"print(json_encode(array('backend' => (config_get_path('dhcpbackend') ?: 'isc'))));",
		phpLiteral(backend),
	)

	var resp struct {
		Backend string `json:"backend"`
	}
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return err
	}

	if resp.Backend != backend {
		return fmt.Errorf("%w: DHCP backend read-back %q, expected %q", ErrServerValidation, resp.Backend, backend)
	}

	return nil
}

// CreateDhcpBackend sets the global DHCP server backend.
func (pf *Client) CreateDhcpBackend(ctx context.Context, backend DhcpBackend) (DhcpBackend, error) {
	defer pf.write(&pf.mutexes.DhcpBackend)()

	if err := pf.setDhcpBackend(ctx, backend); err != nil {
		return "", fmt.Errorf("%w DHCP backend, %w", ErrCreateOperationFailed, err)
	}

	return pf.getDhcpBackend(ctx)
}

// UpdateDhcpBackend sets the global DHCP server backend.
func (pf *Client) UpdateDhcpBackend(ctx context.Context, backend DhcpBackend) (DhcpBackend, error) {
	defer pf.write(&pf.mutexes.DhcpBackend)()

	if err := pf.setDhcpBackend(ctx, backend); err != nil {
		return "", fmt.Errorf("%w DHCP backend, %w", ErrUpdateOperationFailed, err)
	}

	return pf.getDhcpBackend(ctx)
}

// DeleteDhcpBackend resets the global DHCP backend to its default ('isc')
// only when it still matches the previously managed value.
func (pf *Client) DeleteDhcpBackend(ctx context.Context, backend DhcpBackend) error {
	defer pf.write(&pf.mutexes.DhcpBackend)()

	command := fmt.Sprintf(
		"$e = '%s';"+
			"if ((config_get_path('dhcpbackend') ?: 'isc') === $e) {"+
			"config_set_path('dhcpbackend', 'isc');"+
			"write_config('Terraform: DHCP server backend reset');"+
			"}"+
			"print(json_encode(array('ok' => true)));",
		phpLiteral(backend),
	)

	if err := pf.executePHPCommand(ctx, command, nil); err != nil {
		return fmt.Errorf("%w DHCP backend, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}
