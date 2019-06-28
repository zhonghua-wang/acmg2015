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

type filterFunc func(item map[string]string) bool

// colname
var (
	clinvarCol      = "ClinVar Significance"
	hgmdCol         = "HGMD Pred"
	domainDbNSFPCol = "Interpro_domain"
	domainPfamCol   = "pfamId"
)

// regexp
var (
	isMissenseIndel = regexp.MustCompile(`missense|ins|del`)
)

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
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
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
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(rdr)
	for {
		_, err := rdr.Next()
		if err == io.EOF {
			break
		}
		simple_util.CheckErr(err)
		n++
	}
	return n
}

// PM1
func CheckPM1(item map[string]string, dbNSFPDomain, PfamDomain map[string]bool, tbx *bix.Bix) string {
	if !isMissenseIndel.MatchString(item["Function"]) {
		return "0"
	}
	var dbNSFP = item["Interpro_domain"]
	var pfam = item["pfamId"]
	var flag bool

	for _, k := range strings.Split(dbNSFP, ";") {
		if dbNSFPDomain[k] {
			flag = true
		}
	}
	for _, k := range strings.Split(pfam, ";") {
		if PfamDomain[k] {
			flag = true
		}
	}
	if !flag {
		chr := strings.Replace(item["#Chr"], "chr", "", 1)
		start, err := strconv.Atoi(item["Start"])
		simple_util.CheckErr(err)
		end, err := strconv.Atoi(item["Stop"])
		simple_util.CheckErr(err)
		n := countBix(tbx, chr, start-10, end+10)
		if n >= 2 {
			flag = true
		}
	}
	if flag {
		return "1"
	} else {
		return "0"
	}
	return "0"
}

func ComparePM1(item map[string]string, dbNSFPDomain, PfamDomain map[string]bool, tbx *bix.Bix) {
	rule := "PM1"
	val := CheckPM1(item, dbNSFPDomain, PfamDomain, tbx)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Interpro_domain", "pfamId")
	}
}
