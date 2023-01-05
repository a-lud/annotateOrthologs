# Parse UniProt 'idmapping' files

UniProt-SwissProt is a reviewed database of protein sequences. The great thing about this database is that
there are two files that accompany it that contain a whole range of functional meta-data. These are publicly
availble from UniProt's [download][updl] page. The two files in question are:

- idmappping.dat.gz
- idmapping_selected.tab.gz

The **idmapping.dat.gz** file is a large, long-format file with a huge amount of meta-data for each UniProtKB
accession. The **idmapping_selected.tab.gz** file is somewhat of a subset of the first file, but in TSV
format rather than long-format. It also contains a few extra fields not found in the *dat* file.

These files are the backbone to UniProt's `ID mapping` web-portal. Rather than having to upload a list of identifiers
to the web and download the results, I wrote the `parseIdMap.go` program to do it programatically.

## Usage

The help page for the tool is shown below:

```
usage: parseIdMap [-h|--help] -a|--accession "<value>" -m|--mapFile "<value>"
                  -i|--idType "<value>" [-i|--idType "<value>" ...]
                  [-o|--outfile "<value>"]

                  This is a tool to parse the UniProtKB 'idmapping.dat.gz' and
                  'idmapping_selected.tab.gz' files for key information when
                  given a CSV file containing UniprotKB accessions. As these
                  files are cumbersome to work with due to their size, I've
                  tried to speed things up by utilising go-routines. Even so,
                  it can still take a bit of time to process the data.

		  This tool is supposed to be simple, merely providing a somewhat easy method
                  to extract information from these files without having to
                  rely on the web-api.

Arguments:

  -h  --help       Print help information
  -a  --accession  CSV of UniProtKB accessions. Can be generated using
                   'bestBlast.py' script provided in this repository.
  -m  --mapFile    File path to either 'IDmapping.dat.gz' or
                   'IDmapping_selected.tab.gz' files
  -i  --idType     Which ID field/s to extract. This argument can be specified
                   multiple times.
  -o  --outfile    Output CSV file name (with path). Default:
                   IDmapping.parsed.csv
```

The program expects a three column CSV file from `bestBlast.py`, where the third column
is the UniProtKB accessions. You can then pass either the *idmapping.dat.gz* or the
*idmapping_selected.tab.gz* file, along with the fields of interest you want to parse
from either of the files.

## Output

The output from the tool is a CSV file. The UniProtKB accession will always be the first
column, while the remaining columns are dependent on what the user asks for. An example
of outputs is shown below.

Output: 'idmapping.dat.gz'
```
accession,idtype,id
A4K2U9,Gene_Name,YWHAB
P62262,Gene_Name,YWHAE
P61983,Gene_Name,Ywhag
O70456,Gene_Name,Sfn
O77642,Gene_Name,SFN
Q5ZMD1,Gene_Name,YWHAQ
Q52M98,Gene_Name,ywhaq
Q5ZKC9,Gene_Name,YWHAZ
P63102,Gene_Name,Ywhaz
```

Output: 'idmapping_selected.tab.gz'
```
accession,GO
Q60495,GO:0030424; GO:0009986; GO:0005905; ... GO:0006417; GO:0008542
P79307,GO:0030424; GO:0009986; GO:0005905; ... GO:0006417; GO:0008542
P63116,GO:0097440; GO:0016021; GO:0043025; ... GO:0042942; GO:0015804
```

[updl]: https://www.uniprot.org/help/downloads