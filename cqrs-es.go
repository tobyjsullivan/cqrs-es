package cqrs_es

type Event interface {}

type Command interface {
    // This Execute should not modify the original array passed in.
    // Returns a set of new events or an error if the command is invalid
    Execute([]Event) ([]Event, error)
}

type EntityId string

type Service interface {
    Execute(EntityId, Command) error
    Events(EntityId, uint) []Event
}


