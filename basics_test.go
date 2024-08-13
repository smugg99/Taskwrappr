// basics_test.go
package taskwrappr

import (
	"testing"
)

func TestConditionals(t *testing.T) {
	s, err := NewScript("scripts/basics.tw", GetBuiltIn())
	if err != nil {
		t.Errorf("NewScript returned an error: %s", err)
	}

	//var buf bytes.Buffer

    // mw := io.MultiWriter(os.Stdout, &buf)

	if err := s.Run(); err != nil {
		t.Errorf("run returned an error: %s", err)
	}

	//fmt.Printf("Captured logs: %s", buf.String())
}