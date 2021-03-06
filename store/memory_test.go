package store

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "github.com/satori/go.uuid"
    "fmt"
)

func TestMemoryStore(t *testing.T) {
    s := NewMemoryStore()

    entity := cqrs_es.EntityId(uuid.NewV4().String())

    hist, err := s.Events(entity)
    if err != nil {
        t.Error(fmt.Sprintf("Unexpected error reading history: %s", err.Error()))
    }

    if l := len(hist); l != 0 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }

    content1 := "Event 1 content"
    content2 := "Event 2 content"
    s.Commit(entity, []cqrs_es.Event{
        &testEvent{Content: content1},
        &testEvent{Content: content2},
    })

    hist, err = s.Events(entity)
    if err != nil {
        t.Error(fmt.Sprintf("Unexpected error reading history: %s", err.Error()))
    }

    if l := len(hist); l != 2 {
        t.Error(fmt.Sprintf("Unexpected history length after commit: %d", l))
    }

    if testEvent := hist[0].(*testEvent); testEvent.Content != content1 {
        t.Error(fmt.Sprintf("Unexpected content in first event. Expected: %s; Actual: %s", content1, testEvent.Content))
    }

    if testEvent := hist[1].(*testEvent); testEvent.Content != content2 {
        t.Error(fmt.Sprintf("Unexpected content in second event. Expected: %s; Actual: %s", content2, testEvent.Content))
    }
}
