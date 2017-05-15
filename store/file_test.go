package store

import (
    "testing"
    "os"
    "path"
    "fmt"
    "github.com/tobyjsullivan/cqrs-es"
    "github.com/satori/go.uuid"
)

var testDir string

func init() {
    tmp := os.TempDir()
    testDir = path.Join(tmp, "test-out")
    err := os.MkdirAll(testDir, 0700)
    if err != nil {
        panic(err)
    }
}

func initTestDir(name string) string {
    dir := path.Join(testDir, name)
    err := os.MkdirAll(dir, 0700)
    if err != nil {
        panic(err)
    }

    return dir
}

func cleanTestDir(path string) {
    os.RemoveAll(path)
}

func TestFileStore(t *testing.T) {
    dir := initTestDir("t1")
    defer cleanTestDir(dir)

    logger.Println("Test dir is", dir)

    s := NewFileStore(dir, &testSerializer{})

    entity := cqrs_es.EntityId(uuid.NewV4().String())

    hist, err := s.Events(entity)
    if err != nil {
        t.Error(fmt.Sprintf("Unexpected error reading history: %s", err.Error()))
    }

    if l := len(hist); l != 0 {
        t.Error(fmt.Sprintf("Unexpected history length: %d", l))
    }

    content1 := "Event 1 content"
    content2 := "Event 2 content"
    s.Commit(entity, []cqrs_es.Event{
        &testEvent{Content: content1},
        &testEvent{Content: content2},
    })

    hist, err = s.Events(entity)
    if err != nil {
        t.Error(fmt.Sprintf("Unexpected error reading history: %s", err.Error()))
    }

    if l := len(hist); l != 2 {
        t.Error(fmt.Sprintf("Unexpected history length after commit: %d", l))
    }

    for _, e := range hist {
        logger.Println("Found an event: ", e)
    }

    if testEvent := hist[0].(*testEvent); testEvent.Content != content1 {
        t.Error(fmt.Sprintf("Unexpected content in first event. Expected: %s; Actual: %s", content1, testEvent.Content))
    }

    if testEvent := hist[1].(*testEvent); testEvent.Content != content2 {
        t.Error(fmt.Sprintf("Unexpected content in second event. Expected: %s; Actual: %s", content2, testEvent.Content))
    }
}
