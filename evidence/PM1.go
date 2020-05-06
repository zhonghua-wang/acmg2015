package evidence

import (
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/brentp/bix"
	"github.com/brentp/irelate/interfaces"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/stringsUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
)

type filterFunc func(item map[string]string) bool

// colname
var (
	clinvarCol      = "ClinVar Significance"
	hgmdCol         = "HGMD Pred"
	domainDbNSFPCol = "Interpro_domain"
	domainPfamCol   = "pfamId"
	ext             = 3

	PM1pfamId         = make(map[string]bool)
	PM1InterproDomain = make(map[string]bool)
)

// regexp
var (
	isMissenseIndel = regexp.MustCompile(`missense|cds-ins|cds-del|cds-indel`)
)

func LoadPM1(pfamId, interproDomain string) {
	var pfamIdArray = textUtil.File2Array(pfamId)
	var interproDomainArray = textUtil.File2Array(interproDomain)
	for _, key := range pfamIdArray {
		PM1pfamId[key] = true
	}
	for _, key := range interproDomainArray {
		PM1InterproDomain[key] = true
	}
}

func FilterPathogenic(item map[string]string) (keep bool) {
	if IsClinVarPLP.MatchString(item[clinvarCol]) || IsHgmdDM.MatchString(item[hgmdCol]) {
		return true
	}
	return
}

func FilterBenign(item map[string]string) (keep bool) {
	if IsClinVarBLB.MatchString(item[clinvarCol]) || IsHgmdB.MatchString(item[hgmdCol]) {
		return true
	}
	return
}

func FindPM1MutationDomain(fileName string, filter filterFunc) (mutationDomain map[string][]string) {
	mutationDomain = make(map[string][]string)
	itemArray, _ := textUtil.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !filter(item) {
			continue
		}
		if !isMissenseIndel.MatchString(item["Function"]) {
			continue
		}
		var domains []string
		for _, col := range []string{domainDbNSFPCol, domainPfamCol} {
			domains = append(domains, item[col])
		}
		key := strings.Join([]string{item["#Chr"], item["Start"], item["Stop"], item["MutationName"]}, "\t")
		_, ok := mutationDomain[key]
		if ok {
			log.Printf("[Duplicate Mutatin:%s]\n", key)
			//d=append(d,domains...)
		} else {
			mutationDomain[key] = domains
		}
	}
	return
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

// PM1
func CheckPM1(item map[string]string, tbx *bix.Bix) string {
	if !isMissenseIndel.MatchString(item["Function"]) {
		return "0"
	}
	var dbNSFP = item["Interpro_domain"]
	var pfam = item["pfamId"]
	var flag bool

	for _, k := range strings.Split(dbNSFP, ";") {
		if PM1InterproDomain[k] {
			flag = true
		}
	}
	for _, k := range strings.Split(pfam, ";") {
		if PM1pfamId[k] {
			flag = true
		}
	}
	if !flag {
		var chr = strings.Replace(item["#Chr"], "chr", "", 1)
		var start = stringsUtil.Atoi(item["Start"])
		var end = stringsUtil.Atoi(item["Stop"])
		n := countBix(tbx, chr, start-ext, end+ext)
		if n >= 2 {
			flag = true
		}
	}
	if flag {
		return "1"
	} else {
		return "0"
	}
}

func ComparePM1(item map[string]string, tbx *bix.Bix) {
	rule := "PM1"
	val := CheckPM1(item, tbx)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Interpro_domain", "pfamId")
	}
}
