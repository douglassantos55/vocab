package pkg_test

import (
	"testing"

	"example.com/gocab/pkg"
)

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		service.AddWord("German", "Haus", "house", "Dein Haus ist sauber", []string{"noun"})
		service.AddWord("Spanish", "hola", "hello", "Hola, hombre", []string{"greeting"})

		if exists, _ := repository.HasWord("German", "Haus"); !exists {
			t.Error("should have word \"Haus\" in German")
		}

		if exists, _ := repository.HasWord("Spanish", "Haus"); exists {
			t.Error("should not have word \"Haus\" in Spanish")
		}
	})

	t.Run("repeated", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		repository.AddWord(pkg.Word{Lang: "German", Word: "Haus"})

		service := pkg.NewService(repository)

		_, err := service.AddWord("German", "Haus", "house", "Dein Haus ist sauber", []string{"noun"})
		if err != pkg.ErrWordAlreadyRegistered {
			t.Errorf("expected error %v, got %v", pkg.ErrWordAlreadyRegistered, err)
		}
	})
}
