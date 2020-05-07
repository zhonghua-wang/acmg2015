package evidence

// BP1
func CheckBP1(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}
	if bp1GeneList[item["Gene Symbol"]] {
		return "1"
	} else {
		return "0"
	}
}

func CompareBP1(item map[string]string, ClinVarBP1GeneList, HgmdBP1GeneList map[string]float64) {
	rule := "BP1"
	val := CheckBP1(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
