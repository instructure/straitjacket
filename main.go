package main

import (
	"flag"
	"fmt"
	"straitjacket/engine"
)

func main() {
	var flagRunLanguageTests = flag.Bool("test", false, "Run the language sanity tests")
	var testLanguage = flag.String("lang", "", "Run the tests for one specific language")
	var disableApparmor = flag.Bool("disable-apparmor", false, "Disable all apparmor usage. NOT FOR PROD obv.")
	flag.Parse()

	if *disableApparmor {
		fmt.Printf("Running without apparmor support!\n")
	}

	engine, err := engine.New("config", *disableApparmor)
	if err != nil {
		panic(err)
	}

	if *flagRunLanguageTests {
		runLanguageTests(engine, *testLanguage)
	} else {
		startServer(engine, ":8081")
	}
}

func runLanguageTests(engine *engine.Engine, langToRun string) {
	for _, lang := range engine.Languages() {
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
