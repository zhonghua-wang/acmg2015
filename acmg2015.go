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
	inputData["ACMG"] = predACMG2015(info)
	return inputData
}

func predACMG2015(info map[string]string) string {

	PVS1 := info["PVS1"]

	PS1 := info["PS1"]
	PS2 := info["PS2"]
	PS3 := info["PS3"]
	PS4 := info["PS4"]

	PM1 := info["PM1"]
	PM2 := info["PM2"]
	PM3 := info["PM3"]
	PM4 := info["PM4"]
	PM5 := info["PM5"]
	PM6 := info["PM6"]

	PP1 := info["PP1"]
	PP2 := info["PP2"]
	PP3 := info["PP3"]
	PP4 := info["PP4"]
	PP5 := info["PP5"]

	BA1 := info["BA1"]

	BS1 := info["BS1"]
	BS2 := info["BS2"]
	BS3 := info["BS3"]
	BS4 := info["BS4"]

	BP1 := info["BP1"]
	BP2 := info["BP2"]
	BP3 := info["BP3"]
	BP4 := info["BP4"]
	BP5 := info["BP5"]
	BP6 := info["BP6"]
	BP7 := info["BP7"]

	// PVS
	//  PVS1 5 得分
	//  PVS1 6 不得分
	if PVS1 == "5" {
		PVS1 = "1"
	}
	if PVS1 == "6" {
		PVS1 = "0"
	}
	var sumPVS int
	if PVS1 == "1" {
		sumPVS++
	}

	// PS
	//  PS1 2 暂时不得分
	//  PS1 3 得分
	//  PS1 4 不得分
	if PS1 == "2" || PS1 == "4" {
		PS1 = "0"
	}
	if PS1 == "3" {
		PS1 = "1"
	}
	//  PS4 5 得分
	if PS4 == "5" {
		PS4 = "1"
	}
	var sumPS int
	if PS1 == "1" {
		sumPS++
	}
	if PS2 == "1" {
		sumPS++
	}
	if PS3 == "1" {
		sumPS++
	}
	if PS4 == "1" {
		sumPS++
	}

	// PM
	//  不与PS4 同时得分
	if PS4 == "1" {
		PM2 = "0"
	}
	//  PM3 2 升级到PS得分
	if PM3 == "2" {
		sumPS++
		PM3 = "0"
	}
	//  PM5 2 暂时得分
	//  PM5 3 得分
	//  PM5 4 不得分
	//  PM5 5 得分
	if PM5 == "2" || PM5 == "3" || PM5 == "5" {
		PM5 = "1"
	}
	if PM5 == "4" {
		PM5 = "0"
	}
	var sumPM int
	if PM1 == "1" {
		sumPM++
	}
	if PM2 == "1" {
		sumPM++
	}
	if PM3 == "1" {
		sumPM++
	}
	if PM4 == "1" {
		sumPM++
	}
	if PM5 == "1" {
		sumPM++
	}
	if PM6 == "1" {
		sumPM++
	}

	// PP
	//  PP1 2 升级到PM
	if PP1 == "2" {
		sumPM++
		PP1 = "0"
	}
	//  PS1/PM5  PP5 不共存
	if PS1 == "1" || PM5 == "1" {
		PP5 = "0"
	}
	//  ACMG 已取消该证据
	PP5 = "0"
	var sumPP int
	if PP1 == "1" {
		sumPP++
	}
	if PP2 == "1" {
		sumPP++
	}
	if PP3 == "1" {
		sumPP++
	}
	if PP4 == "1" {
		sumPP++
	}
	if PP5 == "1" {
		sumPP++
	}

	// BA
	var sumBA int
	if BA1 == "1" {
		sumBA++
	}
	// BS
	var sumBS int
	if BS1 == "1" {
		sumBS++
	}
	if BS2 == "1" {
		sumBS++
	}
	if BS3 == "1" {
		sumBS++
	}
	if BS4 == "1" {
		sumBS++
	}
	// BP
	var sumBP int
	if BP1 == "1" {
		sumBP++
	}
	if BP2 == "1" {
		sumBP++
	}
	if BP3 == "1" {
		sumBP++
	}
	if BP4 == "1" {
		sumBP++
	}
	if BP5 == "1" {
		sumBP++
	}
	if BP6 == "1" {
		sumBP++
	}
	if BP7 == "1" {
		sumBP++
	}

	var ACMG = make(map[string]bool)
	if sumPVS > 0 {
		if sumPS == 1 || sumPM > 1 || sumPP > 1 || (sumPM == 1 && sumPP == 1) {
			ACMG["P"] = true
		}
		if sumPM == 1 {
			ACMG["LP"] = true
		}
	}
	if sumPS > 1 {
		ACMG["P"] = true
	}
	if sumPS == 1 {
		if sumPM > 2 || (sumPM == 2 && sumPP >= 2) || (sumPM == 1 && sumPP >= 4) {
			ACMG["P"] = true
		}
		if sumPM == 1 || sumPM == 2 || sumPP > 1 {
			ACMG["LP"] = true
		}
	}
	if sumPM > 2 || (sumPM == 2 && sumPP > 1) || (sumPM == 1 && sumPP > 3) {
		ACMG["LP"] = true
	}
	if sumBA > 0 || sumBS > 1 {
		ACMG["B"] = true
	}
	if sumBP > 1 || (sumBP == 1 && sumBS == 1) {
		ACMG["LB"] = true
	}
	var PLP, BLB bool
	if ACMG["P"] || ACMG["LP"] {
		PLP = true
	}
	if ACMG["B"] || ACMG["LB"] {
		BLB = true
	}
	if PLP && BLB {
		return "VUS"
	} else if ACMG["P"] {
		return "P"
	} else if ACMG["LP"] {
		return "LP"
	} else if ACMG["B"] {
		return "B"
	} else if ACMG["LB"] {
		return "LB"
	} else {
		return "VUS"
	}
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
