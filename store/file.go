package store

import (
    "github.com/tobyjsullivan/cqrs-es"
    "sync"
    "fmt"
    "path"
    "os"
    "bufio"
    "encoding/json"
    "bytes"
    "io"
)

type fileStore struct {
    storesDir string
    serializer cqrs_es.EventSerializer
    mx        sync.RWMutex
}

func NewFileStore(dir string, serializer cqrs_es.EventSerializer) cqrs_es.Store {
    return &fileStore{
        serializer: serializer,
        storesDir: dir,
    }
}


func (s *fileStore) Events(id cqrs_es.EntityId) ([]cqrs_es.Event, error) {
    s.mx.RLock()
    defer s.mx.RUnlock()

    logger.Println(fmt.Sprintf("Fetching events for %s", id))

    filePath := path.Join(s.storesDir, string(id))

    if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
        return []cqrs_es.Event{}, nil
    }

    fHandle, err := os.Open(filePath)
    defer fHandle.Close()
    if err != nil {
        return nil, err
    }

    var events []cqrs_es.Event
    scanner := bufio.NewScanner(fHandle)
    for scanner.Scan() {
        var record cqrs_es.EventRecord
        content := scanner.Bytes()
        logger.Println("Store data:", string(content))

        if err := json.Unmarshal(content, &record); err != nil {
            logger.Println("Error while unmarshalling serializable: "+err.Error())
            return nil, err
        }

        e, err := s.serializer.Deserialize(&record)
        if err != nil {
            logger.Println("Error while unmarshalling event: "+err.Error())
            return nil, err
        }

        events = append(events, e)
    }

    return events, nil
}

func (s *fileStore) Commit(id cqrs_es.EntityId, events []cqrs_es.Event) error {
    s.mx.Lock()
    defer s.mx.Unlock()

    logger.Println(fmt.Sprintf("Appending events to %s: %v", id, events))

    // Serialize events into a buffer so that we can drop if there's an error
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    for _, curEvent := range events {
        rec, err := s.serializer.Serialize(curEvent)
        if err != nil {
            return err
        }

        if wrErr := encoder.Encode(rec); wrErr != nil {
            return err
        }
    }

    filePath := path.Join(s.storesDir, string(id))
    fWriter, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
    if err != nil {
        return err
    }
    defer fWriter.Close()

    // Write buffer to file
    if _, err = io.Copy(fWriter, &buf); err != nil {
        return err
    }

    return nil
}
