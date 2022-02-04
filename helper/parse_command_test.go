package helper

import (
	"fmt"
	"testing"
)

func TestParseCommands(T *testing.T) {
	result := ParseCommands()
	fmt.Println(result)
}
