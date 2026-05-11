//go:build debug

package pfsense

import (
	"context"
	"fmt"
	"net/url"
	"testing"
)

func TestDebugNATOneToOne(t *testing.T) {
	ctx := context.Background()
	u, _ := url.Parse("http://10.0.161.1")
	skip := true
	client, err := NewClient(ctx, &Options{
		URL:           u,
		Username:      "emanuelb",
		Password:      "Dd4c0be564cc98",
		TLSSkipVerify: &skip,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check 1:1 NAT config structure
	var raw string
	cmd := "$data = config_get_path('nat/onetoone'); print(json_encode($data));"
	err = client.executePHPCommand(ctx, cmd, &raw)
	fmt.Printf("nat/onetoone: %s\nerr: %v\n", raw, err)

	// Try alternate path
	var raw2 string
	cmd2 := "$data = config_get_path('nat/onetoone/rule'); print(json_encode($data));"
	err = client.executePHPCommand(ctx, cmd2, &raw2)
	fmt.Printf("nat/onetoone/rule: %s\nerr: %v\n", raw2, err)

	// Check NAT keys
	var raw3 string
	cmd3 := "$data = config_get_path('nat'); if(is_array($data)) { print(json_encode(array_keys($data))); } else { print(json_encode($data)); }"
	err = client.executePHPCommand(ctx, cmd3, &raw3)
	fmt.Printf("nat keys: %s\nerr: %v\n", raw3, err)
}
