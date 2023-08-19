package pkg_test

import (
	"reflect"
	"testing"

	"example.com/gocab/pkg"
)

func TestAddCommand(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		expected := &pkg.Word{"german", "Hallo", "Hello", "Hallo, wie gehts", []string{"greetings"}}

		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()), pkg.StdWordArgsParser)

		received, _ := cmd.Execute([]string{"-l", "german", "-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})

		if !reflect.DeepEqual(expected, received) {
			t.Errorf("expected %v, got %v", expected, received)
		}
	})

	t.Run("required", func(t *testing.T) {
		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()), pkg.StdWordArgsParser)

		_, err := cmd.Execute([]string{"-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err.Error() != "missing lang" {
			t.Errorf("expected error: missing lang, got: %s", err)
		}

		_, err = cmd.Execute([]string{"-l", "german", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err.Error() != "missing word" {
			t.Errorf("expected error: missing word, got: %s", err)
		}

		_, err = cmd.Execute([]string{"-l", "german", "-w", "Hallo", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err.Error() != "missing meaning" {
			t.Errorf("expected error: missing meaning, got: %s", err)
		}
	})
}
