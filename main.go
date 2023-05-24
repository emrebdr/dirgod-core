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
		"source":      "test.py",
		"workingMode": "strict",
		"cache":       false,
	}

	err := nativeBuilder.CreateNewOperation("createfile", arguments)
	if err != nil {
		panic(err)
	}

	err = nativeBuilder.Execute([]byte("new commit"))
	fmt.Printf("err: %v\n", err)
}
