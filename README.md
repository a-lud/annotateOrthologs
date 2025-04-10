# Annotate Orthologs

This directory contains tools used to annotate orthogroups/orthologs. There are four tools in this
directory that are designed to help assign gene symbols and GO Terms to orthologs that simply
have OG-Identifiers.

There are many other approaches you could use to do this, but this is how I tackled the problem. Below
I give a breif rundown of the tools in each sub-directory, however more detailed information about
each tool can be found in the *README* for each individual tool.

## ParseNcbiGFF3

NCBI annotations have a range of annotation information in the 9th column of the GFF3 file. The directory
[parseNcbiGFF3][parseNcbi] contains a python script designed to parse the gene symbols for protein coding
genes from the *mRNA* transcripts.

## ParseFunannotate

`Funannotate` is a gene prediction pipeline that generates a range of useful outputs during it's runtime.
One of these is an 'annotations.txt' file that contains a bunch of functional annotation information. The
directory [parseFunannotate][parseFun] contains a little tool that quickly parses the 'annotation.txt' files
generated by `Funannotate`, returning *CSV* files with key information.

## BestBlast

The UniProt-SwissProt protein database is full of high-quality, curated protein sequences with all sorts of
additional functional information associated to them. The tool in [bestBlast][bestBlast] takes a set of `BLAST`
output files in a custom out-format 6 and returns high-confidence hits (best-BLAST-hits).

## ParseIdMap

The tool `parseIdMap` in the [parseIdMap][parseId] directory is a complementary tool to the `bestBlast.py` tool
above. The idea of this tool is to parse the meta-data files (*idmapping{_selected}.{tab,dat}.gz*) provided
by UniProt to get annotation information for our best-BLAST-hits.

[parseNcbi]: https://github.com/a-lud/annotateOrthologs/tree/main/parseNcbiGFF3
[bestBlast]: https://github.com/a-lud/annotateOrthologs/tree/main/bestBlast
[parseFun]: https://github.com/a-lud/annotateOrthologs/tree/main/parseFunannotate
[parseId]: https://github.com/a-lud/annotateOrthologs/tree/main/parseIdMap