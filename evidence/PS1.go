package evidence

import (
	"regexp"

	"github.com/liserjrqlxue/goUtil/textUtil"
)

func FindPathogenicMissense(fileName, key string, pathogenicRegexp *regexp.Regexp) (map[string]int, map[string]int, map[string]int) {
	var varList = make(map[string]int)
	var pHGVSList = make(map[string]int)
	var pPosList = make(map[string]int)
	itemArray, _ := textUtil.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !pathogenicRegexp.MatchString(item[key]) {
			continue
		}
		if item["Function"] != "missense" {
			continue
		}
		var key = item["Transcript"] + ":" + item["pHGVS"]
		pHGVSList[key]++
		varList[item["MutationName"]]++
		AAPos := getAAPos.FindString(item["pHGVS"])
		if AAPos != "" {
			pPosList[item["Transcript"]+":"+AAPos]++
		}
	}
	return varList, pHGVSList, pPosList
}

// PS1
func CheckPS1(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}

	var mutation = item["MutationName"]
	var key = item["Transcript"] + ":" + item["pHGVS"]
	var countHGVS = hgvsCount[mutation]
	var countPHGVS = hgvsCount[key]
	if countPHGVS > countHGVS {
		return "1"
	} else {
		return "0"
	}
}

func ComparePS1(item map[string]string, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist map[string]int) {
	rule := "PS1"
	val := CheckPS1(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Transcript", "pHGVS")
	}
}
