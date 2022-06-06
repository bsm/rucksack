package rucksack_test

import (
	"reflect"
	"testing"

	. "github.com/bsm/rucksack/v4"
)

func expectTags(t *testing.T, s string, exp []string) {
	t.Helper()

	if got := Tags(s); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestTags(t *testing.T) {
	t.Run("blank", func(t *testing.T) {
		expectTags(t, "", []string{})
	})
	t.Run("simple", func(t *testing.T) {
		expectTags(t, "a,b", []string{"a", "b"})
	})
	t.Run("spaced", func(t *testing.T) {
		expectTags(t, "a ,  b", []string{"a", "b"})
	})
}

func expectFields(t *testing.T, s string, exp map[string]interface{}) {
	t.Helper()

	if got := Fields(s); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestFields(t *testing.T) {
	t.Run("blank", func(t *testing.T) {
		expectFields(t, "", nil)
	})
	t.Run("simple", func(t *testing.T) {
		expectFields(t, "k1:v1,k2:v2", map[string]interface{}{"k1": "v1", "k2": "v2"})
	})
	t.Run("with spaces", func(t *testing.T) {
		expectFields(t, " k1:v1 ,   k2:v2", map[string]interface{}{"k1": "v1", "k2": "v2"})
	})
}
