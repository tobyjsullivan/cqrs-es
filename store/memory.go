package store

import (
    "sync"
    "github.com/tobyjsullivan/cqrs-es"
    "log"
    "os"
    "fmt"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[store] ", 0)
}

type memoryStore struct {
    mx sync.RWMutex
    hash map[cqrs_es.EntityId][]cqrs_es.Event
}

func NewMemoryStore() cqrs_es.Store {
    return &memoryStore{
        hash: make(map[cqrs_es.EntityId][]cqrs_es.Event),
    }
}

func (s *memoryStore) Events(id cqrs_es.EntityId) []cqrs_es.Event {
    s.mx.RLock()
    defer s.mx.RUnlock()

    logger.Println(fmt.Sprintf("Fetching events for %s", id))
    return s.hash[id]
}

func (s *memoryStore) Commit(id cqrs_es.EntityId, events []cqrs_es.Event) {
    s.mx.Lock()
    defer s.mx.Unlock()

    logger.Println(fmt.Sprintf("Appending events to %s: %v", id, events))
    s.hash[id] = append(s.hash[id], events...)
}


