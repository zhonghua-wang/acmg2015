package evidence

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/brentp/bix"
	"github.com/brentp/irelate/interfaces"
	"github.com/liserjrqlxue/goUtil/fmtUtil"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/stringsUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
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
	PM1Function  = regexp.MustCompile(`missense|cds-indel`)
	getAAPos     = regexp.MustCompile(`^p\.[A-Z]\d+`)
	IsClinVarPLP = regexp.MustCompile(`Pathogenic|Likely_pathogenic`)

	IsHgmdDM = regexp.MustCompile(`DM$|DM\|`)

	isARDRXLPRDDNA = regexp.MustCompile(`AR|DR|XL|PR|DD|NA|UNK`)
	isADPDYL       = regexp.MustCompile(`AD|PD|YL`)
	isSplice       = regexp.MustCompile(`splice`)
	isSpliceIntron = regexp.MustCompile(`splice|intron`)
	isSplice20     = regexp.MustCompile(`splice[+-]20`)
	isP            = regexp.MustCompile(`P`)
	isD            = regexp.MustCompile(`D`)
	isI            = regexp.MustCompile(`I`)
	isNeutral      = regexp.MustCompile(`neutral`)
	isDeleterious  = regexp.MustCompile(`deleterious`)
	repeatSeq      = regexp.MustCompile(`c\..*\[(\d+)>\d+]`)
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
	fmtUtil.Fprintf(
		os.Stderr,
		"Conflict %s:[%s] vs [%s]\t%s[%s]\n",
		rule,
		val,
		item[rule],
		"MutationName",
		item["MutationName"],
	)
	for _, key := range keys {
		fmtUtil.Fprintf(os.Stderr, "\t%30s:[%s]\n", key, item[key])
	}
}

var (
	hgvsCount         = make(map[string]int)
	phgvsCount        = make(map[string]int)
	aaPostCount       = make(map[string]int)
	pm1PfamId         = make(map[string]bool)
	pm1InterproDomain = make(map[string]bool)
	bp1GeneList       = make(map[string]bool)
	bs2GeneList       = make(map[string]bool)
	ba1Exception      = make(map[string]bool)
	pp2GeneList       = make(map[string]bool)
	pp2ZscoreMap      = make(map[string]float32)
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

func LoadPM1(pfamId, interproDomain string) {
	var pfamIdArray = textUtil.File2Array(pfamId)
	var interproDomainArray = textUtil.File2Array(interproDomain)
	for _, key := range pfamIdArray {
		pm1PfamId[key] = true
	}
	for _, key := range interproDomainArray {
		pm1InterproDomain[key] = true
	}
}

func countBix(tbx *bix.Bix, chr string, start, end int) (n int) {
	rdr, err := tbx.Query(interfaces.AsIPosition(chr, start, end))
	simpleUtil.CheckErr(err)
	defer simpleUtil.DeferClose(rdr)
	for {
		_, err := rdr.Next()
		if err == io.EOF {
			break
		}
		simpleUtil.CheckErr(err)
		n++
	}
	return n
}

func LoadBP1(bp1geneList string) {
	var genes = textUtil.File2Array(bp1geneList)
	bp1GeneList = make(map[string]bool)
	for _, gene := range genes {
		bp1GeneList[gene] = true
	}
}

func LoadBS2(fileName string) {
	for _, key := range textUtil.File2Array(fileName) {
		bs2GeneList[key] = true
	}
}

func LoadBA1(fileName string) {
	for _, key := range textUtil.File2Array(fileName) {
		ba1Exception[key] = true
	}
	return
}

func LoadPP2(pp2geneList string) {
	var genes = textUtil.File2Array(pp2geneList)
	pp2GeneList = make(map[string]bool)
	for _, gene := range genes {
		pp2GeneList[gene] = true
	}
}

// load missense z-score data from file
func LoadPP2Zscore(pp2ZscoreFile string) {
	var file = osUtil.Open(pp2ZscoreFile)
	defer simpleUtil.DeferClose(file)
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var array = strings.Split(scanner.Text(), "\t")
		gene := array[0]
		zscore, err := strconv.ParseFloat(array[4], 32)
		simpleUtil.CheckErr(err)
		pp2ZscoreMap[gene] = float32(zscore)
	}
	simpleUtil.CheckErr(scanner.Err())
}
