package evidence

import "regexp"

var (
	isSplice      = regexp.MustCompile(`splice`)
	isD           = regexp.MustCompile(`D`)
	isDeleterious = regexp.MustCompile(`deleterious`)
)

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckPP3(item map[string]string) string {
	if isD.MatchString(item["GERP++_RS_pred"]) &&
		item["PhyloP Vertebrates Pred"] == "高度保守" &&
		item["PhyloP Placental Mammals Pred"] == "高度保守" {
	} else {
		return "0"
	}
	if isSplice.MatchString(item["Function"]) {
		if item["PVS1"] == "1" || item["PVS1"] == "5" {
			return "0"
		} else {
			if isD.MatchString(item["dbscSNV_ADA_pred"]) {
				return "1"
			} else {
				return "0"
			}
		}
	} else {
		if isD.MatchString(item["SIFT Pred"]) &&
			isD.MatchString(item["Polyphen2 HVAR Pred"]) &&
			isD.MatchString(item["MutationTaster Pred"]) &&
			isDeleterious.MatchString(item["Ens Condel Pred"]) {
			return "1"
		} else {
			return "0"
		}
	}
	return ""
}
