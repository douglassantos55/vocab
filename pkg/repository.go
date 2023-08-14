package pkg

type WordRepository interface {
	AddWord(word Word) (*Word, error)
	HasWord(lang, word string) bool
}

type InMemoryRepository struct {
	words map[string]map[string]Word
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		words: make(map[string]map[string]Word),
	}
}

func (r *InMemoryRepository) AddWord(word Word) (*Word, error) {
	if _, ok := r.words[word.Lang]; !ok {
		r.words[word.Lang] = make(map[string]Word)
	}
	r.words[word.Lang][word.Word] = word
	return &word, nil
}

func (r *InMemoryRepository) HasWord(lang, word string) bool {
	if list, ok := r.words[lang]; ok {
		_, found := list[word]
		return found
	}
	return false
}
