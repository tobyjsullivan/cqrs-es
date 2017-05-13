package bank_accounts

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "fmt"
)

func TestDepositAmountCommand_Execute_Valid(t *testing.T) {
    hist := []cqrs_es.Event{
        &BankAccountOpenedEvent{},
    }

    cmd := &DepositAmountCommand{Amount: 365}

    events, err := cmd.Execute(hist)

    if err != nil {
        t.Error("Encountered an unexpected error: "+err.Error())
    }

    if l := len(events); l != 1 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }

    if event := events[0].(*AmountDepositedEvent); event.Amount != cmd.Amount {
        t.Error(fmt.Sprintf("Unexpected deposit amount. Expected: %d; Actual: %d", cmd.Amount, event.Amount))
    }
}

func TestDepositAmountCommand_Execute_NotOpened(t *testing.T) {
    hist := []cqrs_es.Event{}

    cmd := &DepositAmountCommand{Amount: 123}

    events, err := cmd.Execute(hist)

    if err := err.(*AccountNotOpenError); err == nil {
        t.Error("Expected error didn't happen")
    }

    if l := len(events); l != 0 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }
}
