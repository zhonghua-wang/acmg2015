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
	"path/filepath"
	"strings"
)

// os
var (
	ex, _  = os.Executable()
	exPath = filepath.Dir(ex)
	pSep   = string(os.PathSeparator)
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
	clinvarCol      = "ClinVar Significance"
	hgmdCol         = "HGMD Pred"
	domainDbNSFPCol = "Interpro_domain"
	domainPfamCol   = "pfamId"
)

func main() {

	// spec.var.list anno clinvar hgmd
	if false {
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
		file, err := os.Create("sepc.var.list.txt")
		simple_util.CheckErr(err)
		defer file.Close()

		fmt.Fprintln(file, strings.Join(columns, "\t"))

		for _, key := range mutList {
			var item = make(map[string]string)
			var array []string
			for _, col := range columns {
				item[col] = annotation[key][col]
				array = append(array, item[col])
			}
			fmt.Fprintln(file, strings.Join(array, "\t"))
			annoLite[key] = item
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
	var genomicGFF = parseGff3.File2GFF3array(genomcGffFile)
	acce2chr := simple_util.JsonFile2Map("accession2chr.json")
	var RSGregion = make(map[string][]evidence.Region)
	for _, item := range genomicGFF {
		if item.Type != "transcript" {
			continue
		}
		var region = new(evidence.Region)
		region.Seqid = item.Seqid
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
		}
		LOFIntoleranceGeneList := simple_util.JsonFile2MapInt("LOFIntoleranceGeneList.json")
		fmt.Println("PVS1", evidence.CheckPVS1(item, LOFIntoleranceGeneList))
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
	// load ClinVar
	if false {
		var ClinVarPathogenicDomainDbNSFP = evidence.FindDomain(clinvarAnno, domainDbNSFPCol, clinvarCol, evidence.IsClinVarPLP)
		jsonByte, err := simple_util.JsonIndent(ClinVarPathogenicDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPathogenicDomainDbNSFP.json")

		var ClinVarBenignDomainDbNSFP = evidence.FindDomain(clinvarAnno, domainDbNSFPCol, clinvarCol, evidence.IsClinVarBLB)
		jsonByte, err = simple_util.JsonIndent(ClinVarBenignDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarBenignDomainDbNSFP.json")

		var ClinVarDomainDbNSFP = make(map[string]int)
		for key, val := range ClinVarPathogenicDomainDbNSFP {
			if ClinVarBenignDomainDbNSFP[key] > 0 {
				continue
			}
			if val >= 2 {
				ClinVarDomainDbNSFP[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(ClinVarDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarDomainDbNSFP.json")

		var ClinVarPathogenicDomainPfam = evidence.FindDomain(clinvarAnno, domainPfamCol, clinvarCol, evidence.IsClinVarPLP)
		jsonByte, err = simple_util.JsonIndent(ClinVarPathogenicDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPathogenicDomainPfam.json")

		var ClinVarBenignDomainPfam = evidence.FindDomain(clinvarAnno, domainPfamCol, clinvarCol, evidence.IsClinVarBLB)
		jsonByte, err = simple_util.JsonIndent(ClinVarBenignDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarBenignDomainPfam.json")

		var ClinVarDomainPfam = make(map[string]int)
		for key, val := range ClinVarPathogenicDomainPfam {
			if ClinVarBenignDomainPfam[key] > 0 {
				continue
			}
			if val >= 2 {
				ClinVarDomainPfam[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(ClinVarDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarDomainPfam.json")
	}
	// load HGMD
	if false {
		var HGMDPathogenicDomainDbNSFP = evidence.FindDomain(hgmdAnno, domainDbNSFPCol, hgmdCol, evidence.IsHgmdDM)
		jsonByte, err := simple_util.JsonIndent(HGMDPathogenicDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPathogenicDomainDbNSFP.json")

		var HGMDBenignDomainDbNSFP = evidence.FindDomain(hgmdAnno, domainDbNSFPCol, clinvarCol, evidence.IsClinVarBLB)
		jsonByte, err = simple_util.JsonIndent(HGMDBenignDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDBenignDomainDbNSFP.json")

		var HGMDDomainDbNSFP = make(map[string]int)
		for key, val := range HGMDPathogenicDomainDbNSFP {
			if HGMDBenignDomainDbNSFP[key] > 0 {
				continue
			}
			if val >= 2 {
				HGMDDomainDbNSFP[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(HGMDDomainDbNSFP, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDDomainDbNSFP.json")

		var HGMDPathogenicDomainPfam = evidence.FindDomain(hgmdAnno, domainPfamCol, hgmdCol, evidence.IsHgmdDM)
		jsonByte, err = simple_util.JsonIndent(HGMDPathogenicDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPathogenicDomainPfam.json")

		var HGMDBenignDomainPfam = evidence.FindDomain(hgmdAnno, domainPfamCol, clinvarCol, evidence.IsClinVarBLB)
		jsonByte, err = simple_util.JsonIndent(HGMDBenignDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDBenignDomainPfam.json")

		var HGMDDomainPfam = make(map[string]int)
		for key, val := range HGMDPathogenicDomainPfam {
			if HGMDBenignDomainPfam[key] > 0 {
				continue
			}
			if val >= 2 {
				HGMDDomainPfam[key] = val
			}
		}
		jsonByte, err = simple_util.JsonIndent(HGMDDomainPfam, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDDomainPfam.json")
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
