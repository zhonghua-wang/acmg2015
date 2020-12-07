package evidence

// PM5
func CheckPM5(item map[string]string) string {
	if item["Function"] != "missense" {
		return "0"
	}

	var key = item["Transcript"] + ":" + item["pHGVS1"]
	var AAPos = getAAPos.FindString(item["pHGVS1"])
	if AAPos == "" {
		return "0"
	}
	var key2 = item["Transcript"] + ":" + AAPos
	var countPHGVS = phgvsCount[key]
	var countAAPos = aaPostCount[key2]
	if countAAPos > countPHGVS {
		return "1"
	} else {
		return "0"
	}
}

func ComparePM5(item map[string]string, ClinVarPHGVSlist, ClinVarAAPosList, HGMDPHGVSlist, HGMDAAPosList map[string]int) {
	rule := "PM5"
	val := CheckPM5(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Function", "Transcript", "pHGVS")
	}
}
