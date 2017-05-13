package bank_accounts

import (
    "testing"
    "github.com/tobyjsullivan/cqrs-es"
    "github.com/satori/go.uuid"
    "fmt"
)

func TestIntegration(t *testing.T) {
    svc := cqrs_es.NewService()

    accountId := cqrs_es.EntityId(uuid.NewV4().String())

    // Try making a deposit before account is opened
    if err := svc.Execute(accountId, &DepositAmountCommand{ Amount: 20 }); err == nil {
        t.Error("Did not receive expected error.")
    }

    // Open the account.
    openCmd := &OpenBankAccountCommand{
        ClientName: "Test Account Holder",
    }
    if err := svc.Execute(accountId, openCmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    // Deposit some money.
    depCmd := &DepositAmountCommand{
        Amount: 500,
    }
    if err := svc.Execute(accountId, depCmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    // Try withdrawing too much money
    if err := svc.Execute(accountId, &WithdrawAmountCommand{ Amount: 2000 }); err == nil {
        t.Error("Did not receive expected error.")
    }

    // Withdraw some money.
    wdrCmd := &WithdrawAmountCommand{
        Amount: 350,
    }
    if err := svc.Execute(accountId, wdrCmd); err != nil {
        t.Error("Unexpected error: "+err.Error())
    }

    // Try opening the account again
    if err := svc.Execute(accountId, &OpenBankAccountCommand{ ClientName: "No good" }); err == nil {
        t.Error("Did not receive expected error.")
    }

    // Compute aggregate balance.
    if amount := currentBalance(svc.Events(accountId, 0)); amount != 150 {
        t.Error(fmt.Sprintf("Unexpected current balance: %d", amount))
    }
}
