package initchecker

import (
	"testing"
)

func TestCheck(t *testing.T) {
	checker := Checker{}
	s := []string{"localhost:8080", "localhost:9999"}
	checker.Addrs = s
	checker.Check()
}
