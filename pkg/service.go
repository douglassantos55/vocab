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

type Service interface {
	AddWord(lang, word, meaning, example string, tags []string) (*Word, error)
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

	w := Word{lang, word, meaning, example, tags}
	return s.repository.AddWord(w)
}

func NewService(repository WordRepository) *service {
	return &service{repository}
}
