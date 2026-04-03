package command

type CommandRegistry struct{}

func NewRegistry() *CommandRegistry {
    reg := CommandRegistry{}
    return &reg
}
