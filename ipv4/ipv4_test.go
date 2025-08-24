package ipv4_test

import (
	"net"
	"strings"
	"testing"

	"github.com/ronny/ripcalc/ipv4"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		name       string
		cidr       string
		wantAddr   string
		wantPrefix int
		wantError  bool
	}{
		{
			name:       "valid /24 network",
			cidr:       "192.168.0.0/24",
			wantAddr:   "192.168.0.0",
			wantPrefix: 24,
			wantError:  false,
		},
		{
			name:       "valid /16 network",
			cidr:       "10.0.0.0/16",
			wantAddr:   "10.0.0.0",
			wantPrefix: 16,
			wantError:  false,
		},
		{
			name:      "invalid CIDR",
			cidr:      "192.168.0.0/33",
			wantError: true,
		},
		{
			name:      "invalid IP",
			cidr:      "256.256.256.256/24",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipv4.ParseCIDR(tt.cidr)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseCIDR() expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("ParseCIDR() error = %v", err)
				return
			}

			if got.Address.String() != tt.wantAddr {
				t.Errorf("ParseCIDR() address = %v, want %v", got.Address.String(), tt.wantAddr)
			}

			if got.PrefixLength != tt.wantPrefix {
				t.Errorf("ParseCIDR() prefix = %v, want %v", got.PrefixLength, tt.wantPrefix)
			}
		})
	}
}

func TestNetwork_Calculate(t *testing.T) {
	tests := []struct {
		name          string
		cidr          string
		wantNetmask   string
		wantWildcard  string
		wantNetwork   string
		wantBroadcast string
		wantHostMin   string
		wantHostMax   string
		wantHostCount uint32
		wantClass     string
		wantType      string
	}{
		{
			name:          "192.168.0.0/24",
			cidr:          "192.168.0.0/24",
			wantNetmask:   "255.255.255.0",
			wantWildcard:  "0.0.0.255",
			wantNetwork:   "192.168.0.0",
			wantBroadcast: "192.168.0.255",
			wantHostMin:   "192.168.0.1",
			wantHostMax:   "192.168.0.254",
			wantHostCount: 254,
			wantClass:     "C",
			wantType:      "Private Internet",
		},
		{
			name:          "8.8.8.0/24",
			cidr:          "8.8.8.0/24",
			wantNetmask:   "255.255.255.0",
			wantWildcard:  "0.0.0.255",
			wantNetwork:   "8.8.8.0",
			wantBroadcast: "8.8.8.255",
			wantHostMin:   "8.8.8.1",
			wantHostMax:   "8.8.8.254",
			wantHostCount: 254,
			wantClass:     "A",
			wantType:      "Public Internet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv4.ParseCIDR(tt.cidr)
			if err != nil {
				t.Fatalf("ParseCIDR() error = %v", err)
			}

			err = network.Calculate()
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if net.IP(network.Netmask).String() != tt.wantNetmask {
				t.Errorf("Netmask = %v, want %v", net.IP(network.Netmask).String(), tt.wantNetmask)
			}

			if network.Wildcard.String() != tt.wantWildcard {
				t.Errorf("Wildcard = %v, want %v", network.Wildcard.String(), tt.wantWildcard)
			}

			if network.Network.String() != tt.wantNetwork {
				t.Errorf("Network = %v, want %v", network.Network.String(), tt.wantNetwork)
			}

			if network.Broadcast.String() != tt.wantBroadcast {
				t.Errorf("Broadcast = %v, want %v", network.Broadcast.String(), tt.wantBroadcast)
			}

			if network.HostMin.String() != tt.wantHostMin {
				t.Errorf("HostMin = %v, want %v", network.HostMin.String(), tt.wantHostMin)
			}

			if network.HostMax.String() != tt.wantHostMax {
				t.Errorf("HostMax = %v, want %v", network.HostMax.String(), tt.wantHostMax)
			}

			if network.HostCount != tt.wantHostCount {
				t.Errorf("HostCount = %v, want %v", network.HostCount, tt.wantHostCount)
			}

			if network.Class != tt.wantClass {
				t.Errorf("Class = %v, want %v", network.Class, tt.wantClass)
			}

			if network.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", network.Type, tt.wantType)
			}
		})
	}
}

func TestFormatBinary(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{
			name: "192.168.0.0",
			ip:   "192.168.0.0",
			want: "11000000.10101000.00000000.00000000",
		},
		{
			name: "255.255.255.0",
			ip:   "255.255.255.0",
			want: "11111111.11111111.11111111.00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip).To4()

			got := ipv4.FormatBinary(ip)
			if got != tt.want {
				t.Errorf("FormatBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_FormattedText(t *testing.T) {
	network, err := ipv4.ParseCIDR("192.168.0.1/24")
	if err != nil {
		t.Fatalf("ParseCIDR() error = %v", err)
	}

	err = network.Calculate()
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	output := network.FormattedText()

	// Check that output contains expected elements
	expectedParts := []string{
		"Address:",
		"192.168.0.1",
		"11000000.10101000.00000000. 00000001",
		"Netmask:",
		"255.255.255.0",
		"= 24",
		"11111111.11111111.11111111. 00000000",
		"Wildcard:",
		"0.0.0.255",
		"00000000.00000000.00000000. 11111111",
		"Network:",
		"192.168.0.0/24",
		"11000000.10101000.00000000. 00000000",
		"HostMin:",
		"192.168.0.1",
		"HostMax:",
		"192.168.0.254",
		"11000000.10101000.00000000. 11111110",
		"Broadcast:",
		"192.168.0.255",
		"11000000.10101000.00000000. 11111111",
		"Hosts:",
		"254",
		"Class C",
		"Private Internet",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("FormattedText() missing expected part: %q", part)
		}
	}
}

func TestNetwork_String(t *testing.T) {
	tests := []struct {
		name string
		cidr string
		want string
	}{
		{
			name: "/24 network",
			cidr: "192.168.0.0/24",
			want: "192.168.0.0/24",
		},
		{
			name: "/16 network",
			cidr: "10.0.0.0/16",
			want: "10.0.0.0/16",
		},
		{
			name: "/32 single host",
			cidr: "192.168.1.1/32",
			want: "192.168.1.1/32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv4.ParseCIDR(tt.cidr)
			if err != nil {
				t.Fatalf("ParseCIDR() error = %v", err)
			}

			got := network.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkLocalClassification(t *testing.T) {
	network, err := ipv4.ParseCIDR("169.254.1.1/16")
	if err != nil {
		t.Fatalf("ParseCIDR() error = %v", err)
	}

	err = network.Calculate()
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	if network.Type != "Link Local" {
		t.Errorf("Type = %v, want Link Local", network.Type)
	}
}

func TestClassifyAddress(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		wantClass string
	}{
		// Class A: 0.0.0.0 to 127.255.255.255
		{
			name:      "Class A - 0.0.0.1",
			ip:        "0.0.0.1",
			wantClass: "A",
		},
		{
			name:      "Class A - 10.0.0.1",
			ip:        "10.0.0.1",
			wantClass: "A",
		},
		{
			name:      "Class A - 127.255.255.255",
			ip:        "127.255.255.255",
			wantClass: "A",
		},
		// Class B: 128.0.0.0 to 191.255.255.255
		{
			name:      "Class B - 128.0.0.1",
			ip:        "128.0.0.1",
			wantClass: "B",
		},
		{
			name:      "Class B - 172.16.0.1",
			ip:        "172.16.0.1",
			wantClass: "B",
		},
		{
			name:      "Class B - 191.255.255.255",
			ip:        "191.255.255.255",
			wantClass: "B",
		},
		// Class C: 192.0.0.0 to 223.255.255.255
		{
			name:      "Class C - 192.0.0.1",
			ip:        "192.0.0.1",
			wantClass: "C",
		},
		{
			name:      "Class C - 192.168.1.1",
			ip:        "192.168.1.1",
			wantClass: "C",
		},
		{
			name:      "Class C - 223.255.255.255",
			ip:        "223.255.255.255",
			wantClass: "C",
		},
		// Class D: 224.0.0.0 to 239.255.255.255 (Multicast)
		{
			name:      "Class D - 224.0.0.1",
			ip:        "224.0.0.1",
			wantClass: "D",
		},
		{
			name:      "Class D - 230.0.0.1",
			ip:        "230.0.0.1",
			wantClass: "D",
		},
		{
			name:      "Class D - 239.255.255.255",
			ip:        "239.255.255.255",
			wantClass: "D",
		},
		// Class E: 240.0.0.0 to 255.255.255.255 (Reserved)
		{
			name:      "Class E - 240.0.0.1",
			ip:        "240.0.0.1",
			wantClass: "E",
		},
		{
			name:      "Class E - 250.0.0.1",
			ip:        "250.0.0.1",
			wantClass: "E",
		},
		{
			name:      "Class E - 255.255.255.255",
			ip:        "255.255.255.255",
			wantClass: "E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := ipv4.ParseCIDR(tt.ip + "/24")
			if err != nil {
				t.Fatalf("ParseCIDR() error = %v", err)
			}

			err = network.Calculate()
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if network.Class != tt.wantClass {
				t.Errorf("Class = %v, want %v", network.Class, tt.wantClass)
			}
		})
	}
}

func TestMulticastAddressClassification(t *testing.T) {
	network, err := ipv4.ParseCIDR("224.0.0.1/24")
	if err != nil {
		t.Fatalf("ParseCIDR() error = %v", err)
	}

	err = network.Calculate()
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// Should be Class D
	if network.Class != "D" {
		t.Errorf("Class = %v, want D", network.Class)
	}

	// Should be Multicast type
	if network.Type != "Multicast" {
		t.Errorf("Type = %v, want Multicast", network.Type)
	}
}

func TestReservedAddressClassification(t *testing.T) {
	network, err := ipv4.ParseCIDR("240.0.0.1/24")
	if err != nil {
		t.Fatalf("ParseCIDR() error = %v", err)
	}

	err = network.Calculate()
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	// Should be Class E
	if network.Class != "E" {
		t.Errorf("Class = %v, want E", network.Class)
	}

	// Should be Public Internet (since no special range defined for Class E)
	if network.Type != "Public Internet" {
		t.Errorf("Type = %v, want Public Internet", network.Type)
	}
}
