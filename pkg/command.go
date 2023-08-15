package pkg

import (
	"flag"
	"fmt"
	"strings"
)

type Command interface {
	Execute(args []string) (any, error)
	ParseArgs(args []string) (any, error)
}

type addCommand struct {
	service Service
}

func (c *addCommand) ParseArgs(args []string) (any, error) {
	flagset := flag.NewFlagSet("add", flag.ExitOnError)

	lang := flagset.String("lang", "", "foreign language")
	word := flagset.String("word", "", "foreign word")
	meaning := flagset.String("meaning", "", "translation")
	tags := flagset.String("tags", "", "comma-separated list of tags")
	example := flagset.String("example", "", "example sentence")

	if err := flagset.Parse(args); err != nil {
		return nil, err
	}

	if *lang == "" {
		return nil, fmt.Errorf("missing lang")
	}
	if *word == "" {
		return nil, fmt.Errorf("missing word")
	}
	if *meaning == "" {
		return nil, fmt.Errorf("missing meaning")
	}

	return Word{*lang, *word, *meaning, *example, strings.Split(*tags, ",")}, nil
}

func (c *addCommand) Execute(args []string) (any, error) {
	values, err := c.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	word := values.(Word)
	return c.service.AddWord(word.Lang, word.Word, word.Meaning, word.Example, word.Tags)
}

func CreateAddCommand(service Service) Command {
	return &addCommand{service}
}
