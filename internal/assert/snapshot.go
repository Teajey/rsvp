package assert

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"testing"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

var newline = []byte("\n")

func ensureSuffix(s, suffix []byte) []byte {
	if !bytes.HasSuffix(s, suffix) {
		s = append(s, suffix...)
	}
	return s
}

func writeSnapshot(t *testing.T, path string, actual []byte) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create snapshot file: %s", err)
	}
	_, err = f.Write(actual)
	if err != nil {
		t.Fatalf("Failed to write to snapshot file: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Failed to close snapshot file: %s", err)
	}
}

func snapshot(t *testing.T, path string, actual []byte) {
	if !pathExists(path) {
		writeSnapshot(t, path, actual)
		t.Errorf("Snapshot file created: %s", path)
		return
	}

	actual = ensureSuffix(actual, newline)

	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read snapshot file: %s", err)
	}

	expected = ensureSuffix(expected, newline)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Value doesn't match snapshot %s:\n\nExpected:\n%s\n\nActual:\n%s\n", path, string(expected), string(actual))
	}

	if os.Getenv("UPDATE_SNAPSHOTS") == "" {
		return
	}

	writeSnapshot(t, path, actual)
	t.Logf("Snapshot file %s updated\n", path)
}

func SnapshotText(t *testing.T, actual string) {
	snapshot(t, fmt.Sprintf("%s.snap.txt", t.Name()), []byte(actual))
}

func SnapshotJson(t *testing.T, actual any) {
	data, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf("Snapshot failed to serialize actual to JSON: %s", err)
	}
	snapshot(t, fmt.Sprintf("%s.snap.json", t.Name()), data)
}

func SnapshotXml(t *testing.T, actual any) {
	data, err := xml.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf("Snapshot failed to serialize actual to XML: %s", err)
	}
	snapshot(t, fmt.Sprintf("%s.snap.xml", t.Name()), data)
}
