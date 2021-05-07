package framework

// Represents a command registered via function with context
type Command func(Context)

type CommandStruct struct {
	Command Command
	Help    string
}

type CommandMap map[string]CommandStruct

type CommandHandler struct {
	Cmds CommandMap
}

// Creates a new instance of CommandHandler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CommandMap)}
}

// retrieves the list of registered commands
// from the CommandHandler registrar
func (handler CommandHandler) GetCommands() CommandMap {
	return handler.Cmds
}

// Get a command from the CommandHandler registrar (if one exists)
func (handler CommandHandler) Get(name string) (*Command, bool) {
	cmd, found := handler.Cmds[name]

	return &cmd.Command, found
}

// Register a command to the CommandHandler registrar, using a name, CommandContext,
// and help message for the help command
func (handler CommandHandler) Register(name string, command Command, helpMessage string) {
	cmdStruct := CommandStruct{Command: command, Help: helpMessage}
	handler.Cmds[name] = cmdStruct
	if len(name) > 1 {
		handler.Cmds[name[:1]] = cmdStruct
	}
}

// Retrieves the help message from a command out of the registrar
func (command CommandStruct) GetHelp() string {
	return command.Help
}