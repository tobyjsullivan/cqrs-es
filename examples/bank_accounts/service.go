package bank_accounts

import (
    "github.com/tobyjsullivan/cqrs-es"
    "sync"
    "fmt"
    "log"
    "os"
)

var logger *log.Logger

func init()  {
    logger = log.New(os.Stdout, "[bank_accounts] ", 0)
}

type service struct {
    mx sync.RWMutex
    store map[cqrs_es.EntityId][]cqrs_es.Event
}

func NewService() cqrs_es.Service {
   return &service{
       store: make(map[cqrs_es.EntityId][]cqrs_es.Event),
   }
}

func (svc *service) Execute(id cqrs_es.EntityId, cmd cqrs_es.Command) error {
    svc.mx.Lock()
    defer svc.mx.Unlock()

    newEvents, err := cmd.Execute(svc.store[id])
    logger.Println(fmt.Sprintf("Result of command: (%v, %s)", newEvents, err))

    if err != nil {
        return err
    }

    svc.store[id] = append(svc.store[id], newEvents...)
    logger.Println(fmt.Sprintf("Updated history for %s: %v", id, svc.store[id]))

    return nil
}

func (svc *service) Events(id cqrs_es.EntityId, asOf uint) []cqrs_es.Event {
    svc.mx.RLock()
    defer svc.mx.RUnlock()

    hist := svc.history(id)
    logger.Println(fmt.Sprintf("History for %s: %v", id, hist))
    if uint(len(hist)) <= asOf {
        return []cqrs_es.Event{}
    }


    return hist[asOf:]
}


func (svc *service) history(id cqrs_es.EntityId) []cqrs_es.Event {
    return svc.store[id]
}
