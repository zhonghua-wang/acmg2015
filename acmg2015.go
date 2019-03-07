package acmg2015

import (
	"regexp"
	"strconv"
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

// to be update
var LoFIntoleranceGene = map[string]bool{
	"ZNF574": true,
}

func AddACMG2015(inputData map[string]string) map[string]string {
	var info = make(map[string]string)
	info["PVS1"] = checkPVS1(inputData)
	inputData["ACMG"] = PredACMG2015(info)
	return inputData
}

// PVS1
func checkPVS1(inputData map[string]string) string {
	var function = inputData["Function"]
	if FuncInfo[function] < 3 {
		return "0"
	}
	var geneSymbol = inputData["Gene Symbol"]
	if !LoFIntoleranceGene[geneSymbol] {
		return "0"
	}
	if isSplice.MatchString(function) {
		score, err := strconv.ParseFloat(inputData["dbscSNV_RF_SCORE"], 32)
		if err == nil && score >= 0.6 {
			return "1"
		}
		score, err = strconv.ParseFloat(inputData["dbscSNV_ADA_SCORE"], 32)
		if err == nil && score >= 0.6 {
			return "1"
		}
		return "0"
	}
	return "1"
}
