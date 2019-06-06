package evidence

import (
	"github.com/brentp/bix"
	"github.com/brentp/irelate/interfaces"
	"github.com/liserjrqlxue/simple-util"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var PVS1AFThreshold = 0.05

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
		if !CheckAFAllLowThen(item, PVS1AFlist, PVS1AFThreshold, true) {
			continue
		}
		geneList[item["Gene Symbol"]]++
	}
	return geneList
}

var PVS1AFlist = []string{
	"GnomAD EAS AF",
	"GnomAD AF",
	"1000G AF",
	"ESP6500 AF",
	"ExAC EAS AF",
	"ExAC AF",
}

func CheckPVS1(item map[string]string, LOFList map[string]int, transcriptInfo map[string][]Region, tbx *bix.Bix) string {
	if FuncInfo[item["Function"]] < 3 {
		return "0"
	}
	if LOFList[item["Gene Symbol"]] == 0 {
		return "0"
	}
	if CheckDomain(item) {
		return "1"
	}
	regions, ok := transcriptInfo[item["Transcript"]]
	if ok {
		if CheckOtherPathogenic(tbx, item, regions) {
			return "1"
		}
	} else {
		log.Printf("Transcript(%s)of Gene(%s) not in transcriptInfo\n", item["Transcript"], item["Gene Symbol"])
	}

	return "0"
}

func ComparePVS1(item map[string]string, LOFList map[string]int, transcriptInfo map[string][]Region, tbx *bix.Bix) {
	rule := "PVS1"
	val := CheckPVS1(item, LOFList, transcriptInfo, tbx)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			PrintConflict(item, rule, val, "Function", "Gene Symbol")
		}
	}
}

// 突变位点后有重要的蛋白结构功能区域
func CheckDomain(item map[string]string) bool {
	return false
}

// 突变位点后有其他致病突变（基于公共数据库）位点
func CheckOtherPathogenic(tbx *bix.Bix, item map[string]string, regions []Region) bool {
	chromosome := strings.Replace(item["#Chr"], "chr", "", -1)
	start, err := strconv.Atoi(item["Start"])
	simple_util.CheckErr(err)
	end := start
	var flag bool
	for _, item := range regions {
		if item.Chromosome != chromosome {
			continue
		}
		if item.Strand == "+" {
			if int(item.End) >= end {
				flag = true
				end = int(item.End)
			}
		} else if item.Strand == "-" {
			if int(item.Start) <= start {
				flag = true
				start = int(item.Start)
			}
		} else {
			log.Printf("unknown Strand(%s):%+v", item.Strand, item)
		}
		if start > end {
			flag = false
		}
	}
	if flag {
		return checkPathogenicOfRegion(tbx, chromosome, start, end)
	} else {
		log.Printf("can not set region after this variant:[%s]\n", item["MutationName"])
	}
	return false
}

func checkPathogenicOfRegion(tbx *bix.Bix, chromosome string, start, end int) bool {
	rdr, err := tbx.Query(interfaces.AsIPosition(chromosome, start, end))
	simple_util.CheckErr(err)
	for {
		_, err := rdr.Next()
		if err == io.EOF {
			break
		}
		simple_util.CheckErr(err)
		return true
	}
	return false
}
