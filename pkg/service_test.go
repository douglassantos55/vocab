package pkg_test

import (
	"testing"

	"example.com/gocab/pkg"
)

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		service.AddWord("German", "Haus", "house", "", "Dein Haus ist sauber", []string{"noun"})
		service.AddWord("Spanish", "hola", "hello", "", "Hola, hombre", []string{"greeting"})

		if exists, _ := repository.HasWord("German", "Haus"); !exists {
			t.Error("should have word \"Haus\" in German")
		}

		if exists, _ := repository.HasWord("Spanish", "Haus"); exists {
			t.Error("should not have word \"Haus\" in Spanish")
		}
	})

	t.Run("repeated", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		repository.AddWord("German", "Haus", "", "", "", nil)

		service := pkg.NewService(repository)

		_, err := service.AddWord("German", "Haus", "house", "", "Dein Haus ist sauber", []string{"noun"})
		if err != pkg.ErrWordAlreadyRegistered {
			t.Errorf("expected error %v, got %v", pkg.ErrWordAlreadyRegistered, err)
		}
	})

	t.Run("update not registered", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		if _, err := service.UpdateWord("german", "Haus", "House; Home", "", "Ich habe ein Haus", []string{"nouns"}); err != pkg.ErrWordNotRegistered {
			t.Errorf("expected %v, got %v", pkg.ErrWordNotRegistered, err)
		}
	})

	t.Run("update", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		repository.AddWord("german", "Haus", "House", "", "Mein Haus ist blau", []string{"nouns"})

		word, err := service.UpdateWord("german", "Haus", "House; Home", "", "Ich habe ein Haus", []string{"nouns"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if word.Meaning != "House; Home" {
			t.Errorf("should have updated meaning, got %v", word.Meaning)
		}

		if word.Example != "Ich habe ein Haus" {
			t.Errorf("should have updated example, got %v", word.Example)
		}
	})

	t.Run("quiz no words", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		words, err := service.CreateQuiz("stormtrooper", []string{"pronoun"})

		if err != pkg.ErrNoWordsFound {
			t.Errorf("should get error %v, got %v", pkg.ErrNoWordsFound, err)
		}

		if words != nil {
			t.Errorf("should not get any words, got %v", words)
		}
	})

	t.Run("quiz with words", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		repository.AddWord("german", "Er", "He", "", "Er ist mein Mann", []string{"pronoun"})
		repository.AddWord("german", "Mann", "Man; Husband", "", "Mein Mann ist stark", []string{"noun"})
		repository.AddWord("german", "Frau", "Woman; Wife", "", "Mein Frau ist schon", []string{"noun"})
		repository.AddWord("german", "Stark", "Strong", "", "Mein Mann ist stark", []string{"adjective"})

		words, err := service.CreateQuiz("german", []string{"noun", "pronoun"})

		if err != nil {
			t.Errorf("should not get error, got %v", err)
		}

		if len(words) != 3 {
			t.Errorf("should get %v words, got %v", 3, len(words))
		}
	})

	t.Run("quiz without tags", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		repository.AddWord("german", "Er", "He", "", "Er ist mein Mann", []string{"pronoun"})
		repository.AddWord("german", "Mann", "Man; Husband", "", "Mein Mann ist stark", []string{"noun"})
		repository.AddWord("german", "Frau", "Woman; Wife", "", "Mein Frau ist schon", []string{"noun"})
		repository.AddWord("german", "Stark", "Strong", "", "Mein Mann ist stark", []string{"adjective"})

		words, err := service.CreateQuiz("german", []string{})

		if err != nil {
			t.Errorf("should not get error, got %v", err)
		}

		if len(words) != 4 {
			t.Errorf("should get %v words, got %v", 4, len(words))
		}
	})

	t.Run("import same language", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		failed := service.ImportWords([]*pkg.Word{
			{"german", "Hallo", "Hello", "", "Hallo, wie gehts", nil, 0},
			{"german", "Prost", "Cheers", "", "Prost!", nil, 0},
			{"german", "Haus", "House", "", "Mein Haus ist weit weg", nil, 0},
			//{"spanish", "Hombre", "Man", "", "Un belo hombre", nil, 0},
		})

		if len(failed) != 0 {
			t.Fatalf("expected no errors, got %v", failed)
		}

		words, err := repository.FindWords("german", nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(words) != 3 {
			t.Errorf("expected %d words, got %d", 3, len(words))
		}
	})

	t.Run("import different languages", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		failed := service.ImportWords([]*pkg.Word{
			{"german", "Hallo", "Hello", "", "Hallo, wie gehts", nil, 0},
			{"german", "Prost", "Cheers", "", "Prost!", nil, 0},
			{"german", "Haus", "House", "", "Mein Haus ist weit weg", nil, 0},
			{"spanish", "Hombre", "Man", "", "Un belo hombre", nil, 0},
		})

		if len(failed) != 0 {
			t.Fatalf("expected no errors, got %v", failed)
		}

		words, err := repository.FindWords("german", nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(words) != 3 {
			t.Errorf("expected %d words, got %d", 3, len(words))
		}

		words, err = repository.FindWords("spanish", nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(words) != 1 {
			t.Errorf("expected %d words, got %d", 1, len(words))
		}
	})

	t.Run("import existing words", func(t *testing.T) {
		repository := pkg.NewInMemoryRepository()
		service := pkg.NewService(repository)

		repository.AddWord("german", "Er", "He", "", "Er ist mein Mann", []string{"pronoun"})
		repository.AddWord("german", "Mann", "Man; Husband", "", "Mein Mann ist stark", []string{"noun"})
		repository.AddWord("german", "Frau", "Woman; Wife", "", "Mein Frau ist schon", []string{"noun"})
		repository.AddWord("german", "Stark", "Strong", "", "Mein Mann ist stark", []string{"adjective"})

		failed := service.ImportWords([]*pkg.Word{
			{"german", "Er", "Hello", "", "Hallo, wie gehts", nil, 0},
			{"german", "Mann", "Cheers", "", "Prost!", nil, 0},
			{"german", "Frau", "House", "", "Mein Haus ist weit weg", nil, 0},
			{"german", "Stark", "Man", "", "Un belo hombre", nil, 0},
		})

		if len(failed) != 0 {
			t.Fatalf("expected no errors, got %v", failed)
		}

		words, err := repository.FindWords("german", nil)
		if err != nil {
			t.Fatal(err)
		}

		expected := map[string]string{
			"Er":    "Hello",
			"Mann":  "Cheers",
			"Frau":  "House",
			"Stark": "Man",
		}

		for _, word := range words {
			if expected[word.Word] != word.Meaning {
				t.Errorf("expected meaning %s for word %s, got %s", expected[word.Word], word.Word, word.Meaning)
			}
		}
	})
}
