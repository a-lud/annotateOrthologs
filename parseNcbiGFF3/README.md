# Parse NCBI GFF3 files for gene symbols

This tool was written with the sole purpose of parsing gene symbols from NCBI annotated GFF3 files.
The script utilises the [gffutils][gffutils] Python package to handle loading and parsing of the
NCBI GFF3 files. 

## Usage

The program is called using the following command:

```
python3 parseNcbiGff3.py <arguments>
```

The help page for the program is shown below.


```
usage: parseNcbiGff3.py [-h] /path/to/gff3/dir .gff3 /path/to/out.csv

# -------------------------------------------------------- #
#                       ParseNcbiGFF3                      #
# -------------------------------------------------------- #

Script that parses the gene symbol from NCBI GFF3 files. The
tool uses the DBxref field as a unique identifier as
indicated by the 'gffutils' tool.

------------------------------------------------------------

positional arguments:
  /path/to/gff3/dir  path to directory containing GFF3 files
  .gff3              GFF3 file extension to match (default: '.gff3')
  /path/to/out.csv   Name (and path) of output file (default: ncbi-gene-names.csv)

options:
  -h, --help         show this help message and exit

Alastair J. Ludington
University of Adelaide
2022
```

It takes three inputs: a directory path to where the NCBI GFF3 files are, the extension
of the GFF3 files, and the name of the output file (with directory path if you want it
put somewhere specific). The tool expects that the input GFF3 files are named according
to their sample e.g. *homo_sapiens.gff3*. If you don't change the name to something
simple/informative, that's ok, it'll just be used as the sample id in the first column
of the output file.

Non-NCBI GFF3 files can exist in the given directory, as the tool checks each file in the
directory for the expected NCBI GFF3 header information. Once the NCBI GFF3 files have been
identified, they are each loaded separately into a database using `gffutils`. The `Dbxref`
field in the 9th column of the file used to prevent sequence-id conflicts brought about
by gene isoforms. Once loaded, sequences are filtered out if they're not protein-coding, or
if they only have a locus tag e.g. *LOC...*. This is repeated for each valid GFF3, appending
each samples results to a long-form table.

## Output

The output from `parseNcbiGff3.py` is shown below.

```
pseudonaja_textilis,rna-XM_026703005.1,GATA3
pseudonaja_textilis,rna-XM_026707001.1,CELF2
pseudonaja_textilis,rna-XM_026707264.1,CELF2
pseudonaja_textilis,rna-XM_026707708.1,USP6NL
pseudonaja_textilis,rna-XM_026707832.1,ECHDC3
...
crotalus_tigris,rna-XM_039321002.1,SPATA1
crotalus_tigris,rna-XM_039321004.1,GNG5
crotalus_tigris,rna-XM_039321017.1,RPF1
crotalus_tigris,rna-XM_039321037.1,SAMD13
```

It returns a long-form table with the fields: sample name (taken from input file), transcript
identifier and gene symbol.

[gffutils]: https://daler.github.io/gffutils/index.html