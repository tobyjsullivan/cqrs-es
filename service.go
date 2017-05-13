package cqrs_es

import (
    "fmt"
    "sync"
)

type Service struct {
    mx sync.RWMutex
    store map[EntityId][]Event
}

func NewService() *Service {
    return &Service{
        store: make(map[EntityId][]Event),
    }
}

func (svc *Service) Execute(id EntityId, cmd Command) error {
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

func (svc *Service) Events(id EntityId, asOf uint) []Event {
    svc.mx.RLock()
    defer svc.mx.RUnlock()

    hist := svc.store[id]
    logger.Println(fmt.Sprintf("History for %s: %v", id, hist))
    if uint(len(hist)) <= asOf {
        return []Event{}
    }

    return hist[asOf:]
}
