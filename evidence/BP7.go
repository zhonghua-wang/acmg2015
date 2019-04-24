package evidence

import (
	"regexp"
)

var (
	isBP7Func = regexp.MustCompile(`coding-synon|intron|splice[+-]20`)
)

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckBP7(item map[string]string) string {
	if !isBP7Func.MatchString(item["Function"]) {
		return "0"
	}
	if item["GERP++_RS_pred"] == "不保守" &&
		item["PhyloP Vertebrates Pred"] == "不保守" &&
		item["PhyloP Placental Mammals Pred"] == "不保守" {
	} else {
		return "0"
	}
	if isP.MatchString(item["dbscSNV_RF_pred"]) &&
		isP.MatchString(item["dbscSNV_ADA_pred"]) {
	} else {
		return "0"
	}
	return "1"
	return ""
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
