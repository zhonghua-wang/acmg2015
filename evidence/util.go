package evidence

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/stringsUtil"
)

type Region struct {
	Seqid string
	//Source string
	Type       string
	Chromosome string
	Start      uint64
	End        uint64
	Strand     string
	Gene       string
}

// regexp
var (
	getAAPos      = regexp.MustCompile(`^p\.[A-Z]\d+`)
	IsClinVarPLP  = regexp.MustCompile(`Pathogenic|Likely_pathogenic`)
	IsClinVarBLB  = regexp.MustCompile(`Benign|Likely_benign`)
	IsHgmdDM      = regexp.MustCompile(`DM$|DM\|`)
	IsHgmdB       = regexp.MustCompile(`DP|FP|DFP`)
	isARDRXLNA    = regexp.MustCompile(`AR|DR|XL|NA`)
	isARDRNA      = regexp.MustCompile(`AR|DR|NA`)
	isADYL        = regexp.MustCompile(`AD|YL`)
	isADXLYL      = regexp.MustCompile(`AD|XL|YL`)
	isSplice      = regexp.MustCompile(`splice`)
	isSplice20    = regexp.MustCompile(`splice[+-]20`)
	isD           = regexp.MustCompile(`D`)
	isDeleterious = regexp.MustCompile(`deleterious`)
	repeatSeq     = regexp.MustCompile(`c\..*\[(\d+)>\d+]`)
)

// Tier1 >1
// LoF 3
var FuncInfo = map[string]int{
	"splice-3":     3,
	"splice-5":     3,
	"init-loss":    3,
	"alt-start":    3,
	"frameshift":   3,
	"nonsense":     3,
	"stop-gain":    3,
	"span":         3,
	"missense":     2,
	"cds-del":      2,
	"cds-indel":    2,
	"cds-ins":      2,
	"splice-10":    2,
	"splice+10":    2,
	"coding-synon": 1,
	"splice-20":    1,
	"splice+20":    1,
}

func CheckAFAllLowThen(item map[string]string, AFlist []string, threshold float64, includeEqual bool) bool {
	for _, key := range AFlist {
		af := item[key]
		if af == "" || af == "." || af == "0" {
			continue
		}
		AF, err := strconv.ParseFloat(af, 64)
		simpleUtil.CheckErr(err)
		if includeEqual {
			if AF > threshold {
				return false
			}
		} else {
			if AF >= threshold {
				return false
			}
		}
	}
	return true
}

func PrintConflict(item map[string]string, rule, val string, keys ...string) {
	fmt.Fprintf(
		os.Stderr,
		"Conflict %s:[%s] vs [%s]\t%s[%s]\n",
		rule,
		val,
		item[rule],
		"MutationName",
		item["MutationName"],
	)
	for _, key := range keys {
		fmt.Fprintf(os.Stderr, "\t%30s:[%s]\n", key, item[key])
	}
}

var (
	hgvsCount   = make(map[string]int)
	phgvsCount  = make(map[string]int)
	aaPostCount = make(map[string]int)
)

func LoadPS1PM5(hgvs, pHgvs, aaPos string) {
	hgvsCount = tsv2mapStringInt(hgvs)
	phgvsCount = tsv2mapStringInt(pHgvs)
	aaPostCount = tsv2mapStringInt(aaPos)
}

func tsv2mapStringInt(tsv string) map[string]int {
	var db = make(map[string]int)

	var file = osUtil.Open(tsv)
	defer simpleUtil.DeferClose(file)

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var array = strings.Split(scanner.Text(), "\t")
		array = append(array, "NA", "NA")
		if array[1] == "NA" {
			array[1] = "0"
		}
		var v, ok = db[array[0]]
		if ok {
			var vStr = strconv.Itoa(v)
			if array[1] != vStr {
				panic("dup key[" + array[0] + "],different value:[" + vStr + "]vs[" + array[1] + "]")
			}
		} else {
			db[array[0]] = stringsUtil.Atoi(array[1])
		}
	}
	simpleUtil.CheckErr(scanner.Err())
	return db
}
