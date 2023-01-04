#!/usr/bin/env python3

import argparse
import logging
from os import DirEntry, scandir, path, remove
from textwrap import dedent

from pandas import DataFrame, concat, read_csv

# import pandas as pd

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s %(levelname)-8s %(message)s"
)


def getArgs():
    """Get user arguments and set up parser"""
    desc = """\
    # -------------------------------------------------------- #
    #                        BestBlast                         #
    # -------------------------------------------------------- #

    Given a directory of BLAST outputs (custom outfmt 6), return
    the "best" hit for each query. The "best" hit is subjective,
    depending on what the user requires. I have implemented a
    method that requires "best" hits to meet the following
    requirements:
        - \%-identity of query to target > threshold
        - \%-query-coverage to target > threshold
        - query-length ~ target-length

    Run BLAST with the following argument: 
        '-outfmt6 6 qaccver,saccver,qlen,slen,length,qcovs,
        pident,mismatch,gapopen,qstart,qend,sstart,send,
        evalue,bitscore"

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
        "blast",
        help="Path to directory that contains BLAST results",
        metavar="/path/to/dir",
    )

    parser.add_argument(
        "outfile",
        help="Output file (with path) for best hits (default: %(default)s)",
        metavar="best-hits.csv",
        default="best-hits.csv",
    )

    # Optional arguments
    parser.add_argument(
        "-e",
        "--extension",
        help="Extension of BLAST output files (default: %(default)s).",
        default="outfmt6",
    )

    parser.add_argument(
        "-i",
        "--identity",
        help="Minimum percentage-identity between 0..100 (default: %(default)s).",
        default=80,
        type=int,
    )

    parser.add_argument(
        "-q",
        "--qcoverage",
        help="Minimum query-coverage relative to subject 0..100 (default: %(default)s).",
        default=80,
        type=int,
    )

    parser.add_argument(
        "-p",
        "--prop_qofs",
        help="Minimum proportion of query-length relative to subject-length between 0..100 (default: %(default)s).",
        default=80,
        type=int,
    )

    args = parser.parse_args()
    return args


def listFiles(path: str, extension: str) -> list[DirEntry]:
    """List files in a directory matching a file extension"""
    logging.info("[listFiles]\tGetting GFF3 files")
    files = scandir(path)
    return [f for f in files if f.is_file() and f.name.endswith(extension)]


def getBestHit(files: list[DirEntry], pid: int, qcov: int, qofs: int) -> DataFrame:
    """Return the "best" BLAST hits from the user provided BLAST-outfmt6 tables"""

    logging.info("[getBestHit]\tGetting best-blast-hits")

    dfs = []
    for file in files:
        name = file.name.split(".")[0]

        # Import BLAST table
        df = read_csv(
            filepath_or_buffer=file.path,
            sep="\t",
            names=[
                "qaccver",
                "saccver",
                "qlen",
                "slen",
                "length",
                "qcovs",
                "pident",
                "mismatch",
                "gapopen",
                "qstart",
                "qend",
                "sstart",
                "send",
                "evalue",
                "bitscore",
            ],
        )

        # Get best hits for each gene. Requires the following:
        #   - %ID >= pid%
        #   - Query coverage >= qcov%
        #   - Proportion of query length to subject length is >= qofs%
        #   - Take top hit if these conditions are met (typically highest bit score)
        df_filt = (
            df.pipe(lambda x: x.loc[x.pident >= pid])
            .pipe(lambda x: x.loc[x.qcovs >= qcov])
            .assign(
                qlen_prop_slen=lambda x: (x.qlen / x.slen) * 100,
            )
            .pipe(lambda x: x.loc[x.qlen_prop_slen >= qofs])
            .groupby("qaccver")
            .first()
            .pipe(lambda x: x[["saccver"]])
            .pipe(lambda x: x.reset_index(level=0))
        )
        df_filt.insert(loc=0, column="species", value=name)
        logging.info(f"\t\t\t- Processed {name} - {len(df_filt.index)} best-hits")
        dfs.append(df_filt)

    return concat(dfs, axis=0, ignore_index=True)


if __name__ == "__main__":
    args = getArgs()

    blast = listFiles(args.blast, args.extension)
    df = getBestHit(
        files=blast, pid=args.identity, qcov=args.qcoverage, qofs=args.prop_qofs
    )

    if path.exists(args.outfile):
        logging.info(
            f"[Main]\t- Removing existing output: {args.outfile}"
        )
        remove(args.outfile)
    
    # Write results to file
    df.to_csv(path_or_buf=args.outfile, index=False)

