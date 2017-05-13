package bank_accounts

import (
    "log"
    "os"
)

var logger *log.Logger

func init()  {
    logger = log.New(os.Stdout, "[bank_accounts] ", 0)
}
