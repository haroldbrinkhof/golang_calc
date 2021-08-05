package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"catsandcoding.be/calc/parser"
)

func main() {
	showCalculations := false

	// cli flag parsing
	flag.BoolVar(&showCalculations, "s", false, "show initial term before result, e.g. 1 + 2 = 3")
	flag.Parse()
	// handle any extra terms given
	calculateAndPrint(flag.Args(), showCalculations)

	// handle terms passed via pipe
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err == nil {

		lines := strings.Split(string(bytes), "\n")
		calculateAndPrint(lines, showCalculations)

	} else {
		log.Println(err)
	}
}

func calculateAndPrint(lines []string, showCalculations bool) {
	for i := 0; i < len(lines); i++ {
		if len(strings.TrimSpace(lines[i])) > 0 {
			total, err := parser.Calculate(lines[i])
			var outcome string
			if err == nil {
				outcome = fmt.Sprintf("%f", total)
			} else {
				outcome = err.Error()
			}
			if showCalculations {
				fmt.Printf("%s = %s\n", lines[i], outcome)
			} else {
				fmt.Println(outcome)
			}
		}
	}
}
