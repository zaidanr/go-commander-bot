package helper

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
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
