package bank_accounts

type AccountAlreadyOpenError struct {
    msg string
}
func (e *AccountAlreadyOpenError) Error() string { return e.msg }

type AccountNotOpenError struct {
    msg string
}
func (e *AccountNotOpenError) Error() string { return e.msg }

type InsufficientFundsError struct {
    msg string
}
func (e *InsufficientFundsError) Error() string { return e.msg }
