import requests
from pathlib import Path
import pandas as pd

gnomad_lof_url = 'https://gnomad-public-us-east-1.s3.amazonaws.com/release/2.1.1/constraint/gnomad.v2.1.1.lof_metrics.by_gene.txt.bgz'
gene_info_url = 'https://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/Mammalia/Homo_sapiens.gene_info.gz'

data_dir = Path('db/raw_data')

disease_file = data_dir / 'phoenix-disease.xlsx'

if not disease_file.exists():
    raise ('please provide phoenix disease file')

gnomad_lof_file = data_dir / 'gnomad.v2.1.1.lof_metrics.by_gene.txt.bgz'
gene_info_file = data_dir / 'Homo_sapiens.gene_info.gz'
result_file = data_dir / 'gnomad_zscore.tsv'

if not gnomad_lof_file.exists():
    print('Downloading gnomad lof file...')
    r = requests.get(gnomad_lof_url)
    with open(gnomad_lof_file, 'wb') as f:
        f.write(r.content)

if not gene_info_file.exists():
    print('Downloading gene info file...')
    r = requests.get(gene_info_url)
    with open(gene_info_file, 'wb') as f:
        f.write(r.content)

disease_df = pd.read_excel(disease_file)

disease_dict = disease_df.groupby('entry ID')['Gene/Locus'].first().to_dict()

gnomad_df = pd.read_csv(gnomad_lof_file, sep='\t',
                        compression='gzip')  # type: ignore

gene_info_df = pd.read_csv(gene_info_file, sep='\t',
                           compression='gzip')  # type: ignore


def xRef2dict(value: str):
    # 确认输入值为字符串
    if not isinstance(value, str):
        raise ValueError("Input must be a string")

    # 使用字典推导式来简化代码
    return {key: ':'.join(value) for key, *value in [item.split(':') for item in value.split('|')]}


gene_info_df = gene_info_df.join(pd.json_normalize(
    gene_info_df['dbXrefs'].map(xRef2dict)))


merge_df = gnomad_df.merge(gene_info_df[['GeneID', 'HGNC', 'Ensembl']],
                           left_on='gene_id', right_on='Ensembl', how='left')

# drop empty gene id

merge_df = merge_df.dropna(subset=['GeneID'])

# save GeneID to int

merge_df['entrez_id'] = merge_df['GeneID'].astype(int)

merge_df['gene'] = merge_df.apply(
    lambda x: disease_dict.get(x['entrez_id'], x['gene']), axis=1)

merge_df[['gene', 'entrez_id', 'HGNC', 'Ensembl', 'mis_z']]\
    .rename(columns={'HGNC': 'hgnc', 'Ensembl': 'ensembl'})\
    .to_csv(result_file,
            sep='\t', index=False)
