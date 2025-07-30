package assert

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func writeSnapshot(t *testing.T, path, actual string) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create snapshot file: %s", err)
	}
	_, err = f.Write([]byte(actual))
	if err != nil {
		t.Fatalf("Failed to write to snapshot file: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Failed to close snapshot file: %s", err)
	}
}

func TextSnapshot(t *testing.T, path, actual string) {
	if !pathExists(path) {
		writeSnapshot(t, path, actual)
		t.Errorf("Snapshot file created: %s", path)
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read snapshot file: %s", err)
	}
	expected := string(content)

	if expected != actual {
		t.Errorf("Value doesn't match snapshot %s:\n%s", path, actual)
	}

	if os.Getenv("UPDATE_SNAPSHOTS") == "" {
		return
	}

	writeSnapshot(t, path, actual)
	fmt.Printf("Snapshot file %s updated\n", path)
}
