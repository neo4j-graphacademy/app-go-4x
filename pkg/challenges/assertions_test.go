package challenges_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		t.Fatal(fmt.Errorf("expected nil value"))
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

func assertResultHasNextRecord(t *testing.T, result neo4j.Result) {
	if !result.Next() {
		t.Fatalf("Expected `.Next()` to return true on neo4j.Result.  No next record found.")
	}
}

func assertContains(t *testing.T, str string, contains string) {
	if strings.Contains(str, contains) {
		t.Fatalf("Expected '%s' to contain '%s'", str, contains)
	}
}
