package main

import (
	"flag"
	"os"
	"regexp"

	"github.com/liserjrqlxue/acmg2015/evidence"

	"github.com/liserjrqlxue/simple-util"
)

var (
	clinvar = flag.String(
		"clinvar",
		"",
		"clinvar anno file",
	)
)

var (
	isPathogenic       = regexp.MustCompile(`Pathogenic`)
	islikelyPathogenic = regexp.MustCompile(`Likely_pathogenic`)
)

var BP1geneSet = make(map[string]int)

type BP1info struct {
	Count int
	P     int
	LOF   int
	Ratio float32
}

func main() {
	flag.Parse()
	if *clinvar == "" {
		flag.Parse()
		os.Exit(1)
	}
	clinvarArray, _ := simple_util.File2MapArray(*clinvar, "\t", nil)
	var clinvarDb = make(map[string]BP1info)
	for _, item := range clinvarArray {
		//mainKey:=strings.Join([]string{item["#Chr"],item["Start"],item["Stop"],item["Ref"],item["Call"]},"-")
		clinvarTag := item["ClinVar Significance"]
		geneSymbol := item["Gene Symbol"]
		info, ok := clinvarDb[geneSymbol]
		if !ok {
			info = BP1info{0, 0, 0, 0}
		}
		info.Count++
		if isPathogenic.MatchString(clinvarTag) || islikelyPathogenic.MatchString(clinvarTag) {
			info.P++
			function := item["Function"]
			if evidence.FuncInfo[function] == 3 {
				info.LOF++
			}
		}
		clinvarDb[geneSymbol] = info
	}
	for key, value := range clinvarDb {
		if value.P == 0 {
			value.Ratio = 0
		} else {
			value.Ratio = float32(value.LOF) / float32(value.P)
		}
		if value.P >= 10 && value.Ratio > 0.8 {
			BP1geneSet[key]++
		}
		clinvarDb[key] = value
	}
	err := simple_util.Json2File("BP1.GeneSet.json", BP1geneSet)
	simple_util.CheckErr(err)
	err = simple_util.Json2File("BP1.ClinVar.json", clinvarDb)
	simple_util.CheckErr(err)
}
