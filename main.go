package main

import (
	"fmt"
	"os"
	"path"

	"example.com/gocab/pkg"
	"github.com/jessevdk/go-flags"
)

func main() {
	if len(os.Args) < 2 {
		print("expected 'add', 'update' or 'quiz' subcommands")
		os.Exit(1)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	databaseDir := path.Join(configDir, "gocab")
	_, err = os.ReadDir(databaseDir)
	if err != nil {
		if err := os.Mkdir(databaseDir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	filename := path.Join(databaseDir, "database.db")
	repository, err := pkg.NewSqliteRepository(fmt.Sprintf("file:%s?cache=shared", filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer repository.Close()
	service := pkg.NewService(repository)

	parser := flags.NewNamedParser("gocab", flags.Default)

	addCommand := pkg.CreateAddCommand(service)
	updateCommand := pkg.CreateUpdateCommand(service)
	quizCommand := pkg.CreateQuizCommand(service, os.Stdin, os.Stdout)

	parser.AddCommand("add", "add new word", "", addCommand)
	parser.AddCommand("update", "update word", "", updateCommand)
	parser.AddCommand("quiz", "start quiz", "", quizCommand)

	parser.Parse()
}
