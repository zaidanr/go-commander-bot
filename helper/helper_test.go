package helper

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommands(T *testing.T) {
	result := ParseCommands()
	fmt.Println(result)
}

func TestExec(t *testing.T) {
	cmd := exec.Command("whoami")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func TestValidateTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "ValidHostname",
			input:    "google.com",
			expected: true,
		},
		{
			name:     "InvalidHostname (RCE)",
			input:    "google.com;id",
			expected: false,
		},
		{
			name:     "SingleWord(g.co)",
			input:    "g.co",
			expected: true,
		},
		{
			name:     "ValidAddress (Private)",
			input:    "192.168.9.1",
			expected: true,
		},
		{
			name:     "ValidAddress (Public)",
			input:    "250.253.123.123",
			expected: true,
		},
		{
			name:     "InvalidAddress1",
			input:    "192.168.0.300",
			expected: false,
		},
		{
			name:     "InvalidAddress (RCE)",
			input:    "192.168.0.1;id",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Validate(test.input)
			assert.Equal(t, test.expected, result)
		})
	}

	// fmt.Println("H")
}

func TestValidate(t *testing.T) {
	fmt.Printf("%t\n", Validate("google.com"))
}
