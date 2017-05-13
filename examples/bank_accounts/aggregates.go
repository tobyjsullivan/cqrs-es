package bank_accounts

import "github.com/tobyjsullivan/cqrs-es"

func accountIsOpen(history []cqrs_es.Event) bool {
    for _, e := range history {
        switch e.(type) {
        case *BankAccountOpenedEvent:
            return true
        }
    }

    return false
}

func currentBalance(history []cqrs_es.Event) int {
    accrued := 0

    for _, e := range history {
        switch e := e.(type) {
        case *AmountDepositedEvent:
            accrued += e.Amount
        case *AmountWithdrawnEvent:
            accrued -= e.Amount
        }
    }

    return accrued
}
