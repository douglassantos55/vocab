package pkg

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type WordRepository interface {
	AddWord(word Word) (*Word, error)
	HasWord(lang, word string) (bool, error)
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

func (r *InMemoryRepository) HasWord(lang, word string) (bool, error) {
	if list, ok := r.words[lang]; ok {
		_, found := list[word]
		return found, nil
	}
	return false, nil
}

type SqliteRepository struct {
	conn *sql.DB
}

func NewSqliteRepository(filename string) (*SqliteRepository, error) {
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return &SqliteRepository{conn}, nil
}

func (r *SqliteRepository) Close() {
	r.conn.Close()
}

func (r *SqliteRepository) AddWord(word Word) (*Word, error) {
	stmt, err := r.conn.Prepare("INSERT INTO words (lang, word, meaning, example) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(word.Lang, word.Word, word.Meaning, word.Example)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (r *SqliteRepository) HasWord(lang, word string) (bool, error) {
	stmt, err := r.conn.Prepare("SELECT COUNT(*) FROM words WHERE lang = ? AND word = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(lang, word)
	if err != nil {
		return false, err
	}

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}
