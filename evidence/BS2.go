package evidence

import (
	"github.com/liserjrqlxue/simple-util"
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
	BS2LateOnsetHomoThreshold = 4
	BS2NoLateOnsetThreshold   = 0
	BS2HitCountThreshold      = 2
)

func GetLateOnsetList(fileName string) map[string]int {
	var lateOnsetList = make(map[string]int)
	for _, key := range simple_util.File2Array(fileName) {
		lateOnsetList[key] = 1
	}
	return lateOnsetList
}

func CheckBS2(item map[string]string, lateOnsetList map[string]int) string {
	var homoCount = 0
	for _, key := range BS2HomoList {
		c, e := strconv.Atoi(item[key])
		if e == nil {
			homoCount += c
		}
	}
	if lateOnsetList[item["Gene Symbol"]] > 0 {
		if homoCount >= BS2LateOnsetHomoThreshold {
			return "1"
		}
	} else {
		if homoCount >= BS2NoLateOnsetThreshold {
			return "1"
		}
	}
	var hitCount = 0
	var inherit = item["OMIM inheritance"]
	if !isARXRNA.MatchString(inherit) && isADXD.MatchString(inherit) {
		for _, key := range BS2AF1List {
			if CheckAFAllLowThen(item, []string{key}, 0, true) {
				hitCount++
			}
		}
		if hitCount >= BS2HitCountThreshold {
			return "1"
		}
	}
	return ""
}

func CompareBS2(item map[string]string, lateOnsetList map[string]int) {
	rule := "BS2"
	val := CheckBS2(item, lateOnsetList)
	if val != item[rule] {
		PrintConflict(item, rule, val, append([]string{"Gene Symbol", "OMIM inheritance"}, append(BS2HomoList, BS2AF1List...)...)...)
	}
}
