package evidence

import (
	"strconv"
)

var BS2HomoList = []string{
	"ExAC HomoAlt Count",
	"GnomAD HomoAlt Count",
}

var BS2AF1List = []string{
	"GnomAD EAS AF",
	"GnomAD AF",
	"1000G AF",
	"ESP6500 AF",
	"ExAC EAS AF",
	"ExAC AF",
}

var (
	BS2LateOnsetHomoThreshold = 5
	BS2NoLateOnsetThreshold   = 0
	BS2HitCountThreshold      = 2
)

func CheckBS2(item map[string]string) string {
	var homoCount = 0
	for _, key := range BS2HomoList {
		c, e := strconv.Atoi(item[key])
		if e == nil {
			homoCount += c
		} else {
			c = 0
		}
		if bs2GeneList[item["Gene Symbol"]] {
			if c >= BS2LateOnsetHomoThreshold {
				return "1"
			}
		} else {
			if c > BS2NoLateOnsetThreshold {
				return "1"
			}
		}
	}

	var hitCount = 0
	var inherit = item["ModeInheritance"]
	if !isARDRXLNA.MatchString(inherit) && isADYL.MatchString(inherit) {
		for _, key := range BS2AF1List {
			if !CheckAFAllLowThen(item, []string{key}, 0, true) {
				hitCount++
			}
		}
		if hitCount >= BS2HitCountThreshold {
			return "1"
		}
	}
	return "0"
}

func CompareBS2(item map[string]string) {
	rule := "BS2"
	val := CheckBS2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, append([]string{"Gene Symbol", "ModeInheritance"}, append(BS2HomoList, BS2AF1List...)...)...)
	}
}
