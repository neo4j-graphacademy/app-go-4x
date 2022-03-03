package challenges_test

import (
	"fmt"
	"testing"
)

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Fatal(fmt.Errorf("expected non-nil value"))
	}
}

func assertStringNotEmpty(t *testing.T, str string) {
	t.Helper()
	if str == "" {
		t.Fatal(fmt.Errorf("expected string to be non-empty"))
	}
}
