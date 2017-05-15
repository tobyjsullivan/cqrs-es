package store

import (
    "github.com/tobyjsullivan/cqrs-es"
    "encoding/json"
)

type testEvent struct {
    Content string
}

type testSerializer struct {}

func (s *testSerializer) Serialize(event cqrs_es.Event) (*cqrs_es.EventRecord, error) {
    t := ""
    switch event.(type) {
    case *testEvent:
        t = "testEvent"
    }

    b, err := json.Marshal(event)
    if err != nil {
        return nil, err
    }

    return &cqrs_es.EventRecord{
        Type: t,
        Data: string(b),
    }, nil
}

func (s *testSerializer) Deserialize(record *cqrs_es.EventRecord) (cqrs_es.Event, error) {
    var e cqrs_es.Event
    switch record.Type {
    case "testEvent":
        e = new(testEvent)
    }

    err := json.Unmarshal([]byte(record.Data), e)
    if err != nil {
        return nil, err
    }

    return e, nil
}