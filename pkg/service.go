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
	Type   int
	Word   *Word
	Answer string
}

func NewQuestion(word *Word) *Question {
	rand.Seed(time.Now().UnixNano())
	return &Question{Type: rand.Intn(2), Word: word}
}

func (q *Question) Text() string {
	if q.Type == 0 {
		return fmt.Sprintf("How do you say \"%s\" in %s\n", q.Word.Meaning, q.Word.Lang)
	}
	return fmt.Sprintf("What does %s mean?\n", q.Word.Word)
}

func (q *Question) IsCorrect() bool {
	if q.Type == 0 {
		return strings.TrimSpace(q.Answer) == strings.TrimSpace(q.Word.Word)
	} else {
		for _, meaning := range strings.Split(q.Word.Meaning, ";") {
			if strings.TrimSpace(q.Answer) == strings.TrimSpace(meaning) {
				return true
			}
		}
		return false
	}
}

type Summary struct {
	wrong   []*Question
	correct []*Question
	Total   int
}

func (s *Summary) Correct(question *Question) {
	s.correct = append(s.correct, question)
}

func (s *Summary) Wrong(question *Question) {
	s.wrong = append(s.wrong, question)
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
	SaveResult(summary *Summary) error
}

type service struct {
	repository WordRepository
}

func NewService(repository WordRepository) *service {
	return &service{repository}
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

func (s *service) SaveResult(summary *Summary) error {
	return nil
}
