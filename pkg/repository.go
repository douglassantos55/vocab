package pkg

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type WordRepository interface {
	HasWord(lang, word string) (bool, error)
	AddWord(lang, word, meaning, example string, tags []string) (*Word, error)
	UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error)
}

type InMemoryRepository struct {
	words map[string]map[string]Word
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		words: make(map[string]map[string]Word),
	}
}

func (r *InMemoryRepository) AddWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	if _, ok := r.words[lang]; !ok {
		r.words[lang] = make(map[string]Word)
	}

	w := Word{lang, word, meaning, example, tags}
	r.words[lang][word] = w

	return &w, nil
}

func (r *InMemoryRepository) UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	words, ok := r.words[lang]
	if !ok {
		return nil, fmt.Errorf("no lang found: %s", lang)
	}

	w := Word{lang, word, meaning, example, tags}
	words[word] = w

	return &w, nil
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

func (r *SqliteRepository) AddWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return nil, err
	}

	insertStmt, err := tx.Prepare("INSERT INTO words (lang, word, meaning, example) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	defer insertStmt.Close()

	result, err := insertStmt.Exec(lang, word, meaning, example)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := r.createTags(tx, id, tags); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Word{lang, word, meaning, example, tags}, nil
}

func (r *SqliteRepository) createTags(tx *sql.Tx, id int64, tags []string) error {
	stmt, err := tx.Prepare("INSERT INTO tags (word_id, tag) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, tag := range tags {
		if _, err := stmt.Exec(id, tag); err != nil {
			return err
		}
	}

	return nil
}

func (r *SqliteRepository) UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("UPDATE words SET meaning = ?, example = ? WHERE lang = ? AND word = ?")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(meaning, example, lang, word)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := r.updateTags(tx, lang, word, tags); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Word{lang, word, meaning, example, tags}, nil
}

func (r *SqliteRepository) updateTags(tx *sql.Tx, lang, word string, tags []string) error {
	_, err := tx.Exec(`
        DELETE FROM tags
        WHERE word_id = (SELECT id FROM words WHERE lang = ? AND word = ?)
    `, lang, word)

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO tags (word_id, tag)
        SELECT id, ? FROM words WHERE lang = ? AND word = ?
    `)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, tag := range tags {
		if _, err := stmt.Exec(tag, lang, word); err != nil {
			return err
		}
	}

	return nil
}
		return nil, err
	}

	return &Word{lang, word, meaning, example, tags}, nil
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
