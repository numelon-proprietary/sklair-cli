package commandRegistry

// adapted from CommandRegistry intended for a bot in another numelon-proprietary project

import "fmt"

type Command struct {
	Name        string
	Description string
	Aliases     []string
	Run         func(args []string) int
}

type CommandRegistry struct {
	commands map[string]*Command
}

func New() *CommandRegistry {
	r := &CommandRegistry{
		commands: make(map[string]*Command),
	}

	r.Register(&Command{
		Name:        "help",
		Description: "Shows available commands",
		Run: func(args []string) int {
			r.PrintHelp()
			return 0
		},
	})

	return r
}

func (r *CommandRegistry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

func (r *CommandRegistry) Get(name string) (*Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *CommandRegistry) PrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("	sklair <command> [args]")
	fmt.Println()
	fmt.Println("Available commands:")

	seen := make(map[*Command]bool)
	for _, cmd := range r.commands {
		if seen[cmd] {
			continue
		}
		seen[cmd] = true

		fmt.Printf("  %-12s %s\n", cmd.Name, cmd.Description)
	}
}
