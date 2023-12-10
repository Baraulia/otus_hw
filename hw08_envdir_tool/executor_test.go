package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestData() error {
	err := os.Setenv("HELLO", "SHOULD_REPLACE")
	if err != nil {
		return fmt.Errorf("error while setting env variable %s: %w", "HELLO", err)
	}

	err = os.Setenv("FOO", "SHOULD_REPLACE")
	if err != nil {
		return fmt.Errorf("error while setting env variable %s: %w", "FOO", err)
	}

	err = os.Setenv("UNSET", "SHOULD_REMOVE")
	if err != nil {
		return fmt.Errorf("error while setting env variable %s: %w", "UNSET", err)
	}

	err = os.Setenv("ADDED", "from original env")
	if err != nil {
		return fmt.Errorf("error while setting env variable %s: %w", "ADDED", err)
	}

	err = os.Setenv("EMPTY", "SHOULD_BE_EMPTY")
	if err != nil {
		return fmt.Errorf("error while setting env variable %s: %w", "EMPTY", err)
	}

	return nil
}

func TestRunCmd(t *testing.T) {
	if err := setupTestData(); err != nil {
		t.Fatalf("error while setting up test data: %v", err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error while getting current directory: %v", err)
	}

	arg1 := "arg1=1"
	arg2 := "arg2=2"

	testTable := []struct {
		name         string
		commandPath  string
		env          Environment
		expectedCode int
		expectedEnv  map[string]string
	}{
		{
			name:        "simple case",
			commandPath: fmt.Sprintf("%s/testdata/echo.sh", currentDir),
			env: Environment{
				"HELLO": {Value: "hello", NeedRemove: true},
				"FOO":   {Value: "foo", NeedRemove: true},
				"UNSET": {NeedRemove: true},
				"EMPTY": {Value: " ", NeedRemove: true},
			},
			expectedCode: 0,
			expectedEnv: map[string]string{
				"HELLO": "hello",
				"FOO":   "foo",
				"UNSET": "",
				"EMPTY": " ",
				"ADDED": "from original env",
			},
		},
		{
			name:         "directory not exist",
			commandPath:  fmt.Sprintf("%s/fail_testdata/echo.sh", currentDir),
			env:          Environment{},
			expectedCode: 127,
			expectedEnv:  map[string]string{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			code := RunCmd([]string{"/bin/bash", testCase.commandPath, arg1, arg2}, testCase.env)
			require.Equal(t, code, testCase.expectedCode)
			for k, v := range testCase.expectedEnv {
				fmt.Printf("check %s\n", k)
				require.Equal(t, os.Getenv(k), v)
			}
		})
	}
}
