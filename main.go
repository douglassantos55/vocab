package main

import (
	"fmt"
	"os"
	"path"

	"example.com/gocab/pkg"
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

	commands := make(map[string]pkg.Command)
	commands["add"] = pkg.CreateAddCommand(service, pkg.StdWordArgsParser)
	commands["update"] = pkg.CreateUpdateCommand(service, pkg.StdWordArgsParser)
	commands["quiz"] = pkg.CreateQuizCommand(service, pkg.StdQuizArgsParser, os.Stdin, os.Stdout)

	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Println("INVALID COMAND")
		return
	}

	if _, err := command.Execute(os.Args[2:]); err != nil {
		fmt.Println(err)
	}
}
