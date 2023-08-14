package pkg_test

import (
	"testing"

	"example.com/gocab/pkg"
)

type fakeRepository struct {
	data map[string]pkg.Word
}

func (f *fakeRepository) AddWord(word pkg.Word) (*pkg.Word, error) {
	f.data[word.Word] = word
	return &word, nil
}

func (f *fakeRepository) HasWord(word string) bool {
	_, found := f.data[word]
	return found
}

func NewFakeRepository() *fakeRepository {
	return &fakeRepository{
		data: make(map[string]pkg.Word),
	}
}

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		repository := NewFakeRepository()
		service := pkg.NewService(repository)

		service.AddWord("German", "Haus", "house", "Dein Haus ist sauber", []string{"noun"})
		service.AddWord("Spanish", "hola", "hello", "Hola, hombre", []string{"greeting"})

		if !repository.HasWord("German", "Haus") {
			t.Error("should have word \"Haus\" in German")
		}

		if repository.HasWord("Spanish", "Haus") {
			t.Error("should not have word \"Haus\" in Spanish")
		}
	})

	t.Run("repeated", func(t *testing.T) {
		repository := NewFakeRepository()
		repository.AddWord(pkg.Word{Lang: "German", Word: "Haus"})

		service := pkg.NewService(repository)

		_, err := service.AddWord("German", "Haus", "house", "Dein Haus ist sauber", []string{"noun"})
		if err != pkg.ErrWordAlreadyRegistered {
			t.Errorf("expected error %v, got %v", pkg.ErrWordAlreadyRegistered, err)
		}
	})
}
