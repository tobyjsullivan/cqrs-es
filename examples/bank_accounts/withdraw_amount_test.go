package bank_accounts

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "fmt"
)

func TestWithdrawAmountCommand_Execute_Valid(t *testing.T) {
    hist := []cqrs_es.Event{
        &BankAccountOpenedEvent{},
        &AmountDepositedEvent{Amount: 354},
    }

    cmd := &WithdrawAmountCommand{Amount: 95}

    events, err := cmd.Execute(hist)

    if err != nil {
        t.Error("Encountered an unexpected error: "+err.Error())
    }

    if l := len(events); l != 1 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }

    if event := events[0].(*AmountWithdrawnEvent); event.Amount != cmd.Amount {
        t.Error(fmt.Sprintf("Unexpected deposit amount. Expected: %d; Actual: %d", cmd.Amount, event.Amount))
    }
}

func TestWithdrawAmountCommand_Execute_NotOpened(t *testing.T) {
    hist := []cqrs_es.Event{}

    cmd := &WithdrawAmountCommand{Amount: 123}

    events, err := cmd.Execute(hist)

    if err := err.(*AccountNotOpenError); err == nil {
        t.Error("Expected error didn't happen")
    }

    if l := len(events); l != 0 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }
}

func TestWithdrawAmountCommand_Execute_InsufficientFunds(t *testing.T) {
    hist := []cqrs_es.Event{
        &BankAccountOpenedEvent{},
        &AmountDepositedEvent{Amount: 164},
    }

    cmd := &WithdrawAmountCommand{Amount: 276}

    events, err := cmd.Execute(hist)

    if err := err.(*InsufficientFundsError); err == nil {
        t.Error("Expected error didn't happen")
    }

    if l := len(events); l != 0 {
        t.Error(fmt.Sprintf("Unexpected events length: %d", l))
    }
}