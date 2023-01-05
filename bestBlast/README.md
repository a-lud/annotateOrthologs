# Get best-BLAST hits

`BLAST` is a great tool for assigning crude annotations to unknown sequences. One of the challenges of
using `BLAST` in this way is only considering high-confidence BLAST hits. The tool `bestBlast.py` is my
approach to filtering `BLAST` results for best-blast-hits.

# Usage

Below is the help page for the `bestBlast.py` program.

```
usage: bestBlast.py [-h] [-e EXTENSION] [-i IDENTITY] [-q QCOVERAGE] [-p PROP_QOFS]
                    /path/to/dir best-hits.csv

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

positional arguments:
  /path/to/dir          Path to directory that contains BLAST results
  best-hits.csv         Output file (with path) for best hits (default: best-hits.csv)

options:
  -h, --help            show this help message and exit
  -e EXTENSION, --extension EXTENSION
                        Extension of BLAST output files (default: outfmt6).
  -i IDENTITY, --identity IDENTITY
                        Minimum percentage-identity between 0..100 (default: 80).
  -q QCOVERAGE, --qcoverage QCOVERAGE
                        Minimum query-coverage relative to subject 0..100 (default: 80).
  -p PROP_QOFS, --prop_qofs PROP_QOFS
                        Minimum proportion of query-length relative to subject-length between 0..100
                        (default: 80).

Alastair J. Ludington
University of Adelaide
2022
```

This tool filters BLAST hits on three key criteria:

1. Percentage-identity of query sequence to target sequence
2. Percentage-covereage of query sequence to target sequence
3. The query and target sequence lengths are roughly proportional

The first filter above is to ensure that the query sequence shares high sequence identity to its
best target. The second filter is to ensure that query sequence fully aligns with the target sequence.
As BLAST is a local aligner, it is possible to have a short, but highly similar fragment of the query
align to the target, which is great, but does not mean the sequences are idnetical across their full
length. Therefore, the final filter is to check that the query sequence is roughly proportional in length to
the target. By default, the threshold for each of the filters in 80\%.

The program uses BLAST output in the *outfmt 6* format. However, it expects a few custom columns in the output,
so you can't just use the default *outfmt 6* output. The custom column order is specified in the help
page above. Given a directory containing BLAST files in this format, the tool will process each file
and append the results to a long-format table.

## Output

The program generates a long-format table with three fields: sample (basename of the input file),
transcript-identifier and UniProtKB accession. Example of the output is shown below.

```
python_bivittatus,rna-NC_021479.1:15192..16305,O48106
python_bivittatus,rna-NC_021479.1:2589..3555,O79546
python_bivittatus,rna-NC_021479.1:6397..7998,O79548
...
notechis_scutatus,rna-XM_026695068.1,Q3UHG7
notechis_scutatus,rna-XM_026695073.1,Q95KK4
notechis_scutatus,rna-XM_026695076.1,P63090
```