package pkg

import (
	"bufio"
	"io"
)

type Command interface {
	Execute(args []string) (any, error)
}

type addCommand struct {
	service Service
	parser  AddArgsParser
}

func CreateAddCommand(service Service, parser AddArgsParser) Command {
	return &addCommand{service, parser}
}

func (c *addCommand) Execute(args []string) (any, error) {
	lang, word, meaning, example, tags, err := c.parser(args)
	if err != nil {
		return nil, err
	}
	return c.service.AddWord(lang, word, meaning, example, tags)
}

type updateCommand struct {
	service Service
	parser  UpdateArgsParser
}

func CreateUpdateCommand(service Service, parser UpdateArgsParser) Command {
	return &updateCommand{service, parser}
}

func (c *updateCommand) Execute(args []string) (any, error) {
	lang, word, meaning, example, tags, err := c.parser(args)
	if err != nil {
		return nil, err
	}
	return c.service.UpdateWord(lang, word, meaning, example, tags)
}

type quizCommand struct {
	service Service
	parser  QuizArgsParser
	reader  io.Reader
	writer  io.Writer
}

func CreateQuizCommand(service Service, parser QuizArgsParser, reader io.Reader, writer io.Writer) Command {
	return &quizCommand{service, parser, reader, writer}
}

func (c *quizCommand) Execute(args []string) (any, error) {
	lang, tags, err := c.parser(args)
	if err != nil {
		return nil, err
	}

	questions, err := c.service.CreateQuiz(lang, tags)
	if err != nil {
		return nil, err
	}

	summary, err := c.runQuiz(questions)
	if err != nil {
		return nil, err
	}

	if err := c.service.SaveResult(summary); err != nil {
		return nil, err
	}

	return summary, nil
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
