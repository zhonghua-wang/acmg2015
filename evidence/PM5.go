package evidence

// PS1
func CheckPM5(item map[string]string, ClinVarPHGVSlist, ClinVarAAPosList, HGMDPHGVSlist, HGMDAAPosList map[string]int) string {
	if item["Function"] == "missense" {
		return "0"
	}

	var key = item["Transcript"] + ":" + item["pHGVS"]
	var AAPos = getAAPos.FindString(item["pHGVS"])
	if AAPos == "" {
		return "0"
	}
	var key2 = item["Transcript"] + ":" + AAPos
	var flag1, flag2 bool
	ClinVarPHGVSNo := ClinVarPHGVSlist[key]
	HGMDPHGVSNo := HGMDPHGVSlist[key]
	ClinvarAAPosNo := ClinVarAAPosList[key2]
	HGMDAAPosNo := HGMDAAPosList[key2]
	if ClinvarAAPosNo > ClinVarPHGVSNo {
		flag1 = true
	}
	if HGMDAAPosNo > HGMDPHGVSNo {
		flag2 = true
	}

	if flag1 && flag2 {
		return "2"
	} else if flag1 || flag2 {
		return "1"
	} else {
		return "0"
	}
}

func ComparePM5(item map[string]string, ClinVarPHGVSlist, ClinVarAAPosList, HGMDPHGVSlist, HGMDAAPosList map[string]int) {
	rule := "PM5"
	val := CheckPM5(item, ClinVarPHGVSlist, ClinVarAAPosList, HGMDPHGVSlist, HGMDAAPosList)
	if val != item[rule] {
		PrintConflict(item, rule, val, []string{"Function", "Transcript", "pHGVS"})
	}
}
