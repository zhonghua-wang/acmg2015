package evidence

// PP2
func CheckPP2(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}
	if pp2GeneList[item["Gene Symbol"]] {
		return "1"
	} else {
		return "0"
	}
}

func ComparePP2(item map[string]string) {
	rule := "PP2"
	val := CheckPP2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
