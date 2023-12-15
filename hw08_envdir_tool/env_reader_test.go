package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error while getting current directory: %v", err)
	}

	testTable := []struct {
		name          string
		path          string
		expectedError bool
		expectedEnv   Environment
	}{
		{
			name:          "simple case",
			path:          fmt.Sprintf("%s/testdata/env", currentDir),
			expectedError: false,
			expectedEnv: Environment{
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: true},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: true},
				"UNSET": EnvValue{NeedRemove: true},
				"EMPTY": EnvValue{NeedRemove: true},
			},
		},
		{
			name:          "directory not exist",
			path:          fmt.Sprintf("%s/fail_testdata/env", currentDir),
			expectedError: true,
			expectedEnv:   Environment{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			env, err := ReadDir(testCase.path)
			fmt.Println(env)

			if testCase.expectedError {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
			}

			for k, v := range testCase.expectedEnv {
				fmt.Printf("check %s\n", k)
				require.Equal(t, v.Value, env[k].Value)
				require.Equal(t, v.NeedRemove, env[k].NeedRemove)
			}
		})
	}
}
