package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunWithValidCIDR(t *testing.T) {
	err := runWithArgs([]string{"ripcalc", "192.168.0.0/24"})
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}
}

func TestRunWithInvalidCIDR(t *testing.T) {
	err := runWithArgs([]string{"ripcalc", "invalid-cidr"})
	if err == nil {
		t.Error("Expected run() to fail with invalid CIDR, but it succeeded")
	}

	expectedError := "invalid"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Error message should contain %q, got: %v", expectedError, err)
	}
}

func TestRunWithNoArguments(t *testing.T) {
	err := runWithArgs([]string{"ripcalc"})
	if err == nil {
		t.Error("Expected run() to fail with no arguments, but it succeeded")
	}

	expectedError := "no CIDR argument provided"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Error message should contain %q, got: %v", expectedError, err)
	}
}

func TestRunWithHelpFlags(t *testing.T) {
	helpFlags := []string{"-h", "--help"}

	for _, flag := range helpFlags {
		t.Run("help_flag_"+flag, func(t *testing.T) {
			err := runWithArgs([]string{"ripcalc", flag})
			if err != nil {
				t.Errorf("Help command should not return error for flag %s: %v", flag, err)
			}
		})
	}
}

func TestRunWithDifferentNetworks(t *testing.T) {
	tests := []struct {
		name string
		cidr string
	}{
		// IPv4 tests
		{
			name: "IPv4 Class A Private",
			cidr: "10.0.0.1/8",
		},
		{
			name: "IPv4 Class B Private",
			cidr: "172.16.0.1/16",
		},
		{
			name: "IPv4 Class D Multicast",
			cidr: "224.0.0.1/24",
		},
		{
			name: "IPv4 Class A Public",
			cidr: "8.8.8.8/24",
		},
		// IPv6 tests
		{
			name: "IPv6 Global Unicast",
			cidr: "2001:db8::/64",
		},
		{
			name: "IPv6 Link-Local",
			cidr: "fe80::1/64",
		},
		{
			name: "IPv6 Loopback",
			cidr: "::1/128",
		},
		{
			name: "IPv6 Multicast",
			cidr: "ff02::1/128",
		},
		{
			name: "IPv6 Unique Local",
			cidr: "fd00::1/64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runWithArgs([]string{"ripcalc", tt.cidr})
			if err != nil {
				t.Fatalf("run() failed for %s: %v", tt.cidr, err)
			}
		})
	}
}

func TestIsIPv6CIDR(t *testing.T) {
	tests := []struct {
		name     string
		cidr     string
		expected bool
	}{
		{
			name:     "IPv4 CIDR",
			cidr:     "192.168.1.0/24",
			expected: false,
		},
		{
			name:     "IPv6 CIDR",
			cidr:     "2001:db8::/64",
			expected: true,
		},
		{
			name:     "IPv6 loopback",
			cidr:     "::1/128",
			expected: true,
		},
		{
			name:     "IPv6 link-local",
			cidr:     "fe80::/10",
			expected: true,
		},
		{
			name:     "invalid CIDR",
			cidr:     "invalid",
			expected: false,
		},
		{
			name:     "IPv4-mapped IPv6 should be treated as IPv6",
			cidr:     "::ffff:192.168.1.1/128",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPv6CIDR(tt.cidr)
			if result != tt.expected {
				t.Errorf("isIPv6CIDR(%q) = %v, expected %v", tt.cidr, result, tt.expected)
			}
		})
	}
}

// Integration tests that verify actual CLI output
func TestIntegration_IPv4_Output(t *testing.T) {
	tests := []struct {
		name              string
		cidr              string
		expectedElements  []string
		expectedBinaryElements []string
	}{
		{
			name: "Class C Private",
			cidr: "192.168.1.0/24",
			expectedElements: []string{
				"Address:",
				"192.168.1.0",
				"Prefix:",
				"/24",
				"Netmask:",
				"255.255.255.0",
				"Wildcard:",
				"0.0.0.255",
				"Network:",
				"192.168.1.0/24",
				"First host:",
				"192.168.1.1",
				"Last host:",
				"192.168.1.254",
				"Broadcast:",
				"192.168.1.255",
				"Host count:",
				"254",
				"Class C",
				"Private Internet",
			},
			expectedBinaryElements: []string{
				"11000000.10101000.00000001.",
				"11111111.11111111.11111111.",
				"00000000.00000000.00000000.",
			},
		},
		{
			name: "Class A Public",
			cidr: "8.8.8.8/24",
			expectedElements: []string{
				"Address:",
				"8.8.8.8",
				"Class A",
				"Public Internet",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				err := runWithArgs([]string{"ripcalc", tt.cidr})
				if err != nil {
					t.Fatalf("run() failed: %v", err)
				}
			})

			// Verify all expected text elements are present
			for _, element := range tt.expectedElements {
				if !strings.Contains(output, element) {
					t.Errorf("Output missing expected element: %q\nFull output:\n%s", element, output)
				}
			}

			// Verify binary elements are present
			for _, binaryElement := range tt.expectedBinaryElements {
				if !strings.Contains(output, binaryElement) {
					t.Errorf("Output missing expected binary element: %q\nFull output:\n%s", binaryElement, output)
				}
			}

			// Verify network/host boundary space exists in binary output
			if strings.Contains(output, " 0000") || strings.Contains(output, " 1111") {
				// Good - found network/host boundary space
			} else {
				t.Error("Output should contain network/host boundary space in binary representation")
			}
		})
	}
}

func TestIntegration_IPv6_Output(t *testing.T) {
	tests := []struct {
		name             string
		cidr             string
		expectedElements []string
		expectBoundary   bool
	}{
		{
			name: "IPv6 Documentation Range",
			cidr: "2001:db8::/64",
			expectedElements: []string{
				"Address:",
				"2001:db8::",
				"Prefix:",
				"/64",
				"Network:",
				"2001:db8::/64",
				"First host:",
				"Last host:",
				"2001:db8::ffff:ffff:ffff:ffff",
				"Host count:",
				"2^64",
				"Documentation",
				"RFC Example",
			},
			expectBoundary: true,
		},
		{
			name: "IPv6 Loopback",
			cidr: "::1/128",
			expectedElements: []string{
				"Address:",
				"::1",
				"/128",
				"Loopback",
				"Host-only",
				"Host count:",
				"1",
			},
			expectBoundary: false, // /128 has no host bits
		},
		{
			name: "IPv6 Link-Local",
			cidr: "fe80::/10",
			expectedElements: []string{
				"Address:",
				"fe80::",
				"/10",
				"Link-Local Unicast",
				"Auto-configured",
				"2^118 (astronomical)",
			},
			expectBoundary: true,
		},
		{
			name: "IPv6 Multicast",
			cidr: "ff02::1/128",
			expectedElements: []string{
				"ff02::1",
				"Multicast Link-Local",
				"Group Communication",
			},
			expectBoundary: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				err := runWithArgs([]string{"ripcalc", tt.cidr})
				if err != nil {
					t.Fatalf("run() failed: %v", err)
				}
			})

			// Verify all expected elements are present
			for _, element := range tt.expectedElements {
				if !strings.Contains(output, element) {
					t.Errorf("Output missing expected element: %q\nFull output:\n%s", element, output)
				}
			}

			// Since IPv6 binary is now optional, we don't check for binary by default
			// These tests are for the default output format
		})
	}
}

// Helper function to capture stdout during test execution
func captureStdout(t *testing.T, fn func()) string {
	// Save original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Replace stdout with our pipe
	os.Stdout = w

	// Channel to receive captured output
	outputCh := make(chan string, 1)

	// Start reading from the pipe in a goroutine
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	// Execute the function
	fn()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the captured output
	output := <-outputCh
	r.Close()

	return output
}

func TestIPv6Flags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		shouldHaveMask bool
		shouldHaveBinary bool
	}{
		{
			name: "default - no flags",
			args: []string{"ripcalc", "2001:db8::/64"},
			shouldHaveMask: false,
			shouldHaveBinary: false,
		},
		{
			name: "ipv6-mask only",
			args: []string{"ripcalc", "--ipv6-mask", "2001:db8::/64"},
			shouldHaveMask: true,
			shouldHaveBinary: false,
		},
		{
			name: "ipv6-binary only",
			args: []string{"ripcalc", "--ipv6-binary", "2001:db8::/64"},
			shouldHaveMask: false,
			shouldHaveBinary: true,
		},
		{
			name: "both flags",
			args: []string{"ripcalc", "--ipv6-mask", "--ipv6-binary", "2001:db8::/64"},
			shouldHaveMask: true,
			shouldHaveBinary: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				err := runWithArgs(tt.args)
				if err != nil {
					t.Fatalf("run() failed: %v", err)
				}
			})

			// Check mask presence
			hasMask := strings.Contains(output, "Netmask:") && strings.Contains(output, "Wildcard:")
			if hasMask != tt.shouldHaveMask {
				t.Errorf("Mask presence = %v, expected %v", hasMask, tt.shouldHaveMask)
			}

			// Check binary presence
			hasBinary := strings.Contains(output, "0010000000000001:0000110110111000")
			if hasBinary != tt.shouldHaveBinary {
				t.Errorf("Binary presence = %v, expected %v", hasBinary, tt.shouldHaveBinary)
			}

			// If mask is shown, verify correct values
			if tt.shouldHaveMask {
				if !strings.Contains(output, "ffff:ffff:ffff:ffff::") {
					t.Error("Output should contain correct netmask")
				}
				if !strings.Contains(output, "::ffff:ffff:ffff:ffff") {
					t.Error("Output should contain correct wildcard")
				}
			}
		})
	}
}
