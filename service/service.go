package service

import (
    "fmt"
    "sync"
    "github.com/tobyjsullivan/cqrs-es"
    "log"
    "os"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[service] ", 0)
}

type Service struct {
    mx sync.RWMutex
    store cqrs_es.Store
}

func NewService(store cqrs_es.Store) *Service {
    return &Service{
        store: store,
    }
}

func (svc *Service) Execute(id cqrs_es.EntityId, cmd cqrs_es.Command) error {
    svc.mx.Lock()
    defer svc.mx.Unlock()

    newEvents, err := cmd.Execute(svc.store.Events(id))
    logger.Println(fmt.Sprintf("Result of command: (%v, %s)", newEvents, err))

    if err != nil {
        return err
    }

    svc.store.Commit(id, newEvents)
    logger.Println(fmt.Sprintf("Appended events to %s: %v", id, newEvents))

    return nil
}

func (svc *Service) Events(id cqrs_es.EntityId, asOf uint) []cqrs_es.Event {
    svc.mx.RLock()
    defer svc.mx.RUnlock()

    hist := svc.store.Events(id)
    logger.Println(fmt.Sprintf("History for %s: %v", id, hist))
    if uint(len(hist)) <= asOf {
        return []cqrs_es.Event{}
    }

    return hist[asOf:]
}
