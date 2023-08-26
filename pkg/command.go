package pkg

import (
	"bufio"
	"fmt"
	"io"
)

type Command interface {
	Execute(args []string) error
}

type addCommand struct {
	service Service

	Lang    string   `short:"l" long:"lang" required:"true" description:"foreign language"`
	Word    string   `short:"w" long:"word" required:"true" description:"foreign word"`
	Meaning string   `short:"m" long:"meaning" required:"true" description:"translation"`
	Tags    []string `short:"t" long:"tags" required:"true" description:"topics of the word"`
	Example string   `short:"e" long:"example" description:"example sentence"`
}

func CreateAddCommand(service Service) *addCommand {
	return &addCommand{service: service}
}

func (c *addCommand) Execute(args []string) error {
	_, err := c.service.AddWord(c.Lang, c.Word, c.Meaning, c.Example, c.Tags)
	return err
}

type updateCommand struct {
	service Service

	Lang    string   `short:"l" long:"lang" required:"true" description:"language of the word to update"`
	Word    string   `short:"w" long:"word" required:"true" description:"word to update"`
	Meaning string   `short:"m" long:"meaning" required:"true" description:"translation"`
	Tags    []string `short:"t" long:"tags" required:"true" description:"topics of the word"`
	Example string   `short:"e" long:"example" description:"example sentence"`
}

func CreateUpdateCommand(service Service) Command {
	return &updateCommand{service: service}
}

func (c *updateCommand) Execute(args []string) error {
	_, err := c.service.UpdateWord(c.Lang, c.Word, c.Meaning, c.Example, c.Tags)
	return err
}

type quizCommand struct {
	service Service
	reader  io.Reader
	writer  io.Writer

	Lang string   `short:"l" long:"lang" required:"true" description:"foreign language"`
	Tags []string `short:"t" long:"tags" description:"topics of the quiz"`
}

func CreateQuizCommand(service Service, reader io.Reader, writer io.Writer) Command {
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
		c.writer.Write([]byte(question.Text()))

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
