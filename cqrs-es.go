package cqrs_es

import (
    "log"
    "os"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[cqrs-es] ", 0)
}

type Store interface {
    Events(EntityId) ([]Event, error)
    Commit(EntityId, []Event) error
}

type Event interface {}

type Command interface {
    // This Execute should not modify the original array passed in.
    // Returns a set of new events or an error if the command is invalid
    Execute([]Event) ([]Event, error)
}

type EntityId string

type EventSerializer interface {
    Serialize(Event) (*EventRecord, error)
    Deserialize(*EventRecord) (Event, error)
}

type EventRecord struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

