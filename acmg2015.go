package acmg2015

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
)

// colname
var (
	clinvarCol = "ClinVar Significance"
	hgmdCol    = "HGMD Pred"
)

// regexp
var (
	isSplice     = regexp.MustCompile(`splice[+-35]?$`)
	IsClinVarPLP = regexp.MustCompile(`Pathogenic|Likely_pathogenic`)
	IsHgmdDM     = regexp.MustCompile(`DM`)
)

// Tier1 >1
// LoF 3
var FuncInfo = map[string]int{
	"splice-3":     3,
	"splice-5":     3,
	"init-loss":    3,
	"alt-start":    3,
	"frameshift":   3,
	"nonsense":     3,
	"stop-gain":    3,
	"span":         3,
	"missense":     2,
	"cds-del":      2,
	"cds-indel":    2,
	"cds-ins":      2,
	"splice-10":    2,
	"splice+10":    2,
	"coding-synon": 1,
	"splice-20":    1,
	"splice+20":    1,
}

func AddACMG2015(inputData map[string]string) map[string]string {
	var LOFIntoleranceGeneList = simple_util.JsonFile2Map("db/LOFIntoleranceGeneList.json")
	var info = make(map[string]string)
	info["PVS1"] = CheckPVS1(inputData, LOFIntoleranceGeneList)
	inputData["ACMG"] = PredACMG2015(info)
	return inputData
}
