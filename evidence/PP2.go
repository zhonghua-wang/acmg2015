package evidence

import "github.com/liserjrqlxue/goUtil/textUtil"

var (
	PP2geneList map[string]bool
)

func GetPP2geneList(fileName string) {
	var genes = textUtil.File2Array(fileName)
	PP2geneList = make(map[string]bool)
	for _, gene := range genes {
		PP2geneList[gene] = true
	}
}

// PP2
func CheckPP2(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}
	if PP2geneList[item["Gene Symbol"]] {
		return "1"
	} else {
		return "0"
	}
}

func ComparePP2(item map[string]string) {
	rule := "PP2"
	val := CheckPP2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
