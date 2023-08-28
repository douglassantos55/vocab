package pkg_test

import (
	"bytes"
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

func TestUpdateCommand(t *testing.T) {
	t.Run("update not registered", func(t *testing.T) {
		cmd := pkg.CreateUpdateCommand(pkg.NewService(pkg.NewInMemoryRepository()))

		_, err := flags.ParseArgs(cmd, []string{"-l", "german", "-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = cmd.Execute([]string{})
		if err != pkg.ErrWordNotRegistered {
			t.Errorf("expected error %v, got %v", pkg.ErrWordNotRegistered, err)
		}
	})

	t.Run("update", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		cmd := pkg.CreateUpdateCommand(pkg.NewService(repository))

		repository.AddWord("german", "Hallo", "Hello", "", "", []string{})

		_, err := flags.ParseArgs(cmd, []string{"-l", "german", "-w", "Hallo", "-m", "Hello", "-e", "Hallo, wie gehts", "-t", "greetings"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := cmd.Execute([]string{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		words, _ := repository.FindWords("german", []string{})

		if words[0].Example != "Hallo, wie gehts" {
			t.Errorf("expected example %s, got %s", "Hallo, wie gehts", words[0].Example)
		}

		if words[0].Tags[0] != "greetings" {
			t.Errorf("expected tag %s, got %s", "greetings", words[0].Tags[0])
		}
	})

	t.Run("required", func(t *testing.T) {
		cmd := pkg.CreateUpdateCommand(pkg.NewService(pkg.NewInMemoryRepository()))

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

func TestQuizCommand(t *testing.T) {
	t.Run("no words", func(t *testing.T) {
		reader := bytes.NewBuffer(nil)
		writer := bytes.NewBuffer(nil)

		cmd := pkg.CreateQuizCommand(pkg.NewService(pkg.NewInMemoryRepository()), reader, writer)

		_, err := flags.ParseArgs(cmd, []string{"-l", "german"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := cmd.Execute([]string{}); err != pkg.ErrNoWordsFound {
			t.Errorf("expected error %v, got %v", pkg.ErrNoWordsFound, err)
		}
	})

	t.Run("quiz", func(t *testing.T) {
		reader := bytes.NewBuffer([]byte("Hello\n"))
		writer := bytes.NewBuffer([]byte(""))

		repository := pkg.NewInMemoryRepository()
		cmd := pkg.CreateQuizCommand(pkg.NewService(repository), reader, writer)

		repository.AddWord("german", "Hallo", "Hello", "", "", []string{})

		_, err := flags.ParseArgs(cmd, []string{"-l", "german"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := cmd.Execute([]string{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		words, err := repository.FindWords("german", []string{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if words[0].Score != 0.5 {
			t.Errorf("expected score %f, got %f", 0.5, words[0].Score)
		}
	})
}
