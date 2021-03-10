package evidence

// ture	:	"1"
// flase:	"0"
func CheckPP3(item map[string]string, autoPVS1 bool) string {
	// PVS1排除
	if autoPVS1 {
		switch item["AutoPVS1 Adjusted Strength"] {
		case "VeryStrong":
			return "0"
		case "Strong":
			return "0"
		case "Moderate":
			return "0"
		case "Supporting":
			return "0"
		}
	} else if item["PVS1"] == "1" {
		return "0"
	}

	// 保守性
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

	// 有害性
	// b1
	var b1 = true
	count = 0
	if isNeutral.MatchString(item["Ens Condel Pred"]) {
		b1 = false
	} else if isDeleterious.MatchString(item["Ens Condel Pred"]) {
		count++
	}
	for _, pred := range []string{
		item["SIFT Pred"],
		item["MutationTaster Pred"],
		item["Polyphen2 HVAR Pred"],
	} {
		if isP.MatchString(pred) || isI.MatchString(pred) {
			b1 = false
		} else if isD.MatchString(pred) {
			count++
		}
	}
	if count < 2 {
		b1 = false
	}

	// b2
	var b2 = true
	count = 0
	for _, pred := range []string{
		item["dbscSNV_RF_pred"],
		item["dbscSNV_ADA_pred"],
		item["SpliceAI Pred"],
	} {
		if isP.MatchString(pred) || isI.MatchString(pred) {
			b2 = false
		} else if isD.MatchString(pred) {
			count++
		}
	}
	if count < 2 && !isD.MatchString(item["SpliceAI Pred"]) {
		b2 = false
	}
	if b1 || b2 {
		return "1"
	}
	return "0"
}

func ComparePP3(item map[string]string, lostOnly, autoPVS1 bool) {
	rule := "PP3"
	val := CheckPP3(item, autoPVS1)
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
