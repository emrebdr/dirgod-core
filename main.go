// main.go for testing

package main

import (
	"fmt"

	"github.com/emrebdr/dirgod-core/builder"
)

func main() {
	nativeBuilder := builder.NewNativeBuilder()
	nativeBuilder.SetBasePath("test")
	nativeBuilder.SetWorkingMode("force")
	nativeBuilder.SetCacheMode("true")
	arguments := map[string]any{
		"source":      "asd",
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
	fmt.Printf("operation: %v\n", operation)
}
