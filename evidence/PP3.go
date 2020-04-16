package evidence

// ture	:	"1"
// flase:	"0"
func CheckPP3(item map[string]string) string {
	var count = 0
	for _, pred := range []string{
		item["GERP++_RS_pred"],
		item["PhyloP Vertebrates Pred"],
		item["PhyloP Placental Mammals Pred"],
	} {
		if pred == "不保守" {
			return "0"
		} else if pred == "保守" {
			count++
		}
	}
	if count < 2 {
		return "0"
	}

	count = 0
	if isNeutral.MatchString(item["Ens Condel Pred"]) {
		return "0"
	} else if isDeleterious.MatchString(item["Ens Condel Pred"]) {
		count++
	}
	for _, pred := range []string{
		item["SIFT Pred"],
		item["MutationTaster Pred"],
		item["Polyphen2 HVAR Pred"],
		item["dbscSNV_RF_pred"],
		item["dbscSNV_ADA_pred"],
		item["SpliceAI Pred"],
	} {
		if isP.MatchString(pred) {
			return "0"
		} else if isD.MatchString(pred) {
			count++
		}
	}
	if item["Function"] == "splice-3" || item["Function"] == "splice-5" {
		if item["PVS1"] == "1" {
			return "0"
		}
	} else if isSplice.MatchString(item["Function"]) || item["Function"] == "intron" {
		if isP.MatchString(item["SpliceAI Pred"]) {
			return "1"
		}
	}
	if count >= 2 {
		return "1"
	}
	return "0"
}

func ComparePP3(item map[string]string, lostOnly bool) {
	rule := "PP3"
	val := CheckPP3(item)
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
				"PVS1",
				"dbscSNV_RF_pred",
				"dbscSNV_ADA_pred",
				"SIFT Pred",
				"Polyphen2 HVAR Pred",
				"MutationTaster Pred",
				"Ens Condel Pred",
			)
		}
	}
}
