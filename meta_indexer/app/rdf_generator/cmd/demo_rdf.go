package cmd

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
	"time"
)

// DemoRDFCommand command
type DemoRDF struct {
	cli.Helper
	Count int `cli:"c,count" usage:"count of rdf to generate" dft:"10000"`
}

var DemoRDFCommand = &cli.Command{
	Name: "rdf-examples",
	Desc: "generate N RDF mutations and output to file",
	Argv: func() interface{} { return new(DemoRDF) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*DemoRDF)
		return generateRDFfile(argv.Count)
	},
}

func generateRDFfile(count int) error {
	file, err := os.Create(fmt.Sprintf("kv_%s.rdf", time.Now().Format("2006_01_02T15_04_05")))
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < count; i++ {
		rdf := fmt.Sprintf("_:key%d <wo-key> \"6B3aa543AkotypMaLCeuWDTXFLuG9UKyZCSdJBPStJzEe\" .", i)
		_, err := file.Write([]byte(rdf))
		if err != nil {
			return err
		}
		file.WriteString("\n")
	}
	return nil
}
