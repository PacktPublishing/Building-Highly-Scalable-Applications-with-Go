package domain

import "context"

// Command represents a command in the CQRS sense i.e. commands that can be sent to the application to change state.
//
// It was created so all such commands have a common interface so that we can utilize polymorphism.
type Command interface {
	// ID is a unique name for the command
	ID() CommandID
}

// CommandID is a unique identifier for a Command type
type CommandID string

func (i CommandID) String() string {
	return string(i)
}

// CommandHandler will execute a Command
type CommandHandler[C Command] interface {
	Handle(ctx context.Context, command C) error
}

// CommandHandlerFunc satisfies the CommandHandler interface and so allows a normal function to become a CommandHandler.
type CommandHandlerFunc[C Command] func(ctx context.Context, command C) error

func (f CommandHandlerFunc[C]) Handle(ctx context.Context, c C) error {
	return f(ctx, c)
}
