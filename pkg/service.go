package pkg

import "errors"

var ErrWordAlreadyRegistered = errors.New("word already registered")

type Word struct {
	Lang    string
	Word    string
	Meaning string
	Example string
	Tags    []string
}

type WordRepository interface {
	AddWord(word Word) (*Word, error)
	HasWord(word string) bool
}

type service struct {
	repository WordRepository
}

func (s *service) AddWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	if s.repository.HasWord(lang, word) {
		return nil, ErrWordAlreadyRegistered
	}
	w := Word{lang, word, meaning, example, tags}
	return s.repository.AddWord(w)
}

func NewService(repository WordRepository) *service {
	return &service{repository}
}
