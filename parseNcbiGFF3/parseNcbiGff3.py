#!/usr/bin/env python3

import argparse
import logging
from os import DirEntry, path, remove, scandir
from textwrap import dedent

from gffutils import create_db

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s %(levelname)-8s %(message)s"
)


def getArgs():
    """Get user arguments and set up parser"""
    desc = """\
    # -------------------------------------------------------- #
    #                       ParseNcbiGFF3                      #
    # -------------------------------------------------------- #

    Script that parses the gene symbol from NCBI GFF3 files. The
    tool uses the DBxref field as a unique identifier as
    indicated by the 'gffutils' tool.

    ------------------------------------------------------------
    """

    epi = """\
    Alastair J. Ludington
    University of Adelaide
    2022
    """

    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description=dedent(desc),
        epilog=dedent(epi),
    )

    # Required, positional input file arguments
    parser.add_argument(
        "gff",
        help="path to directory containing GFF3 files",
        metavar="/path/to/gff3/dir",
    )

    parser.add_argument(
        "extension",
        help="GFF3 file extension to match (default: '%(default)s')",
        metavar=".gff3",
        default=".gff3",
    )

    parser.add_argument(
        "outfile",
        help="Name (and path) of output file (default: %(default)s)",
        metavar="/path/to/out.csv",
        default="ncbi-gene-names.csv",
    )

    args = parser.parse_args()
    return args


def listFiles(path: str, extension: str) -> list[DirEntry]:
    """List files in a directory matching a file extension"""
    logging.info("[listFiles]\t\tGetting GFF3 files")
    files = scandir(path)
    return [f for f in files if f.is_file() and f.name.endswith(extension)]


def isNCBI(files: list[DirEntry]) -> list[DirEntry]:
    """Determines if a GFF3 is annotated by NCBI by checking the third line of GFF3 header"""
    logging.info("[isNCBI]\t\tFiltering for NCBI GFF3 files")
    ncbi: list[DirEntry] = []
    for gff in files:
        with open(gff.path, "r") as f:
            for i, line in enumerate(f):
                if i < 2:
                    continue

                if i > 2:
                    break

                if "NCBI" in line:
                    ncbi.append(gff)

    return ncbi


def getGeneSymbol(gffs: list[DirEntry], outfile: str):
    """For each gene in an NCBI GFF3 file, get the gene symbol (if present). Returns a two-column CSV file: gene ID, gene symbol"""
    logging.info("[getGeneSymbol]\tGetting gene symbols for each gene model")

    # Check if output file exists, remove if it does
    if path.exists(outfile):
        logging.info(
            "\t\t\t\t- Removing existing 'ncbi-gene-names.csv' file in output directory"
        )
        remove(outfile)

    for gff in gffs:
        name = gff.name.split(".")[0]

        # Load GFF3 into database
        logging.info(f"\n\t\t\t\t\t\t\t\t- Loading {name} into database")
        db = create_db(gff.path, ":memory:", id_spec={"gene": "Dbxref"})
        gene = db.features_of_type("mRNA")

        csv = []
        for g in gene:
            if "gene" not in g.attributes.keys():
                continue
            if g["gene"][0].startswith("LOC1"):
                continue

            geneID = g["ID"][0]
            geneSymbol = g["gene"][0]

            # Append CSV string to list
            csv.append(f"{name},{geneID},{geneSymbol}")

        # Write list to file
        logging.info(f"\t\t\t\t- Writing results for {name}")
        with open(outfile, "a") as out:
            for line in csv:
                out.write(f"{line}\n")


if __name__ == "__main__":
    # Get arguments
    args = getArgs()

    # Get GFF3 files
    gffs = listFiles(args.gff, args.extension)

    # Filter for NCBI files
    ncbi = isNCBI(gffs)

    # Create CSV: sample,transcriptID,geneSymbol
    getGeneSymbol(ncbi, args.outfile)
