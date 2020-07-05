package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	routerFileName = flag.String("i", "router.go", "router config file name")
	outputFileName = flag.String("o", "router_gen.go", "generated router file name")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = Usage
	flag.Parse()

	cfg, err := Analyze(*routerFileName)
	if err != nil {
		err = fmt.Errorf("failed to analyze router file: %w", err)
		log.Fatalln(err)
	}
	g := &Generator{}
	if err := g.Generate(cfg); err != nil {
		err = fmt.Errorf("failed to generate Go source: %w", err)
		log.Fatalln(err)
	}

	f, err := os.Create(*outputFileName)
	defer f.Close()
	if err != nil {
		err = fmt.Errorf("failed to create file: %w", err)
		log.Fatalln(err)
	}
	_, err = f.Write(g.format())
	if err != nil {
		err = fmt.Errorf("failed to write bytes to the file: %w", err)
		log.Fatalln(err)
	}

	log.Printf("Router file generated to %s\n", *outputFileName)
}
