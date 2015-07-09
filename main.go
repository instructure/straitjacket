package main

import (
	"flag"
	"fmt"
	"straitjacket/engine"
)

func main() {
	var flagRunLanguageTests = flag.Bool("test", false, "Run the language sanity tests")
	var testLanguage = flag.String("lang", "", "Run the tests for one specific language")
	flag.Parse()

	if *flagRunLanguageTests {
		runLanguageTests(*testLanguage)
	} else {
		startServer(":8081")
	}
}

func runLanguageTests(langToRun string) {
	theEngine, err := engine.LoadConfig("config")
	if err != nil {
		panic(err)
	}

	for _, lang := range theEngine.Languages {
		if langToRun == "" || langToRun == lang.Name {
			fmt.Printf("Testing %s\n", lang.VisibleName)
			err := lang.RunTests()
			if err != nil {
				panic(err)
			}
		}
	}

	if langToRun == "" {
		fmt.Printf("All languages A-OK\n")
	}
}
