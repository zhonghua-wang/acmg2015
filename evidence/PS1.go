package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
)

// regexp
var (
	IsClinVarPLP = regexp.MustCompile(`Pathogenic|Likely_pathogenic`)
	IsHgmdDM     = regexp.MustCompile(`DM$|DM\|`)
)

func FindPathogenicMissense(fileName, key string, pathogenicRegexp *regexp.Regexp) (map[string]int, map[string]int) {
	var varList = make(map[string]int)
	var pHGVSlist = make(map[string]int)
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !pathogenicRegexp.MatchString(item[key]) {
			continue
		}
		if item["Function"] != "missense" {
			continue
		}
		var key = item["Transcript"] + ":" + item["pHGVS"]
		pHGVSlist[key]++
		varList[item["MutationName"]]++
	}
	return varList, pHGVSlist
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
