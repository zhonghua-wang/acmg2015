package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
)

var (
	BP1LoFRatioThreshold = 0.80
)

func CalGeneLoFRatio(fileName, key string, filter *regexp.Regexp, threshold int) map[string]float64 {
	var allCount = make(map[string]int)
	var targetCount = make(map[string]int)
	var ratioCount = make(map[string]float64)
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !filter.MatchString(item[key]) {
			continue
		}
		gene := item["Gene Symbol"]
		allCount[gene]++
		if FuncInfo[item["Function"]] == 3 {
			targetCount[gene]++
		}
	}
	for key, val := range targetCount {
		if threshold > 0 {
			if allCount[key] >= threshold {
				ratioCount[key] = float64(val) / float64(allCount[key])
			}
		} else {
			ratioCount[key] = float64(val) / float64(allCount[key])
		}
	}
	return ratioCount
}

// BP1
func CheckBP1(item map[string]string, ClinVarBP1GeneList, HgmdBP1GeneList map[string]float64) string {
	if item["Function"] == "missense" {
		return "0"
	}
	gene := item["Gene Symbol"]
	if ClinVarBP1GeneList[gene] > BP1LoFRatioThreshold || HgmdBP1GeneList[gene] > BP1LoFRatioThreshold {
		return "1"
	}

	return ""
}

func CompareBP1(item map[string]string, ClinVarBP1GeneList, HgmdBP1GeneList map[string]float64) {
	rule := "BP1"
	val := CheckBP1(item, ClinVarBP1GeneList, HgmdBP1GeneList)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
