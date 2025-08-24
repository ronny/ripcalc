package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunWithValidCIDR(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	
	// Set up args for testing
	os.Args = []string{"ripcalc", "192.168.0.0/24"}

	err := run()
	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}
}

func TestRunWithInvalidCIDR(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{"ripcalc", "invalid-cidr"}

	err := run()
	if err == nil {
		t.Error("Expected run() to fail with invalid CIDR, but it succeeded")
	}

	expectedError := "invalid CIDR notation"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Error message should contain %q, got: %v", expectedError, err)
	}
}

func TestRunWithNoArguments(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	
	// Set up args with no CIDR argument
	os.Args = []string{"ripcalc"}

	err := run()
	if err == nil {
		t.Error("Expected run() to fail with no arguments, but it succeeded")
	}

	expectedError := "no CIDR argument provided"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Error message should contain %q, got: %v", expectedError, err)
	}
}

func TestRunWithHelpFlags(t *testing.T) {
	helpFlags := []string{"-h", "--help", "help"}
	
	for _, flag := range helpFlags {
		t.Run("help_flag_"+flag, func(t *testing.T) {
			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()
			
			os.Args = []string{"ripcalc", flag}

			err := run()
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
		{
			name: "Class A Private",
			cidr: "10.0.0.1/8",
		},
		{
			name: "Class B Private", 
			cidr: "172.16.0.1/16",
		},
		{
			name: "Class D Multicast",
			cidr: "224.0.0.1/24",
		},
		{
			name: "Class A Public",
			cidr: "8.8.8.8/24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()
			
			os.Args = []string{"ripcalc", tt.cidr}

			err := run()
			if err != nil {
				t.Fatalf("run() failed for %s: %v", tt.cidr, err)
			}
		})
	}
}