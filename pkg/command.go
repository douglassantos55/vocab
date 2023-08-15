package pkg

import (
	"flag"
	"fmt"
	"strings"
)

type Command interface {
	Execute(args []string) (any, error)
}

type AddArgsParser func(args []string) (string, string, string, string, []string, error)

func StdAddArgsParser(args []string) (string, string, string, string, []string, error) {
	flagset := flag.NewFlagSet("add", flag.ExitOnError)

	lang := flagset.String("lang", "", "foreign language")
	word := flagset.String("word", "", "foreign word")
	meaning := flagset.String("meaning", "", "translation")
	tags := flagset.String("tags", "", "comma-separated list of tags")
	example := flagset.String("example", "", "example sentence")

	if err := flagset.Parse(args); err != nil {
		return "", "", "", "", nil, err
	}

	if *lang == "" {
		return "", "", "", "", nil, fmt.Errorf("missing lang")
	}

	if *word == "" {
		return "", "", "", "", nil, fmt.Errorf("missing word")
	}

	if *meaning == "" {
		return "", "", "", "", nil, fmt.Errorf("missing meaning")
	}

	tagList := strings.Split(*tags, ",")
	return *lang, *word, *meaning, *example, tagList, nil
}

type addCommand struct {
	service Service
	parser  AddArgsParser
}

func (c *addCommand) Execute(args []string) (any, error) {
	lang, word, meaning, example, tags, err := c.parser(args)
	if err != nil {
		return nil, err
	}
	return c.service.AddWord(lang, word, meaning, example, tags)
}

func CreateAddCommand(service Service, parser AddArgsParser) Command {
	return &addCommand{service, parser}
}
