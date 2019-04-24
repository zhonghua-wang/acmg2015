package evidence

import (
	"github.com/liserjrqlxue/simple-util"
	"regexp"
	"strings"
)

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
	var dbNSFP = item["Interpro_domain"]
	var pfam = item["pfamId"]
	if ClinVarDbNSFP[dbNSFP] > 0 || ClinVarPfam[pfam] > 0 || HGMDDbNSFP[dbNSFP] > 0 || HGMDPfam[pfam] > 0 {
		return "1"
	}
	return "0"
}
