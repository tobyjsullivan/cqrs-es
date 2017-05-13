package bank_accounts

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "fmt"
)

func TestOpenBankAccountCommand_Execute_Valid(t *testing.T) {
    history := []cqrs_es.Event{}

    cmd := &OpenBankAccountCommand{ ClientName: "Cool Customer" }

    events, err := cmd.Execute(history)

    if err != nil {
        t.Error("Encountered an unexpected error: "+err.Error())
    }

    if l := len(events); l != 1 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }

    if event := events[0].(*BankAccountOpenedEvent); event.ClientName != cmd.ClientName {
        t.Error(fmt.Sprintf("Unexpected client name. Expected: %s; Actual: %s", cmd.ClientName, event.ClientName))
    }
}

func TestOpenBankAccountCommand_Execute_AlreadyOpen(t *testing.T) {
    history := []cqrs_es.Event{
        &BankAccountOpenedEvent{},
    }

    cmd := &OpenBankAccountCommand{ ClientName: "Cool Customer" }

    events, err := cmd.Execute(history)

    if err := err.(*AccountAlreadyOpenError); err == nil {
        t.Error("Expected error didn't happen")
    }

    if l := len(events); l != 0 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }
}
