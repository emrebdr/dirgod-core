// main.go for testing

package main

import (
	"ena/dirgod/builder"
)

func main() {
	nativeBuilder := &builder.NativeBuilder{}
	arguments := map[string]string{
		"path":        "test",
		"workingMode": "Force",
		"cache":       "true",
	}
	nativeBuilder.CreateNewOperation("CreateFolder", arguments)
}
