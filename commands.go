package main

import "fmt"

type Command struct{
	name string
	args []string
}

type Commands struct{
	handlers map[string]func(*state, Command) error
}

func (c *Commands) run(s *state, cmd Command) error{
	handler, ok := c.handlers[cmd.name]; if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}
	err := handler(s, cmd)
	return err
}

func (c *Commands) register(name string, f func(*state, Command) error){
	c.handlers[name] = f
}