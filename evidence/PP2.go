package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
)

var (
	PP2MissenseRatioThreshold = 0.50
)

func CalGeneMissenseRatio(fileName, key string, filter *regexp.Regexp, threshold int) map[string]float64 {
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
		if item["Function"] == "missense" {
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

// PP2
func CheckPP2(item map[string]string, ClinVarPP2GeneList, HgmdPP2GeneList map[string]float64) string {
	if item["Function"] != "missense" {
		return "0"
	}
	gene := item["Gene Symbol"]
	if ClinVarPP2GeneList[gene] > PP2MissenseRatioThreshold || HgmdPP2GeneList[gene] > PP2MissenseRatioThreshold {
		return "1"
	}
	return ""
}

func ComparePP2(item map[string]string, ClinVarPP2GeneList, HgmPP2GeneList map[string]float64) {
	rule := "PP2"
	val := CheckPP2(item, ClinVarPP2GeneList, HgmPP2GeneList)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
