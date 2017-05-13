package bank_accounts

import (
    "github.com/tobyjsullivan/cqrs-es"
)

type WithdrawAmountCommand struct {
    Amount int
}

type AmountWithdrawnEvent struct {
    Amount int
}

func (cmd *WithdrawAmountCommand) Execute(history []cqrs_es.Event) ([]cqrs_es.Event, error) {
    logger.Println("Executing DepositAmount")

    if !accountIsOpen(history) {
        return []cqrs_es.Event{}, &AccountNotOpenError{ msg: "Account has not been opened" }
    }

    if currentBalance(history) < cmd.Amount {
        return []cqrs_es.Event{}, &InsufficientFundsError{ msg: "Insufficient funds"}
    }

    return []cqrs_es.Event{
        &AmountWithdrawnEvent{
            Amount: cmd.Amount,
        },
    }, nil
}
