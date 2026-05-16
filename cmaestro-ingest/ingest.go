package cmaestro_ingest

import "fmt"

func Ingest(data string) []Function {
	fmt.Println("Ingesting", data)
	return []Function{
		NewFunction("simple_entrypoint", "string"),
	}
}
