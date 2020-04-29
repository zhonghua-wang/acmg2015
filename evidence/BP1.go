package evidence

import (
	"github.com/liserjrqlxue/goUtil/textUtil"
)

var (
	BP1geneList map[string]bool
)

func GetBP1geneList(fileName string) {
	var genes = textUtil.File2Array(fileName)
	BP1geneList = make(map[string]bool)
	for _, gene := range genes {
		BP1geneList[gene] = true
	}
}

// BP1
func CheckBP1(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}
	if BP1geneList[item["Gene Symbol"]] {
		return "1"
	} else {
		return "0"
	}
}

func CompareBP1(item map[string]string, ClinVarBP1GeneList, HgmdBP1GeneList map[string]float64) {
	rule := "BP1"
	val := CheckBP1(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
