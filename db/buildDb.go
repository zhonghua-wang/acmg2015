package main

import (
	"flag"
	simple_util "github.com/liserjrqlxue/simple-util"
)

var (
	PVS1db = flag.String(
		"pvs1",
		"",
		"PVS1 LOF不耐受基因集",
	)
	PVS1sheet = flag.String(
		"pvs1sheet",
		"！原有总基因+新增",
		"PVS1 sheet name",
	)
)

func main() {
	flag.Parse()
	if *PVS1db != "" {
		PVS1GeneList := xlsx2mapInt(*PVS1db, *PVS1sheet)
		simple_util.CheckErr(simple_util.Json2File("PVS1GeneList.json", PVS1GeneList))
	}
}
