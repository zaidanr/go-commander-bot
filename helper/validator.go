package helper

import (
	"regexp"
)

// func validateIP(s string) bool {
// 	octet := strings.Split(s, ".")
// 	if r.TypeOf(octet[0] == r.Int.String())
// }

func Validate(s string) bool {

	// Hostname validation
	r1, _ := regexp.Compile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
		]{2,3})$`)

	// IP Validation
	r2, _ := regexp.Compile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`)
	// r2 :=
	return r1.MatchString(s) || r2.MatchString(s)
}
