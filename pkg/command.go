package pkg

import (
	"flag"
	"strings"
)

type Command interface {
	Execute(args []string) (any, error)
}

type addCommand struct {
	service Service
}

func (c *addCommand) parseArgs(args []string) (string, string, string, string, []string, error) {
	flagset := flag.NewFlagSet("add", flag.ExitOnError)

	lang := flagset.String("lang", "", "foreign language")
	word := flagset.String("word", "", "foreign word")
	meaning := flagset.String("meaning", "", "translation")
	example := flagset.String("example", "", "usage example sentence")
	tags := flagset.String("tags", "", "comma-separated list of tags")

	if err := flagset.Parse(args); err != nil {
		return "", "", "", "", []string{}, err
	}

	return *lang, *word, *meaning, *example, strings.Split(*tags, ","), nil
}

func (c *addCommand) Execute(args []string) (any, error) {
	lang, word, meaning, example, tags, err := c.parseArgs(args)
	if err != nil {
		return nil, err
	}
	return c.service.AddWord(lang, word, meaning, example, tags)
}

func CreateAddCommand(service Service) Command {
	return &addCommand{service}
}
