package evidence

import (
	"fmt"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"regexp"
	"strconv"
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

var (
	IsClinVarPLP  = regexp.MustCompile(`Pathogenic|Likely_pathogenic`)
	IsClinVarBLB  = regexp.MustCompile(`Benign|Likely_benign`)
	IsHgmdDM      = regexp.MustCompile(`DM$|DM\|`)
	IsHgmdB       = regexp.MustCompile(`DP|FP|DFP`)
	isARXLNA      = regexp.MustCompile(`AR|XL|NA`)
	isAD          = regexp.MustCompile(`AD`)
	isSplice      = regexp.MustCompile(`splice`)
	isSplice20    = regexp.MustCompile(`splice[+-]20`)
	isD           = regexp.MustCompile(`D`)
	isDeleterious = regexp.MustCompile(`deleterious`)
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
		simple_util.CheckErr(err)
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
