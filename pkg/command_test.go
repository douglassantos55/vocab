package pkg_test

import (
	"reflect"
	"testing"

	"example.com/gocab/pkg"
)

func TestAddCommand(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		expected := &pkg.Word{"german", "Hallo", "Hello", "Hallo, wie gehts", []string{"greetings"}}

		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()))

		received, _ := cmd.Execute([]string{"-lang", "german", "-word", "Hallo", "-meaning", "Hello", "-example", "Hallo, wie gehts", "-tags", "greetings"})

		if !reflect.DeepEqual(expected, received) {
			t.Errorf("expected %v, got %v", expected, received)
		}
	})

	t.Run("required", func(t *testing.T) {
		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()))

		_, err := cmd.Execute([]string{"-word", "Hallo", "-meaning", "Hello", "-example", "Hallo, wie gehts", "-tags", "greetings"})
		if err.Error() != "missing lang" {
			t.Errorf("expected error: missing lang, got: %s", err)
		}

		_, err = cmd.Execute([]string{"-lang", "german", "-meaning", "Hello", "-example", "Hallo, wie gehts", "-tags", "greetings"})
		if err.Error() != "missing word" {
			t.Errorf("expected error: missing word, got: %s", err)
		}

		_, err = cmd.Execute([]string{"-lang", "german", "-word", "Hallo", "-example", "Hallo, wie gehts", "-tags", "greetings"})
		if err.Error() != "missing meaning" {
			t.Errorf("expected error: missing meaning, got: %s", err)
		}
	})
}
