package pkg

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	ErrWordAlreadyRegistered = errors.New("word already registered")
	ErrWordNotRegistered     = errors.New("word not registered")
	ErrNoWordsFound          = errors.New("no words found")
)

type Question struct {
	Type    int
	Word    *Word
	Correct bool
}

func NewQuestion(word *Word) *Question {
	rand.Seed(time.Now().UnixNano())
	return &Question{Type: rand.Intn(2), Word: word}
}

func (q *Question) Answer(word string) bool {
	if q.Type == 0 {
		q.Correct = strings.TrimSpace(word) == strings.TrimSpace(q.Word.Word)
	} else {
		q.Correct = strings.TrimSpace(word) == strings.TrimSpace(q.Word.Meaning)
	}
	return q.Correct
}

func (q *Question) Text() string {
	if q.Type == 0 {
		return fmt.Sprintf("How do you say \"%s\" in %s\n", q.Word.Meaning, q.Word.Lang)
	}
	return fmt.Sprintf("What does %s mean?\n", q.Word.Word)
}

func (q *Question) IsCorrect() bool {
	return q.Correct
}
type Word struct {
	Lang    string
	Word    string
	Meaning string
	Example string
	Tags    []string
	Score   float64
}

type Service interface {
	AddWord(lang, word, meaning, example string, tags []string) (*Word, error)
	UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error)
	CreateQuiz(lang string, tags []string) ([]*Question, error)
}

type service struct {
	repository WordRepository
}

func (s *service) AddWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	exists, err := s.repository.HasWord(lang, word)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrWordAlreadyRegistered
	}

	return s.repository.AddWord(lang, word, meaning, example, tags)
}

func (s *service) UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	exists, err := s.repository.HasWord(lang, word)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrWordNotRegistered
	}

	return s.repository.UpdateWord(lang, word, meaning, example, tags)
}

func (s *service) CreateQuiz(lang string, tags []string) ([]*Question, error) {
	words, err := s.repository.FindWords(lang, tags)
	if err != nil {
		return nil, err
	}

	if len(words) == 0 {
		return nil, ErrNoWordsFound
	}

	questions := make([]*Question, 0)
	for _, word := range words {
		questions = append(questions, NewQuestion(word))
	}

	return questions, nil
}
}

func NewService(repository WordRepository) *service {
	return &service{repository}
}
