package ipv6_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ronny/ripcalc/ipv6"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantErr bool
	}{
		{
			name:    "valid global unicast",
			cidr:    "2001:db8:85a3::8a2e:370:7334/64",
			wantErr: false,
		},
		{
			name:    "valid link-local",
			cidr:    "fe80::1/64",
			wantErr: false,
		},
		{
			name:    "valid loopback",
			cidr:    "::1/128",
			wantErr: false,
		},
		{
			name:    "invalid CIDR",
			cidr:    "not-an-ip",
			wantErr: true,
		},
		{
			name:    "IPv4 address should fail",
			cidr:    "192.168.1.1/24",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv6.ParseCIDR(tt.cidr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCIDR() expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("ParseCIDR() unexpected error: %v", err)
				return
			}

			if network == nil {
				t.Errorf("ParseCIDR() returned nil network")
				return
			}

			// Basic validation that we can calculate
			err = network.Calculate()
			if err != nil {
				t.Errorf("Calculate() unexpected error: %v", err)
			}
		})
	}
}

func TestAddressClassification(t *testing.T) {
	tests := []struct {
		name          string
		cidr          string
		expectedClass string
		expectedType  string
	}{
		{
			name:          "loopback",
			cidr:          "::1/128",
			expectedClass: "Loopback",
			expectedType:  "Host-only",
		},
		{
			name:          "link-local",
			cidr:          "fe80::1/64",
			expectedClass: "Link-Local Unicast",
			expectedType:  "Auto-configured",
		},
		{
			name:          "documentation",
			cidr:          "2001:db8::1/32",
			expectedClass: "Documentation",
			expectedType:  "RFC Example",
		},
		{
			name:          "global unicast",
			cidr:          "2001:470::1/64",
			expectedClass: "Global Unicast",
			expectedType:  "Internet Routable",
		},
		{
			name:          "unique local",
			cidr:          "fd00::1/64",
			expectedClass: "Unique Local Address",
			expectedType:  "Private",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv6.ParseCIDR(tt.cidr)
			if err != nil {
				t.Fatalf("ParseCIDR() unexpected error: %v", err)
			}

			err = network.Calculate()
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}

			if network.Class != tt.expectedClass {
				t.Errorf("Expected class %q, got %q", tt.expectedClass, network.Class)
			}

			if network.Type != tt.expectedType {
				t.Errorf("Expected type %q, got %q", tt.expectedType, network.Type)
			}
		})
	}
}

func TestFormattedText(t *testing.T) {
	network, err := ipv6.ParseCIDR("2001:db8::1/64")
	if err != nil {
		t.Fatalf("ParseCIDR() unexpected error: %v", err)
	}

	err = network.Calculate()
	if err != nil {
		t.Fatalf("Calculate() unexpected error: %v", err)
	}

	output := network.FormattedText()
	if output == "" {
		t.Error("FormattedText() returned empty string")
	}

	// Basic checks that expected elements are present
	expectedElements := []string{
		"Address:",
		"Network:",
		"First host:",
		"Last host:",
		"Host count:",
		"Documentation",
		"RFC Example",
	}

	for _, element := range expectedElements {
		if !containsString(output, element) {
			t.Errorf("FormattedText() missing expected element: %q", element)
		}
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

func TestNetworkCalculations(t *testing.T) {
	tests := []struct {
		name           string
		cidr           string
		expectedPrefix int
	}{
		{
			name:           "standard /64",
			cidr:           "2001:db8::/64",
			expectedPrefix: 64,
		},
		{
			name:           "single host /128",
			cidr:           "::1/128",
			expectedPrefix: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv6.ParseCIDR(tt.cidr)
			if err != nil {
				t.Fatalf("ParseCIDR() unexpected error: %v", err)
			}

			if network.PrefixLength != tt.expectedPrefix {
				t.Errorf("Expected prefix length %d, got %d", tt.expectedPrefix, network.PrefixLength)
			}

			err = network.Calculate()
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}

			// Verify we have non-nil results
			if network.Network == nil {
				t.Error("Network address is nil after Calculate()")
			}

			if network.HostMin == nil {
				t.Error("HostMin is nil after Calculate()")
			}

			if network.HostMax == nil {
				t.Error("HostMax is nil after Calculate()")
			}
		})
	}
}

func TestFormatBinary(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "loopback",
			address:  "::1",
			expected: "0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000001",
		},
		{
			name:     "all zeros",
			address:  "::",
			expected: "0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000:0000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.address)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.address)
			}

			result := ipv6.FormatBinary(ip)
			if result != tt.expected {
				t.Errorf("FormatBinary() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestFormatBinaryWithMask(t *testing.T) {
	tests := []struct {
		name            string
		address         string
		prefixLength    int
		shouldHaveSpace bool
	}{
		{
			name:            "standard /64",
			address:         "2001:db8::1",
			prefixLength:    64,
			shouldHaveSpace: true,
		},
		{
			name:            "single host /128",
			address:         "::1",
			prefixLength:    128,
			shouldHaveSpace: false,
		},
		{
			name:            "/48 prefix",
			address:         "2001:db8::1",
			prefixLength:    48,
			shouldHaveSpace: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.address)
			if ip == nil {
				t.Fatalf("Failed to parse IP: %s", tt.address)
			}

			result := ipv6.FormatBinaryWithMask(ip, tt.prefixLength)

			if result == "" {
				t.Error("FormatBinaryWithMask() returned empty string")
			}

			hasSpace := containsString(result, " ")
			if hasSpace != tt.shouldHaveSpace {
				t.Errorf("FormatBinaryWithMask() space presence = %v, expected %v", hasSpace, tt.shouldHaveSpace)
			}

			// Verify the result contains binary digits and colons
			if !containsString(result, ":") && tt.prefixLength < 128 {
				t.Error("FormatBinaryWithMask() should contain colon separators")
			}
		})
	}
}

func TestMulticastClassification(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		expectedClass string
	}{
		{
			name:          "link-local multicast",
			address:       "ff02::1",
			expectedClass: "Multicast Link-Local",
		},
		{
			name:          "global multicast",
			address:       "ff0e::1",
			expectedClass: "Multicast Global",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv6.ParseCIDR(fmt.Sprintf("%s/128", tt.address))
			if err != nil {
				t.Fatalf("ParseCIDR() unexpected error: %v", err)
			}

			err = network.Calculate()
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}

			if network.Class != tt.expectedClass {
				t.Errorf("Expected class %q, got %q", tt.expectedClass, network.Class)
			}
		})
	}
}
