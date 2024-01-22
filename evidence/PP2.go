package evidence

// PP2 deprecated PP2 rule
func CheckPP2_OLD(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}
	if pp2GeneList[item["Gene Symbol"]] {
		return "1"
	} else {
		return "0"
	}
}

func CheckPP2(item map[string]string) string {
	// get gene score from pp2ZscoreMap
	gene := item["Gene Symbol"]
	score, ok := pp2ZscoreMap[gene]
	// if ok and score > 3.09, return 1
	if ok && score >= 3.09 {
		return "1"
	}
	return "0"
}

func ComparePP2(item map[string]string) {
	rule := "PP2"
	val := CheckPP2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Gene Symbol")
	}
}
