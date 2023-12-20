// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"

	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
)

const (
	fileServerTestDir = "vault-agent-file-test"
)

func testFileSink(t *testing.T, log hclog.Logger) (*sink.SinkConfig, string) {
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s.", fileServerTestDir))
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(tmpDir, "token")

	config := &sink.SinkConfig{
		Logger: log.Named("sink.file"),
		Config: map[string]interface{}{
			"path": path,
		},
	}

	s, err := NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = s

	return config, tmpDir
}

func TestFileSink(t *testing.T) {
	log := corehelpers.Logger

	fs, tmpDir := testFileSink(t, log)
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "token")

	uuidStr, _ := uuid.GenerateUUID()
	if err := fs.WriteToken(uuidStr); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != os.FileMode(0o640) {
		t.Fatalf("wrong file mode was detected at %s", path)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileBytes) != uuidStr {
		t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
	}
}

func testFileSinkMode(t *testing.T, log hclog.Logger) (*sink.SinkConfig, string) {
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s.", fileServerTestDir))
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(tmpDir, "token")

	config := &sink.SinkConfig{
		Logger: log.Named("sink.file"),
		Config: map[string]interface{}{
			"path": path,
			"mode": 0o644,
		},
	}

	s, err := NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = s

	return config, tmpDir
}

func TestFileSinkMode(t *testing.T) {
	log := corehelpers.Logger

	fs, tmpDir := testFileSinkMode(t, log)
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "token")

	uuidStr, _ := uuid.GenerateUUID()
	if err := fs.WriteToken(uuidStr); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != os.FileMode(0o644) {
		t.Fatalf("wrong file mode was detected at %s", path)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileBytes) != uuidStr {
		t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
	}
}
