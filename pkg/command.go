package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type WordCommand struct {
	Lang    string   `short:"l" long:"lang" required:"true" description:"foreign language"`
	Word    string   `short:"w" long:"word" required:"true" description:"foreign word"`
	Meaning string   `short:"m" long:"meaning" required:"true" description:"translation"`
	Tags    []string `short:"t" long:"tags" required:"true" description:"topics of the word"`

	Pronunciation string `short:"p" long:"pronunciation" description:"how to pronounce the word"`
	Example       string `short:"e" long:"example" description:"example sentence"`
}

type addCommand struct {
	WordCommand
	service Service
}

func CreateAddCommand(service Service) *addCommand {
	return &addCommand{service: service}
}

func (c *addCommand) Execute(args []string) error {
	_, err := c.service.AddWord(c.Lang, c.Word, c.Meaning, c.Pronunciation, c.Example, c.Tags)
	return err
}

type updateCommand struct {
	WordCommand
	service Service
}

func CreateUpdateCommand(service Service) *updateCommand {
	return &updateCommand{service: service}
}

func (c *updateCommand) Execute(args []string) error {
	_, err := c.service.UpdateWord(c.Lang, c.Word, c.Meaning, c.Pronunciation, c.Example, c.Tags)
	return err
}

type quizCommand struct {
	service Service
	reader  io.Reader
	writer  io.Writer

	Lang string   `short:"l" long:"lang" required:"true" description:"foreign language"`
	Tags []string `short:"t" long:"tags" description:"topics of the quiz"`
}

func CreateQuizCommand(service Service, reader io.Reader, writer io.Writer) *quizCommand {
	return &quizCommand{service: service, reader: reader, writer: writer}
}

func (c *quizCommand) Execute(args []string) error {
	questions, err := c.service.CreateQuiz(c.Lang, c.Tags)
	if err != nil {
		return err
	}

	summary, err := c.runQuiz(questions)
	if err != nil {
		return err
	}

	if err := c.service.SaveResult(summary); err != nil {
		return err
	}

	_, err = fmt.Fprintln(c.writer, summary)

	return err
}

func (c *quizCommand) runQuiz(questions []*Question) (*Summary, error) {
	reader := bufio.NewReader(c.reader)
	summary := &Summary{Total: len(questions)}

	for _, question := range questions {
		_, err := c.writer.Write([]byte(question.Text()))
		if err != nil {
			return nil, err
		}

		answer, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// Set question's answer
		question.Answer = answer

		if question.IsCorrect() {
			summary.Correct(question)
		} else {
			summary.Wrong(question)
		}
	}

	return summary, nil
}

type importCommand struct {
	writer  io.Writer
	service Service

	Language string `short:"l" long:"lang" required:"true" description:"foreign language"`
	Filename string `short:"f" long:"file" required:"true" description:"csv file containing words to import"`
}

func CreateImportCommand(service Service, writer io.Writer) *importCommand {
	return &importCommand{
		writer:  writer,
		service: service,
	}
}

func (c *importCommand) Execute(args []string) error {
	file, err := os.OpenFile(c.Filename, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	words := make([]*Word, 0)
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i > 0 {
			word, meaning, pronunciation, example, tags := c.parseLine(line)
			if word == "" {
				continue
			}

			words = append(words, &Word{
				Lang:          c.Language,
				Word:          word,
				Meaning:       meaning,
				Pronunciation: pronunciation,
				Example:       example,
				Tags:          strings.Split(tags, ","),
			})
		}
	}

	failed := c.service.ImportWords(words)
	if len(failed) != 0 {
		c.writer.Write([]byte("could not import words:\n"))
		for word, reason := range failed {
			c.writer.Write([]byte(fmt.Sprintf("%s: %s\n", word, reason)))
		}
	}

	return nil
}

func (c *importCommand) parseLine(line string) (string, string, string, string, string) {
	parts := strings.Split(line, ";")
	if len(parts) < 5 {
		return "", "", "", "", ""
	}
	return parts[0], parts[1], parts[2], parts[3], parts[4]
}
