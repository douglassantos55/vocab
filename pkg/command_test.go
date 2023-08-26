package pkg_test

import (
	"testing"

	"example.com/gocab/pkg"
	"github.com/jessevdk/go-flags"
)

func TestAddCommand(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()))

		_, err := flags.ParseArgs(cmd, []string{"-l", "german", "-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err != nil {
			t.Errorf("should add word, got error %v", err)
		}
	})

	t.Run("required", func(t *testing.T) {
		cmd := pkg.CreateAddCommand(pkg.NewService(pkg.NewInMemoryRepository()))

		_, err := flags.ParseArgs(cmd, []string{"-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err == nil {
			t.Fatal("should error, no lang provided")
		}

		_, err = flags.ParseArgs(cmd, []string{"-l", "german", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err == nil {
			t.Fatal("should error, no word provided")
		}

		_, err = flags.ParseArgs(cmd, []string{"-l", "german", "-w", "Hallo", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err == nil {
			t.Fatal("should error, no meaning provided")
		}

		_, err = flags.ParseArgs(cmd, []string{"-l", "german", "-w", "Hallo", "-m", "Hello", "-t", "greetings"})
		if err != nil {
			t.Fatalf("should not error, got %v", err)
		}
	})
}
