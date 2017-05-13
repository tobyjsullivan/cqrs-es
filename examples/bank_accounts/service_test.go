package bank_accounts

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "github.com/satori/go.uuid"
    "fmt"
)

func TestNewService(t *testing.T) {
    svc := NewService()

    history := svc.Events(cqrs_es.EntityId(uuid.NewV4().String()), 0)

    if l := len(history); l != 0 {
        t.Error(fmt.Sprintf("History of new service had unexpected length (%d)", l))
    }
}

func TestService_Execute_Success(t *testing.T) {
    svc := NewService()

    accountId := cqrs_es.EntityId(uuid.NewV4().String())

    cmd := &OpenBankAccountCommand{
        ClientName: "Test Customer",
    }

    if err := svc.Execute(accountId, cmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    newHistory := svc.Events(accountId, 0)

    if l := len(newHistory); l != 1 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }

    if event := newHistory[0].(*BankAccountOpenedEvent); event.ClientName != "Test Customer" {
        t.Error(fmt.Sprintf("Unexpected client name: %s", event.ClientName))
    }
}

func TestService_Execute_Failure(t *testing.T) {
    svc := NewService()

    accountId := cqrs_es.EntityId(uuid.NewV4().String())

    cmd := &OpenBankAccountCommand{
        ClientName: "Test Customer",
    }

    if err := svc.Execute(accountId, cmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    // Execute the OpenAccount command a second time which is invalid
    if err := svc.Execute(accountId, cmd); err == nil {
        t.Error("Expected error did not occur")
    }

    newHistory := svc.Events(accountId, 0)

    if l := len(newHistory); l != 1 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }

    if event := newHistory[0].(*BankAccountOpenedEvent); event.ClientName != "Test Customer" {
        t.Error(fmt.Sprintf("Unexpected client name: %s", event.ClientName))
    }
}