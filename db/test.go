package main

import (
	"encoding/json"
	"fmt"
	"github.com/liserjrqlxue/acmg2015/evidence"
	"github.com/liserjrqlxue/parse-gff3"
	"github.com/liserjrqlxue/simple-util"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// os
var (
	ex, _ = os.Executable()
	//exPath = filepath.Dir(ex)
	//pSep   = string(os.PathSeparator)
	//dbPath       = exPath + pSep + "db" + pSep
	//templatePath = exPath + pSep + "template" + pSep
)

var (
	clinvarAnno = "clinvar.2019-02-19.vcf.gz.2019-02-26.anno.bed.update"
	hgmdAnno    = "hgmd_pro_2018.3_hg19.vcf.2019-02-27.anno.bed.update"
	exacLOF     = "LOF_Intolerance.ExAC.lst"
)

// colname
var (
	clinvarCol = "ClinVar Significance"
	hgmdCol    = "HGMD Pred"
)

var columns = []string{
	"#Chr",
	"Start",
	"Stop",
	"Ref",
	"Call",
	"MutationName",
	clinvarCol,
	hgmdCol,
}

func main() {
	l, err := os.Create("log")
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(l)
	log.SetOutput(l)

	// lite Pathgenic tabix database
	if true {
		// load ClinVar
		ClinVarPathogenicLite := FindPathogenic(clinvarAnno, isPathogenic, clinvarCol, evidence.IsClinVarPLP, columns)
		sort.Sort(Bed(ClinVarPathogenicLite))
		f1, err := os.Create("ClinVarPathogenicLite.bed")
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(f1)

		_, err = fmt.Fprintln(f1, strings.Join(columns, "\t"))
		simple_util.CheckErr(err)
		for _, item := range ClinVarPathogenicLite {
			_, err = fmt.Fprintln(f1, strings.Join(item, "\t"))
			simple_util.CheckErr(err)
		}

		// load HGMD
		HgmdPathogenicLite := FindPathogenic(hgmdAnno, isPathogenic, hgmdCol, evidence.IsHgmdDM, columns)
		sort.Sort(Bed(HgmdPathogenicLite))
		f2, err := os.Create("HgmdPathogenicLite.bed")
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(f2)

		_, err = fmt.Fprintln(f2, strings.Join(columns, "\t"))
		simple_util.CheckErr(err)
		for _, item := range HgmdPathogenicLite {
			_, err = fmt.Fprintln(f2, strings.Join(item, "\t"))
			simple_util.CheckErr(err)
		}

		// merge ClinVar and HGMD
		var merge [][]string
		var dup = make(map[string]bool)
		for _, item := range ClinVarPathogenicLite {
			var key = strings.Join(item[0:5], "\t")
			if dup[key] {
				continue
			}
			merge = append(merge, item)
			dup[key] = true
		}
		for _, item := range HgmdPathogenicLite {
			var key = strings.Join(item[0:5], "\t")
			if dup[key] {
				continue
			}
			merge = append(merge, item)
			dup[key] = true
		}
		sort.Sort(Bed(merge))
		f, err := os.Create("PathogenicLite.bed")
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(f)
		_, err = fmt.Fprintln(f, strings.Join(columns, "\t"))
		simple_util.CheckErr(err)
		for _, item := range merge {
			_, err = fmt.Fprintln(f, strings.Join(item, "\t"))
			simple_util.CheckErr(err)
		}
	}

	// spec.var.list anno clinvar hgmd and filter
	if false {
		mutYzyList := simple_util.File2Array("spec.var.yzy.list")
		var yzyCover = make(map[string]bool)
		mutList := simple_util.File2Array("spec.var.list")
		var annotation = make(map[string]map[string]string)
		var annoLite = make(map[string]map[string]string)
		loadClinVar, _ := simple_util.File2MapArray(clinvarAnno, "\t", nil)
		loadHGMD, _ := simple_util.File2MapArray(hgmdAnno, "\t", nil)
		for _, item := range loadClinVar {
			key := item["MutationName"]
			annotation[key] = item
		}
		for _, item := range loadHGMD {
			key := item["MutationName"]
			info, ok := annotation[key]
			if ok {
				info["HGMD Pred"] = item[hgmdCol]
			} else {
				annotation[key] = item
			}
		}
		var columns = []string{
			"#Chr",
			"Start",
			"Stop",
			"Ref",
			"Call",
			"MutationName",
			"VarType",
			clinvarCol,
			hgmdCol,
		}
		file, err := os.Create("spec.var.list.txt")
		simple_util.CheckErr(err)
		defer file.Close()

		lite, err := os.Create("spec.var.list.lite.txt")
		simple_util.CheckErr(err)
		defer lite.Close()
		liteList, err := os.Create("spec.var.lite.list")
		simple_util.CheckErr(err)
		defer liteList.Close()

		fmt.Fprintln(file, strings.Join(columns, "\t"))
		fmt.Fprintln(lite, strings.Join(columns, "\t"))

		for _, key := range mutList {
			var item = make(map[string]string)
			var array []string
			for _, col := range columns {
				item[col] = annotation[key][col]
				array = append(array, item[col])
			}
			fmt.Fprintln(file, strings.Join(array, "\t"))
			if evidence.IsClinVarPLP.MatchString(annotation[key][clinvarCol]) || evidence.IsHgmdDM.MatchString(annotation[key][hgmdCol]) {
				fmt.Fprintln(lite, strings.Join(array, "\t"))
				fmt.Fprintln(liteList, key)
				yzyCover[key] = true
			}
			annoLite[key] = item
		}
		for _, key := range mutYzyList {
			if !yzyCover[key] {
				fmt.Printf("%s\t[%v]\n", key, annoLite[key])
				fmt.Fprintln(liteList, key)
			}
		}
		jsonByte, err := simple_util.JsonIndent(annoLite, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "spec.var.list.json")
	}

	// get gene region
	// chromosome accessions
	var chrAcce = "http://ftp.ncbi.nlm.nih.gov/genomes/H_sapiens/ARCHIVE/BUILD.37.3/Assembled_chromosomes/chr_accessions_GRCh37.p5"
	var chrAcceFile = "chr_accessions_GRCh37.p5"
	if false {
		simple_util.DownloadFile(chrAcceFile, chrAcce)

	}
	if false {
		var acce2chr = make(map[string]string)
		chrAcceMap := simple_util.File2MapMap(chrAcceFile, "RefSeq Accession.version", "\t")
		for key, item := range chrAcceMap {
			acce2chr[key] = item["#Chromosome"]
		}
		simple_util.Json2File("accession2chr.json", acce2chr)
	}

	var genomicGffUrl = "http://ftp.ncbi.nih.gov/refseq/H_sapiens/annotation/GRCh37_latest/refseq_identifiers/GRCh37_latest_genomic.gff.gz"
	var genomcGffFile = "GRCh37_latest_genomic.gff.gz"
	if false {
		simple_util.DownloadFileProgress(genomcGffFile, genomicGffUrl)
	}
	if false {
		var genomicGFF = parseGff3.File2GFF3array(genomcGffFile)
		acce2chr := simple_util.JsonFile2Map("accession2chr.json")
		var RSGregion = make(map[string][]evidence.Region)
		for _, item := range genomicGFF {
			if item.Type != "transcript" && item.Type != "mRNA" {
				continue
			}
			var region = new(evidence.Region)
			region.Seqid = item.Seqid
			region.Type = item.Type
			region.Chromosome = acce2chr[region.Seqid]
			if region.Chromosome == "" {
				continue
			}
			region.Start = item.Start
			region.End = item.End
			region.Strand = item.Strand
			region.Gene = item.Attributes["gene"]
			name := item.Attributes["Name"]
			old, ok := RSGregion[name]
			if ok {
				log.Printf("Duplicate Transcript(%s):\t%+v vs. %+v", name, old, *region)
			} else {
			}
			RSGregion[name] = append(RSGregion[name], *region)
		}
		err := simple_util.Json2File("transcript.info.json", RSGregion)
		simple_util.CheckErr(err)
	}

	// build PVS1 db
	if false {
		var LOFIntoleranceGeneList = make(map[string]int)

		// load ClinVar
		var ClinVarGeneList = evidence.FindLOFIntoleranceGeneList(clinvarAnno, clinvarCol, evidence.IsClinVarPLP)
		for key, val := range ClinVarGeneList {
			if val > 0 {
				LOFIntoleranceGeneList[key] += val
			}
		}
		jsonByte, err := json.MarshalIndent(ClinVarGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarLOFIntoleranceGeneList.json")

		// load HGMD
		var HGMDGeneList = evidence.FindLOFIntoleranceGeneList(hgmdAnno, hgmdCol, evidence.IsHgmdDM)
		for key, val := range HGMDGeneList {
			if val > 0 {
				LOFIntoleranceGeneList[key] += val
			}
		}
		jsonByte, err = json.MarshalIndent(HGMDGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDLOFIntoleranceGeneList.json")

		var exacList = simple_util.File2Array(exacLOF)
		for _, key := range exacList {
			LOFIntoleranceGeneList[key]++
		}
		jsonByte, err = json.MarshalIndent(LOFIntoleranceGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "LOFIntoleranceGeneList.json")

		// load LOFIntoleranceGeneList
		b, err := ioutil.ReadFile("LOFIntoleranceGeneList.json")
		simple_util.CheckErr(err)
		err = json.Unmarshal(b, &LOFIntoleranceGeneList)
		simple_util.CheckErr(err)
	}
	// test PVS1
	if false {
		var item = map[string]string{
			"Function":    "splice-3",
			"Gene Symbol": "ACD",
			"Start":       "67693132",
			"#Chr":        "16",
			"Transcript":  "NM_001082486.1",
		}
		LOFIntoleranceGeneList := simple_util.JsonFile2MapInt("LOFIntoleranceGeneList.json")
		var transcriptInfo map[string][]evidence.Region
		b, err := ioutil.ReadFile("transcript.info.json")
		simple_util.CheckErr(err)
		err = json.Unmarshal(b, &transcriptInfo)
		simple_util.CheckErr(err)
		fmt.Println("PVS1", evidence.CheckPVS1(item, LOFIntoleranceGeneList, transcriptInfo))
	}

	// build PS1/PM5 db
	if false {
		// load ClinVar
		var ClinVarMissense, ClinVarPHGVSList, ClinVarAAPosList = evidence.FindPathogenicMissense(clinvarAnno, clinvarCol, evidence.IsClinVarPLP)
		jsonByte, err := simple_util.JsonIndent(ClinVarMissense, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPathogenicMissense.json")
		jsonByte, err = simple_util.JsonIndent(ClinVarPHGVSList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPHGVSList.json")
		jsonByte, err = simple_util.JsonIndent(ClinVarAAPosList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarAAPosList.json")

		// load HGMD
		var HGMDMissense, HGMDPHGVSlist, HGMDAAPosList = evidence.FindPathogenicMissense(clinvarAnno, hgmdCol, evidence.IsHgmdDM)
		jsonByte, err = simple_util.JsonIndent(HGMDMissense, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPathogenicMissense.json")
		jsonByte, err = simple_util.JsonIndent(HGMDPHGVSlist, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPHGVSList.json")
		jsonByte, err = simple_util.JsonIndent(HGMDAAPosList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDAAPosList.json")
	}
	// test PS1
	if false {
		var item = map[string]string{
			"MutationName": "NM_000142.4(FGFR3): c.1138G>A (p.Gly380Arg)",
			"Transcript":   "NM_000142.4",
			"pHGVS":        "p.G380R | p.Gly380Arg",
		}
		var ClinVarMissense = simple_util.JsonFile2MapInt("ClinVarPathogenicMissense.json")
		var ClinVarPHGVSlist = simple_util.JsonFile2MapInt("ClinVarPHGVSList.json")
		var HGMDMissense = simple_util.JsonFile2MapInt("HGMDPathogenicMissense.json")
		var HGMDPHGVSlist = simple_util.JsonFile2MapInt("HGMDPHGVSList.json")
		fmt.Println("PS1", evidence.CheckPS1(item, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist))
	}
	// test PM5
	if false {
		var item = map[string]string{
			"MutationName": "NM_000016.4(ACADM): c.616C>T (p.Arg206Cys)",
			"Transcript":   "NM_000016.4",
			"pHGVS":        "p.R206C | p.Arg206Cys",
		}
		var ClinVarPHGVSlist = simple_util.JsonFile2MapInt("ClinVarPHGVSList.json")
		var ClinVarAAPosList = simple_util.JsonFile2MapInt("ClinVarAAPosList.json")
		var HGMDPHGVSlist = simple_util.JsonFile2MapInt("HGMDPHGVSList.json")
		var HGMDAAPosList = simple_util.JsonFile2MapInt("HGMDAAPosList.json")
		fmt.Println("PM5", evidence.CheckPM5(item, ClinVarPHGVSlist, ClinVarAAPosList, HGMDPHGVSlist, HGMDAAPosList))
	}

	// build PM1 db
	if false {
		var ClinVarPathogenicDomain = evidence.FindPM1MutationDomain(clinvarAnno, evidence.FilterPathogenic)
		jsonByte, err := simple_util.JsonIndent(ClinVarPathogenicDomain, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPathogenicDomain.json")

		var HGMDPathogenicDomain = evidence.FindPM1MutationDomain(hgmdAnno, evidence.FilterPathogenic)
		jsonByte, err = simple_util.JsonIndent(HGMDPathogenicDomain, "", "\t")
		simple_util.CheckErr(simple_util.Json2rawFile("HGMDPathogenicDomain.json", HGMDPathogenicDomain))

		var mutationList = make(map[string]int)
		var dbNSFPPathogenicDomain = make(map[string]int)
		var PfamPathogenicDomain = make(map[string]int)

		for mutation := range ClinVarPathogenicDomain {
			mutationList[mutation]++
		}
		for mutation := range HGMDPathogenicDomain {
			mutationList[mutation]++
		}

		simple_util.CheckErr(simple_util.Json2rawFile("PathogenicMutation.json", mutationList))

		for mutation := range mutationList {
			clinvarDomain, ok1 := ClinVarPathogenicDomain[mutation]
			hgmdDomain, ok2 := HGMDPathogenicDomain[mutation]
			if ok1 && ok2 {
				if clinvarDomain[0] == hgmdDomain[0] && clinvarDomain[1] == hgmdDomain[1] {
					for _, domain := range strings.Split(clinvarDomain[0], ";") {
						if domain != "" && domain != "." {
							dbNSFPPathogenicDomain[domain]++
						}
					}
					for _, domain := range strings.Split(clinvarDomain[1], ";") {
						if domain != "" && domain != "." {
							PfamPathogenicDomain[domain]++
						}
					}
				} else {
					log.Printf("[Conflicet Domain:%v vs. %v]\n", clinvarDomain, hgmdDomain)
				}
			} else if ok1 {
				for _, domain := range strings.Split(clinvarDomain[0], ";") {
					if domain != "" && domain != "." {
						dbNSFPPathogenicDomain[domain]++
					}
				}
				for _, domain := range strings.Split(clinvarDomain[1], ";") {
					if domain != "" && domain != "." {
						PfamPathogenicDomain[domain]++
					}
				}
			} else if ok2 {
				for _, domain := range strings.Split(hgmdDomain[0], ";") {
					if domain != "" && domain != "." {
						dbNSFPPathogenicDomain[domain]++
					}
				}
				for _, domain := range strings.Split(hgmdDomain[1], ";") {
					if domain != "" && domain != "." {
						PfamPathogenicDomain[domain]++
					}
				}
			}
		}
		simple_util.CheckErr(simple_util.Json2rawFile("dbNSFPPathogenicDomain.json", dbNSFPPathogenicDomain))
		simple_util.CheckErr(simple_util.Json2rawFile("PfamPathogenicDomain.json", PfamPathogenicDomain))
	}
	if false {
		var ClinVarBenignDomain = evidence.FindPM1MutationDomain(clinvarAnno, evidence.FilterBenign)
		jsonByte, err := simple_util.JsonIndent(ClinVarBenignDomain, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarBenignDomain.json")

		var HGMDBenignDomain = evidence.FindPM1MutationDomain(hgmdAnno, evidence.FilterBenign)
		jsonByte, err = simple_util.JsonIndent(HGMDBenignDomain, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDBenignDomain.json")

		var mutationList = make(map[string]int)
		var dbNSFPBenignDomain = make(map[string]int)
		var PfamBenignDomain = make(map[string]int)

		for mutation := range ClinVarBenignDomain {
			mutationList[mutation]++
		}
		for mutation := range HGMDBenignDomain {
			mutationList[mutation]++
		}

		simple_util.Json2File("BenignMutation.json", mutationList)
		for mutation := range mutationList {
			clinvarDomain, ok1 := ClinVarBenignDomain[mutation]
			hgmdDomain, ok2 := HGMDBenignDomain[mutation]
			if ok1 && ok2 {
				if clinvarDomain[0] == hgmdDomain[0] && clinvarDomain[1] == hgmdDomain[1] {
					for _, domain := range strings.Split(clinvarDomain[0], ";") {
						if domain != "" && domain != "." {
							dbNSFPBenignDomain[domain]++
						}
					}
					for _, domain := range strings.Split(clinvarDomain[1], ";") {
						if domain != "" && domain != "." {
							PfamBenignDomain[domain]++
						}
					}
				} else {
					log.Printf("[Conflicet Domain:%v vs. %v]\n", clinvarDomain, hgmdDomain)
				}
			} else if ok1 {
				for _, domain := range strings.Split(clinvarDomain[0], ";") {
					if domain != "" && domain != "." {
						dbNSFPBenignDomain[domain]++
					}
				}
				for _, domain := range strings.Split(clinvarDomain[1], ";") {
					if domain != "" && domain != "." {
						PfamBenignDomain[domain]++
					}
				}
			} else if ok2 {
				for _, domain := range strings.Split(hgmdDomain[0], ";") {
					if domain != "" && domain != "." {
						dbNSFPBenignDomain[domain]++
					}
				}
				for _, domain := range strings.Split(hgmdDomain[1], ";") {
					if domain != "" && domain != "." {
						PfamBenignDomain[domain]++
					}
				}
			}
		}
		simple_util.Json2rawFile("dbNSFPBenignDomain.json", dbNSFPBenignDomain)
		simple_util.Json2rawFile("PfamBenignDomain.json", PfamBenignDomain)
	}
	if false {
		var dbNSFPPathogenicDomain = make(map[string]int)
		var dbNSFPBenignDomain = make(map[string]int)
		var PfamPathogenicDomain = make(map[string]int)
		var PfamBenignDomain = make(map[string]int)
		dbNSFPPathogenicDomain = simple_util.JsonFile2MapInt("dbNSFPPathogenicDomain.json")
		dbNSFPBenignDomain = simple_util.JsonFile2MapInt("dbNSFPBenignDomain.json")
		PfamPathogenicDomain = simple_util.JsonFile2MapInt("PfamPathogenicDomain.json")
		PfamBenignDomain = simple_util.JsonFile2MapInt("PfamBenignDomain.json")

		var PM1dbNSFPDomain = make(map[string]bool)
		var PM1PfamDomain = make(map[string]bool)

		for domain, count := range dbNSFPPathogenicDomain {
			if count >= 2 {
				PM1dbNSFPDomain[domain] = true
			}
		}
		for domain, count := range dbNSFPBenignDomain {
			if count > 0 {
				PM1dbNSFPDomain[domain] = false
			}
		}

		for domain, count := range PfamPathogenicDomain {
			if count >= 2 {
				PM1PfamDomain[domain] = true
			}
		}
		for domain, count := range PfamBenignDomain {
			if count > 0 {
				PM1PfamDomain[domain] = false
			}
		}

		simple_util.CheckErr(simple_util.Json2rawFile("PM1PfamDomain.json", PM1PfamDomain))
		simple_util.CheckErr(simple_util.Json2rawFile("PM1dbNSFPDomain.json", PM1dbNSFPDomain))
	}

	// build PP2 db
	// load ClinVar
	if false {
		var ClinVarGenePathogenicMissenseRatio = evidence.CalGeneMissenseRatio(clinvarAnno, clinvarCol, evidence.IsClinVarPLP, 10)
		jsonByte, err := simple_util.JsonIndent(ClinVarGenePathogenicMissenseRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarGenePathogenicMissenseRatio.json")

		var ClinVarGeneBenignMissenseRatio = evidence.CalGeneMissenseRatio(clinvarAnno, clinvarCol, evidence.IsClinVarBLB, 0)
		jsonByte, err = simple_util.JsonIndent(ClinVarGeneBenignMissenseRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarGeneBenignMissenseRatio.json")

		var ClinVarPP2GeneList = make(map[string]float64)
		for key, val := range ClinVarGenePathogenicMissenseRatio {
			if ClinVarGeneBenignMissenseRatio[key] < 0.1 {
				ClinVarPP2GeneList[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(ClinVarPP2GeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPP2GeneList.json")
	}
	// load HGMD
	if false {
		var hgmdGenePathogenicMissenseRatio = evidence.CalGeneMissenseRatio(hgmdAnno, hgmdCol, evidence.IsHgmdDM, 10)
		jsonByte, err := simple_util.JsonIndent(hgmdGenePathogenicMissenseRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HgmdGenePathogenicMissenseRatio.json")

		var hgmdGeneBenignMissenseRatio = evidence.CalGeneMissenseRatio(hgmdAnno, hgmdCol, evidence.IsHgmdB, 0)
		jsonByte, err = simple_util.JsonIndent(hgmdGeneBenignMissenseRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HgmdGeneBenignMissenseRatio.json")

		var hgmPP2GeneList = make(map[string]float64)
		for key, val := range hgmdGenePathogenicMissenseRatio {
			if hgmdGeneBenignMissenseRatio[key] < 0.1 {
				hgmPP2GeneList[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(hgmPP2GeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HgmPP2GeneList.json")
	}

	//build BP1 db
	// load ClinVar
	if false {
		var ClinVarGenePathogenicLoFRatio = evidence.CalGeneLoFRatio(clinvarAnno, clinvarCol, evidence.IsClinVarPLP, 10)
		jsonByte, err := simple_util.JsonIndent(ClinVarGenePathogenicLoFRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarGenePathogenicLoFRatio.json")
	}
	if false {
		var HgmdGeneBenignLoFRatio = evidence.CalGeneLoFRatio(hgmdAnno, hgmdCol, evidence.IsHgmdDM, 10)
		jsonByte, err := simple_util.JsonIndent(HgmdGeneBenignLoFRatio, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HgmdGenePathogenicLoFRatio.json")
	}
}

type filterRule func(item map[string]string, key string, filter *regexp.Regexp) bool

func isPathogenic(item map[string]string, key string, filter *regexp.Regexp) bool {
	if filter.MatchString(item[key]) {
		return true
	}
	return false
}

func FindPathogenic(fileName string, filterPathgenic filterRule, key string, filter *regexp.Regexp, keyList []string) (mapArray [][]string) {
	itemArray, _ := simple_util.File2MapArray(fileName, "\t", nil)
	for _, item := range itemArray {
		var lite []string
		if !filterPathgenic(item, key, filter) {
			continue
		}
		for _, key := range keyList {
			lite = append(lite, item[key])
		}
		mapArray = append(mapArray, lite)
	}
	return
}

func chr2int(chromosome string) int {
	chr := strings.Replace(chromosome, "chr", "", -1)
	i, err := strconv.Atoi(chr)
	if err == nil {
		return i
	} else if chr == "X" {
		return 23
	} else if chr == "Y" {
		return 24
	} else if chr == "MT" {
		return 25
	}
	return 26
}

func compareChr(chr1, chr2 string) bool {
	chr1 = strings.Replace(chr1, "chr", "", 1)
	chr2 = strings.Replace(chr2, "chr", "", 1)
	if chr1 == chr2 {
		return false
	}
	chr1Int := chr2int(chr1)
	chr2Int := chr2int(chr2)
	if chr1Int == 26 && chr2Int == 26 {
		return strings.Compare(chr1, chr2) < 0
	} else {
		return chr1Int < chr2Int
	}

}

func compareIntString(a, b string) int {
	if a == b {
		return 0
	}
	i, err1 := strconv.Atoi(a)
	j, err2 := strconv.Atoi(b)
	if err1 != nil && err2 != nil {
		if i == j {
			return 0
		} else {
			return i - j
		}
	} else {
		return strings.Compare(a, b)
	}
}

type Bed [][]string

func (a Bed) Len() int      { return len(a) }
func (a Bed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Bed) Less(i, j int) bool {
	chr1 := a[i][0]
	chr2 := a[j][0]
	if chr1 != chr2 {
		return compareChr(chr1, chr2)
	}

	start1, err := strconv.Atoi(a[i][1])
	simple_util.CheckErr(err)
	start2, err := strconv.Atoi(a[j][1])
	simple_util.CheckErr(err)
	if start1 < start2 {
		return true
	} else if start1 > start2 {
		return false
	}

	stop1, err := strconv.Atoi(a[i][1])
	simple_util.CheckErr(err)
	stop2, err := strconv.Atoi(a[j][1])
	simple_util.CheckErr(err)
	if stop1 < stop2 {
		return true
	} else if stop1 > stop2 {
		return false
	}

	return false
}
