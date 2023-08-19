package pkg

type Command interface {
	Execute(args []string) (any, error)
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
