package pkg_test

import (
	"reflect"
	"testing"

	"example.com/gocab/pkg"
)

func TestAddCommand(t *testing.T) {
	expected := &pkg.Word{"german", "Hallo", "Hello", "Hallo, wie gehts", []string{"greetings"}}

	cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()))

	received, _ := cmd.Execute([]string{"-lang", "german", "-word", "Hallo", "-meaning", "Hello", "-example", "Hallo, wie gehts", "-tags", "greetings"})

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("expected %v, got %v", expected, received)
	}
}
