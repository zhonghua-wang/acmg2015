package evidence

import (
	"fmt"
	"os"
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
	if isP.MatchString(item["GERP++_RS_pred"]) &&
		(item["PhyloP Vertebrates Pred"] == "保守" || item["PhyloP Vertebrates Pred"] == "不保守") &&
		(item["PhyloP Placental Mammals Pred"] == "保守" || item["PhyloP Placental Mammals Pred"] == "不保守") {
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
			fmt.Fprintf(
				os.Stderr,
				"Conflict %s:[%s] vs [%s]\t%s[%s]\n",
				rule,
				val,
				item[rule],
				"MutationName",
				item["MutationName"],
			)
			for _, key := range []string{
				"GERP++_RS_pred",
				"PhyloP Vertebrates Pred",
				"PhyloP Placental Mammals Pred",
				"Function",
				"dbscSNV_RF_pred",
				"dbscSNV_ADA_pred",
				"SIFT Pred",
				"Polyphen2 HVAR Pred",
				"MutationTaster Pred",
				"Ens Condel Pred",
			} {
				fmt.Fprintf(os.Stderr, "\t%30s:[%s]\n", key, item[key])
			}
		}
	}
}
