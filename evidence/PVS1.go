package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
	"strconv"
)

func FindLOFIntoleranceGeneList(fileName, key string, pathogenicRegexp *regexp.Regexp) map[string]int {
	var geneList = make(map[string]int)
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !pathogenicRegexp.MatchString(item[key]) {
			continue
		}
		if FuncInfo[item["Function"]] < 3 {
			continue
		}
		if !CheckAFLowThen(item, 0.05) {
			continue
		}
		geneList[item["Gene Symbol"]]++
	}
	return geneList
}

var AFlist = []string{
	"GnomAD EAS AF",
	"GnomAD AF",
	"1000G AF",
	"ESP6500 AF",
	"ExAC EAS AF",
	"ExAC AF",
	"PVFD AF",
	"Panel AlleleFreq",
}

func CheckAFLowThen(item map[string]string, threshold float64) bool {
	for _, key := range AFlist {
		af := item[key]
		if af == "" || af == "." {
			continue
		}
		AF, err := strconv.ParseFloat(af, 64)
		simple_util.CheckErr(err)
		if AF > threshold {
			return false
		}
	}
	return true
}

func CheckPVS1(item map[string]string, LOFList map[string]int) string {
	if FuncInfo[item["Function"]] < 3 {
		return "0"
	}
	if LOFList[item["Gene Symbol"]] == 0 {
		return "0"
	}
	if CheckDomain(item) {
		return "1"
	}
	if CheckOtherPathogenic(item) {
		return "1"
	}
	return "0"
}

func ComparePVS1(item map[string]string, LOFList map[string]int) {
	rule := "PVS1"
	val := CheckPVS1(item, LOFList)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			PrintConflict(item, rule, val, []string{"Function", "Gene Symbol"})
		}
	}
}

// 突变位点后有重要的蛋白结构功能区域
func CheckDomain(item map[string]string) bool {
	return false
}

// 突变位点后有其他致病突变（基于公共数据库）位点
func CheckOtherPathogenic(item map[string]string) bool {
	return false
}
