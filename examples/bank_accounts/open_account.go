package bank_accounts

import (
    "github.com/tobyjsullivan/cqrs-es"
)

type OpenBankAccountCommand struct {
    ClientName string
}

type BankAccountOpenedEvent struct {
    ClientName string
}

func (cmd *OpenBankAccountCommand) Execute(history []cqrs_es.Event) ([]cqrs_es.Event, error) {
    logger.Println("Executing OpenBankAccount")

    if accountIsOpen(history) {
        return []cqrs_es.Event{}, &AccountAlreadyOpenError{ msg: "Account already exists" }
    }

    return []cqrs_es.Event{
        &BankAccountOpenedEvent{
            ClientName: cmd.ClientName,
        },
    }, nil
}
