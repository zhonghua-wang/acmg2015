package evidence

import (
	"regexp"
)

var (
	isBP7Func = regexp.MustCompile(`coding-synon`)
)

// ture	:	"1"
// flase:	"0"
func CheckBP7(item map[string]string) string {
	if !isBP7Func.MatchString(item["Function"]) {
		return "0"
	}
	if !(item["GERP++_RS_pred"] == "不保守" &&
		item["PhyloP Vertebrates Pred"] == "不保守" &&
		item["PhyloP Placental Mammals Pred"] == "不保守") {
		return "0"
	}
	if isD.MatchString(item["dbscSNV_RF_pred"]) ||
		isD.MatchString(item["dbscSNV_ADA_pred"]) ||
		isD.MatchString(item["SpliceAI Pred"]) {
		return "0"
	}
	if isP.MatchString(item["SpliceAI Pred"]) {
		return "1"
	}
	return "0"
}

func CompareBP7(item map[string]string, lostOnly bool) {
	rule := "BP7"
	val := CheckBP7(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			if lostOnly && val != "1" {
				return
			}
			PrintConflict(
				item,
				rule,
				val,
				"GERP++_RS_pred",
				"PhyloP Vertebrates Pred",
				"PhyloP Placental Mammals Pred",
				"Function",
				"dbscSNV_RF_pred",
				"dbscSNV_ADA_pred",
			)
		}
	}
}
