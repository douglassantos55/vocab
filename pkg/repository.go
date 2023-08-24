package pkg

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type WordRepository interface {
	HasWord(lang, word string) (bool, error)
	FindWords(lang string, tags []string) ([]*Word, error)
	AddWord(lang, word, meaning, example string, tags []string) (*Word, error)
	UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error)
	SaveResult(summary *Summary) error
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

	w := Word{lang, word, meaning, example, tags, 0}
	r.words[lang][word] = w

	return &w, nil
}

func (r *InMemoryRepository) UpdateWord(lang, word, meaning, example string, tags []string) (*Word, error) {
	words, ok := r.words[lang]
	if !ok {
		return nil, fmt.Errorf("no lang found: %s", lang)
	}

	w := Word{lang, word, meaning, example, tags, 0}
	words[word] = w

	return &w, nil
}

func (r *InMemoryRepository) FindWords(lang string, tags []string) ([]*Word, error) {
	words, ok := r.words[lang]
	if !ok {
		return nil, nil
	}

	found := make([]*Word, 0)

	for _, word := range words {
	outer:
		for _, tag := range word.Tags {
			for _, expected := range tags {
				if tag == expected {
					found = append(found, &word)
					break outer
				}
			}
		}
	}

	return found, nil
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

	return &Word{lang, word, meaning, example, tags, 0}, nil
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

	return &Word{lang, word, meaning, example, tags, 0}, nil
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

func (r *SqliteRepository) FindWords(lang string, tags []string) ([]*Word, error) {
	var err error
	var rows *sql.Rows

	if len(tags) > 0 {
		rows, err = r.findWordsForTags(lang, tags)
	} else {
		rows, err = r.findWords(lang)
	}

	if err != nil {
		return nil, err
	}

	var words []*Word

	for rows.Next() {
		var lang string
		var word string
		var meaning string
		var example string
		var score float64

		if err := rows.Scan(&lang, &word, &meaning, &example, &score); err != nil {
			return nil, err
		}

		words = append(words, &Word{lang, word, meaning, example, nil, score})
	}

	return words, nil
}

func (r *SqliteRepository) findWords(lang string) (*sql.Rows, error) {
	return r.conn.Query(`
        -- Hard
        SELECT * FROM (
            SELECT lang, word, meaning, example, score
            FROM words
            WHERE lang = ? AND score = 0
            ORDER BY RANDOM()
            LIMIT 5
        )
        UNION

        -- Medium
        SELECT * FROM (
            SELECT lang, word, meaning, example, score
            FROM words
            WHERE lang = ? AND score = 0.5
            ORDER BY RANDOM()
            LIMIT 5
        )
        UNION

        -- Easy
        SELECT * FROM (
            SELECT lang, word, meaning, example, score
            FROM words
            WHERE lang = ? AND score = 1
            ORDER BY RANDOM()
            LIMIT 5
        )
    `, lang, lang, lang)
}

func (r *SqliteRepository) findWordsForTags(lang string, tags []string) (*sql.Rows, error) {
	args := make([]any, 0)

	for i := 0; i < 3; i++ {
		args = append(args, lang)
		for _, tag := range tags {
			args = append(args, tag)
		}
	}

	return r.conn.Query(`
        -- Hard
        SELECT * FROM (
            SELECT lang, word, meaning, example, score FROM words
            WHERE lang = ? AND score = 0 AND id IN (SELECT word_id FROM tags WHERE tag IN (?`+strings.Repeat(",?", len(tags)-1)+`))
            ORDER BY RANDOM()
            LIMIT 5
        )
        UNION

        -- Medium
        SELECT * FROM (
            SELECT lang, word, meaning, example, score FROM words
            WHERE lang = ? AND score = 0.5 AND id IN (SELECT word_id FROM tags WHERE tag IN (?`+strings.Repeat(",?", len(tags)-1)+`))
            ORDER BY RANDOM()
            LIMIT 5
        )
        UNION

        -- Easy
        SELECT * FROM (
            SELECT lang, word, meaning, example, score FROM words
            WHERE lang = ? AND score = 1 AND id IN (SELECT word_id FROM tags WHERE tag IN (?`+strings.Repeat(",?", len(tags)-1)+`))
            ORDER BY RANDOM()
            LIMIT 5
        )
    `, args...)
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

func (r *SqliteRepository) SaveResult(summary *Summary) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	for _, question := range summary.Questions {
		if question.IsCorrect() {
			_, err := tx.Exec(`
                UPDATE words SET score = score + 0.5
                WHERE lang = ? AND word = ? AND score < 1
            `, question.Word.Lang, question.Word.Word)

			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err := tx.Exec(`
                UPDATE words SET score = score - 0.5
                WHERE lang = ? AND word = ? AND score > 0
            `, question.Word.Lang, question.Word.Word)

			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
