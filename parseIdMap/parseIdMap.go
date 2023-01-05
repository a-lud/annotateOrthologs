package main

// A massive help! https://marcellanz.com/post/file-read-challenge/

import (
	"fmt"
	"log"
	"os"
	"parseIdMap/dat"
	"parseIdMap/tab"
	"parseIdMap/utility"
	"path/filepath"
	"strings"

	"github.com/akamensky/argparse"
)

// Global variables
var (
	// datTypes/tabTypes taken from
	// https://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/idmapping/README
	datTypes = []string{
		"UniProtKB-ID", "Allergome", "ArachnoServer", "Araport", "BioCyc", "BioGRID", "BioMuta", "CCDS",
		"CGD", "ChEMBL", "ChiTaRS", "CLAE", "ComplexPortal", "CPTAC", "CRC64", "dictyBase",
		"DIP", "DisProt", "DMDM", "DNASU", "DrugBank", "EchoBASE", "eggNOG", "EMBL", "EMBL-CDS",
		"Ensembl", "EnsemblGenome", "EnsemblGenome_PRO", "EnsemblGenome_TRS", "Ensembl_PRO",
		"Ensembl_TRS", "ESTHER", "FlyBase", "GeneCards", "GeneID", "Gene_Name", "Gene_OrderedLocusName",
		"Gene_ORFName", "GeneReviews", "Gene_Synonym", "GeneTree", "GeneWiki", "GenomeRNAi", "GI",
		"GlyConnect", "GuidetoPHARMACOLOGY", "HGNC", "HOGENOM", "IDEAL", "KEGG", "LegioList", "Leproma",
		"MaizeGDB", "MEROPS", "MGI", "MIM", "MINT", "NCBI_TaxID", "neXtProt", "OMA", "Orphanet", "OrthoDB",
		"PATRIC", "PDB", "PeroxiBase", "PharmGKB", "PHI-base", "PlantReactome", "PomBase", "ProteomicsDB",
		"PseudoCAP", "Reactome", "RefSeq", "RefSeq_NT", "RGD", "SGD", "STRING", "SwissLipids", "TAIR", "TCDB",
		"TreeFam", "TubercuList", "UCSC", "UniParc", "UniPathway", "UniRef100", "UniRef50", "UniRef90", "VEuPathDB",
		"VGNC", "WBParaSite", "WBParaSite_TRS_PRO", "World-2DPAGE", "WormBase", "WormBase_PRO", "WormBase_TRS", "Xenbase", "ZFIN",
	}

	tabTypes = []string{
		"UniProtKB-AC", "UniProtKB-ID", "EntrezGene",
		"RefSeq", "GI", "PDB", "GO", "UniRef100", "UniRef90", "UniRef50", "UniParc",
		"PIR", "NCBI-taxon", "MIM", "UniGene", "PubMed", "EMBL", "EMBL-CDS", "Ensembl",
		"Ensembl_TRS", "Ensembl_PRO", "Additional_PubMed",
	}

	// Column index of idmapping_selected file
	tabIndex = map[string]int{
		"UniProtKB-AC":      0,
		"UniProtKB-ID":      1,
		"EntrezGene":        2,
		"RefSeq":            3,
		"GI":                4,
		"PDB":               5,
		"GO":                6,
		"UniRef100":         7,
		"UniRef90":          8,
		"UniRef50":          9,
		"UniParc":           10,
		"PIR":               11,
		"NCBI-taxon":        12,
		"MIM":               13,
		"UniGene":           15,
		"PubMed":            16,
		"EMBL":              17,
		"EMBL-CDS":          18,
		"Ensembl":           19,
		"Ensembl_TRS":       20,
		"Ensembl_PRO":       21,
		"Additional_PubMed": 22,
	}
)

func main() {

	// Arguments
	parser := argparse.NewParser(
		"parseIdMap",
		"This is a tool to parse the UniProtKB 'idmapping.dat.gz' and 'idmapping_selected.tab.gz' "+
			"files for key information when given a CSV file containing UniprotKB accessions. "+
			"As these files are cumbersome to work with due to their size, I've tried to speed "+
			"things up by utilising go-routines. Even so, it can still take a bit of time to process the data."+
			"\n\n"+
			"\t\t  "+
			"This tool is supposed to be simple, merely providing a somewhat easy method to extract "+
			"information from these files without having to rely on the web-api.",
	)

	accessions := parser.String(
		"a",
		"accession",
		&argparse.Options{
			Required: true,
			Help:     "CSV of UniProtKB accessions. Can be generated using 'bestBlast.py' script provided in this repository.",
		},
	)

	idmap := parser.String(
		"m",
		"mapFile",
		&argparse.Options{
			Required: true,
			Help:     "File path to either 'IDmapping.dat.gz' or 'IDmapping_selected.tab.gz' files",
		},
	)

	idtype := parser.StringList(
		"i",
		"idType",
		&argparse.Options{
			Required: true,
			Help:     "Which ID field/s to extract. This argument can be specified multiple times.",
		},
	)

	outCsv := parser.String(
		"o",
		"outfile",
		&argparse.Options{
			Required: false,
			Help:     "Output CSV file name (with path)",
			Default:  "IDmapping.parsed.csv",
		},
	)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// Parse best-blast-hits CSV file for UniProtKB accessions
	uniProtAcc := utility.GetUniprotIDs(*accessions)

	// Pipeline depending on input file
	basename := strings.Split(strings.ToLower(filepath.Base(*idmap)), ".")[0]

	// Parsing approach based on basename of file
	switch basename {
	case "idmapping":
		utility.CheckIdType(*idtype, datTypes)
		matches := dat.ParseDatParallel(*idmap, *idtype, uniProtAcc)
		utility.WriteSliceToCsv(basename, *idtype, matches, *outCsv)
	case "idmapping_selected":
		utility.CheckIdType(*idtype, tabTypes)
		matches := tab.ParseTabParallel(*idmap, *idtype, tabIndex, uniProtAcc)
		utility.WriteSliceToCsv(basename, *idtype, matches, *outCsv)
	}

	log.Println("Finished!")
}
