package main

import (
	"flag"
	"fmt"
	"straitjacket/engine"
)

func main() {
	var testLanguages = flag.Bool("test", false, "Run the language sanity tests")
	flag.Parse()

	if *testLanguages {
		runLanguageTests()
	} else {
		startServer(":8081")
	}
}

func runLanguageTests() {
	theEngine, err := engine.LoadConfig("config")
	if err != nil {
		panic(err)
	}

	for _, lang := range theEngine.Languages {
		fmt.Printf("Testing %s\n", lang.VisibleName)
		err := lang.RunTests()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("All languages A-OK\n")
}
