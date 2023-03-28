// main.go for testing

package main

import (
	"github.com/emrebdr/dirgod-core/builder"
)

func main() {
	nativeBuilder := builder.NewNativeBuilder()
	nativeBuilder.SetBasePath("test")
	nativeBuilder.SetWorkingMode("force")
	nativeBuilder.SetCacheMode("true")
	arguments := map[string]any{
		"source":      "models",
		"destination": "tests",
		"workingMode": "strict",
		"cache":       false,
	}

	err := nativeBuilder.CreateNewOperation("copyfolder", arguments)
	if err != nil {
		panic(err)
	}

	operation := nativeBuilder.GetOperations()
	operation[0].Exec()
}
