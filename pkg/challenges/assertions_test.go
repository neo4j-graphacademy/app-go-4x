package challenges_test

import (
	"fmt"
	"reflect"
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

func assertEquals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}

func assertNotEquals(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Fatalf("Received %v (type %v), expected NOT %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
	}
}
