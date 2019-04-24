package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
)

// regexp
var (
	getAAPos = regexp.MustCompile(`^p\.[A-Z]\d+`)
)

func FindPathogenicMissense(fileName, key string, pathogenicRegexp *regexp.Regexp) (map[string]int, map[string]int, map[string]int) {
	var varList = make(map[string]int)
	var pHGVSList = make(map[string]int)
	var pPosList = make(map[string]int)
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
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
func CheckPS1(item map[string]string, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist map[string]int) string {
	if item["Function"] == "missense" {
		return "0"
	}

	var mutation = item["MutationName"]
	var key = item["Transcript"] + ":" + item["pHGVS"]
	var flag1, flag2 bool
	ClinVarMissenseNo := ClinVarMissense[mutation]
	HGMDMissenseNo := HGMDMissense[mutation]
	ClinVarPHGVSNo := ClinVarPHGVSlist[key]
	HGMDPHGVSNo := HGMDPHGVSlist[key]
	if ClinVarMissenseNo < ClinVarPHGVSNo {
		flag1 = true
	}
	if HGMDMissenseNo < HGMDPHGVSNo {
		flag2 = true
	}

	if flag1 && flag2 {
		return "1"
	} else if flag1 || flag2 {
		return "2"
	} else {
		return "0"
	}
}

func ComparePS1(item map[string]string, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist map[string]int) {
	rule := "PS1"
	val := CheckPS1(item, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Transcript", "pHGVS")
	}
}
