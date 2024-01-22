# acmg2015
acmg2015 pred

## PP2 Missense z-score config

- run script `util/prep_pp2_file.py` to download required PP2 Missense z-score file from GnomAD V2.1 and gene info file from NCBI ftp, then merge the two file into `gnomad_zscore.tsv`
- add new PP2 config term `PP2MissenseZScore acmg/gnomad_zscore.tsv` to acmg cfg in `anno2xlsx` (`etc/acmg.db.list.txt`)
