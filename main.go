// main.go for testing

package main

import (
	"ena/dirgod/builder"
)

func main() {
	nativeBuilder := &builder.NativeBuilder{}
	nativeBuilder.Operation = "CreateFolder"
	nativeBuilder.Argumentes = map[string]interface{}{
		"path": "test",
		"workingMode": "Force",
		"cache": "true",
	}
	nativeBuilder.CreateNewOperation()
}
