package service

import (
    "testing"
    "github.com/satori/go.uuid"
    "fmt"
    "errors"
    "github.com/tobyjsullivan/cqrs-es/store"
    "github.com/tobyjsullivan/cqrs-es"
)

type testCommand struct {
    succeed bool
    content string
}

type testEvent struct {
    content string
}

func (cmd *testCommand) Execute(history []cqrs_es.Event) ([]cqrs_es.Event, error) {
    if !cmd.succeed {
        return []cqrs_es.Event{}, errors.New("Command failed")
    }

    return []cqrs_es.Event{
        &testEvent{
            content: cmd.content,
        },
    }, nil
}

func TestNewService(t *testing.T) {
    svc := NewService(store.NewMemoryStore())

    history := svc.Events(cqrs_es.EntityId(uuid.NewV4().String()), 0)

    if l := len(history); l != 0 {
        t.Error(fmt.Sprintf("History of new service had unexpected length (%d)", l))
    }
}

func TestService_Execute_Success(t *testing.T) {
    svc := NewService(store.NewMemoryStore())

    entityId := cqrs_es.EntityId(uuid.NewV4().String())

    cmd := &testCommand{
        succeed: true,
        content: "Ipsum lorem",
    }

    if err := svc.Execute(entityId, cmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    newHistory := svc.Events(entityId, 0)

    if l := len(newHistory); l != 1 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }

    if event := newHistory[0].(*testEvent); event.content != cmd.content {
        t.Error(fmt.Sprintf("Unexpected event content. Expected: %s; Actual: %s", cmd.content, event.content))
    }
}

func TestService_Execute_Failure(t *testing.T) {
    svc := NewService(store.NewMemoryStore())

    entityId := cqrs_es.EntityId(uuid.NewV4().String())

    cmd := &testCommand{
        succeed: false,
    }

    if err := svc.Execute(entityId, cmd); err == nil {
        t.Error("Did not receive expected error: ")
    }

    newHistory := svc.Events(entityId, 0)

    if l := len(newHistory); l != 0 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }
}