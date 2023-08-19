package pkg

import (
	"flag"
	"fmt"
	"strings"
)

type AddArgsParser func(args []string) (string, string, string, string, []string, error)
type UpdateArgsParser func(args []string) (string, string, string, string, []string, error)

func StdWordArgsParser(args []string) (string, string, string, string, []string, error) {
	flagset := flag.NewFlagSet("add", flag.ExitOnError)

	lang := flagset.String("l", "", "foreign language")
	word := flagset.String("w", "", "foreign word")
	meaning := flagset.String("m", "", "translation")
	tags := flagset.String("t", "", "comma-separated list of tags")
	example := flagset.String("e", "", "example sentence")

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
