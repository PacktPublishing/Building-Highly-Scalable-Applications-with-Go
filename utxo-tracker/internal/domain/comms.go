package domain

import (
	"context"
	"errors"
)

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

// Query represents a query in the CQRS sense i.e. data that can be requested from the application.
//
// It was created so all such queries have a common interface so that we can utilize polymorphism.
type Query interface {
	// ID is a unique name for the query
	ID() QueryID
}

// QueryID is a unique identifier for a Query type
type QueryID string

func (i QueryID) String() string {
	return string(i)
}

// QueryHandler will execute a Query
type QueryHandler[Q Query, R any] interface {
	// Handle executes the query and calls the provided QueryCallback for each item in the result set. If returns
	// ErrIterationExit when the QueryCallback opted to end iteration early, nil if all good and some opaque error if
	// something else went wrong during inside Handle or inside the provided QueryCallback.
	Handle(ctx context.Context, q Q, callback QueryCallback[R]) error
}

// ErrIterationExit can be returned by a QueryCallback to prematurely end iteration
var ErrIterationExit = errors.New("query: early exit of iteration")

// QueryCallback functions are given to QueryHandler.Handle. It receives the values of the query result set,
// and it should respond with nil to continue iteration, ErrIterationExit to prematurely end iteration or an actual
// error to abort things when something went wrong.
type QueryCallback[R any] func(value R) error
