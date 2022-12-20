// main.go for testing

package main

import (
	"ena/dirgod/builder"
	"fmt"
)

func main() {
	nativeBuilder := &builder.NativeBuilder{}
	nativeBuilder.SetBasePath("test")
	nativeBuilder.SetWorkingMode("force")
	nativeBuilder.SetCacheMode("true")
	arguments := map[string]any{
		"path":        "asd",
		"workingMode": "strict",
		"cache":       false,
	}

	err := nativeBuilder.CreateNewOperation("CreateFolder", arguments)
	if err != nil {
		panic(err)
	}

	err = nativeBuilder.CreateNewOperation("CreateFile")
	if err != nil {
		panic(err)
	}

	operation := nativeBuilder.GetOperations()

	fmt.Println(operation)
}
