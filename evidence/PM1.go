package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"log"
	"regexp"
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

func FindDomain(fileName, key, filterKey string, filter *regexp.Regexp) map[string]int {
	var DomainCount = make(map[string]int)
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		if !filter.MatchString(item[filterKey]) {
			continue
		}
		domain := item[key]
		if domain == "" || domain == "." {
			continue
		}
		domains := strings.Split(domain, ";")
		for _, d := range domains {
			if d == "" || d == "." {
				continue
			}
			DomainCount[d]++
		}
	}
	return DomainCount
}

// PM1
func CheckPM1(item map[string]string, ClinVarDbNSFP, ClinVarPfam, HGMDDbNSFP, HGMDPfam map[string]int) string {
	if !isMissenseIndel.MatchString(item["Function"]) {
		return ""
	}
	var dbNSFP = item["Interpro_domain"]
	var pfam = item["pfamId"]
	if ClinVarDbNSFP[dbNSFP] > 0 || ClinVarPfam[pfam] > 0 || HGMDDbNSFP[dbNSFP] > 0 || HGMDPfam[pfam] > 0 {
		return "1"
	}
	return "0"
}

func ComparePM1(item map[string]string, ClinVarDbNSFP, ClinVarPfam, HGMDDbNSFP, HGMDPfam map[string]int) {
	rule := "PM1"
	val := CheckPM1(item, ClinVarDbNSFP, ClinVarPfam, HGMDDbNSFP, HGMDPfam)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Interpro_domain", "pfamId")
	}
}
