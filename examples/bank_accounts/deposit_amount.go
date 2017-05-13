package bank_accounts

import (
    "github.com/tobyjsullivan/cqrs-es"
)

type DepositAmountCommand struct {
    Amount int
}

type AmountDepositedEvent struct {
    Amount int
}

func (cmd *DepositAmountCommand) Execute(history []cqrs_es.Event) ([]cqrs_es.Event, error) {
    logger.Println("Executing DepositAmount")

    if !accountIsOpen(history) {
        return []cqrs_es.Event{}, &AccountNotOpenError{ msg: "Account has not been opened" }
    }

    return []cqrs_es.Event{
        &AmountDepositedEvent{
            Amount: cmd.Amount,
        },
    }, nil
}