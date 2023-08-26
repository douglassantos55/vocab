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
		return fmt.Sprintf("[%s] How do you say \"%s\" in %s\n", q.Word.Level(), q.Word.Meaning, q.Word.Lang)
	}

	if q.Word.Pronunciation != "" {
		return fmt.Sprintf("[%s] What does %s [%s] mean?\n", q.Word.Level(), q.Word.Word, q.Word.Pronunciation)
	}

	return fmt.Sprintf("[%s] What does %s mean?\n", q.Word.Level(), q.Word.Word)
}

func (q *Question) ExpectedAnswer() string {
	if q.Type == 0 {
		return q.Word.Word
	}
	return q.Word.Meaning
}

func (q *Question) IsCorrect() bool {
	for _, meaning := range strings.Split(q.ExpectedAnswer(), ";") {
		if strings.TrimSpace(strings.ToLower(q.Answer)) == strings.TrimSpace(strings.ToLower(meaning)) {
			return true
		}
	}
	return false
}

type Summary struct {
	Total     int
	Mistakes  int
	Questions []*Question
}

func (s *Summary) Correct(question *Question) {
	s.Questions = append(s.Questions, question)
}

func (s *Summary) Wrong(question *Question) {
	s.Mistakes++
	s.Questions = append(s.Questions, question)
}

func (s *Summary) String() string {
	correct := s.Total - s.Mistakes
	performance := (1 - float64(s.Mistakes)/float64(s.Total)) * 100

	str := fmt.Sprintf("\nTotal: %d, Correct: %d, Mistakes: %d, Performance: %.0f%%\n", s.Total, correct, s.Mistakes, performance)

	if s.Mistakes > 0 {
		for _, question := range s.Questions {
			if !question.IsCorrect() {
				str += fmt.Sprintf("%s -> %s\n", strings.TrimSpace(question.Answer), question.ExpectedAnswer())
			}
		}
	}

	return str
}

type Word struct {
	Lang          string
	Word          string
	Meaning       string
	Pronunciation string
	Example       string
	Tags          []string
	Score         float64
}

func (w *Word) Level() string {
	if w.Score < 0.5 {
		return "Hard"
	} else if w.Score < 1 {
		return "Medium"
	}
	return "Easy"
}

type Service interface {
	AddWord(lang, word, meaning, pronunciation, example string, tags []string) (*Word, error)
	UpdateWord(lang, word, meaning, pronunciation, example string, tags []string) (*Word, error)
	CreateQuiz(lang string, tags []string) ([]*Question, error)
	SaveResult(summary *Summary) error
}

type service struct {
	repository WordRepository
}

func NewService(repository WordRepository) *service {
	return &service{repository}
}

func (s *service) AddWord(lang, word, meaning, pronunciation, example string, tags []string) (*Word, error) {
	exists, err := s.repository.HasWord(lang, word)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrWordAlreadyRegistered
	}

	return s.repository.AddWord(lang, word, meaning, pronunciation, example, tags)
}

func (s *service) UpdateWord(lang, word, meaning, pronunciation, example string, tags []string) (*Word, error) {
	exists, err := s.repository.HasWord(lang, word)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrWordNotRegistered
	}

	return s.repository.UpdateWord(lang, word, meaning, pronunciation, example, tags)
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
	return s.repository.SaveResult(summary)
}
