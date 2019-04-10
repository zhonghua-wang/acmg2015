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

// colname
var (
	clinvarCol = "ClinVar Significance"
	hgmdCol    = "HGMD Pred"
)

// PS1
// to be update
var missenseHit = map[string]map[string]bool{
	"1-957604-957605-G-A": {
		clinvarCol: true,
		hgmdCol:    true,
	},
}

// PS1
func CheckPS1(item map[string]string) int {
	var function = item["Function"]
	if function == "missense" {
		return 0
	}
	var flag1, flag2 bool
	if IsClinVarPLP.MatchString(item[clinvarCol]) {
		flag1 = true
	}
	if IsHgmdDM.MatchString(item[hgmdCol]) {
		flag2 = true
	}

	var (
		chr   = item["#Chr"]
		start = item["Start"]
		stop  = item["Stop"]
		ref   = item["Ref"]
		alt   = item["Call"]
		//transcript=inputData["Transcript"]
		//aaref=inputData["AAref"]
		//aaalt=inputData["AAalt"]
	)
	var key1 = chr + "-" + start + "-" + stop + "-" + ref + "-" + alt
	if missenseHit[key1][clinvarCol] {
		flag1 = true
	}
	if missenseHit[key1][hgmdCol] {
		flag2 = true
	}

	if flag1 && flag2 {
		return 1
	} else if flag1 || flag2 {
		return 2
	} else {
		return 0
	}
}
