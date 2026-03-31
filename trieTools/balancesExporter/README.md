# Balances exporter

This tool exports account balances, as found in the trie database, under a specific roothash. The roothash is automatically selected by the tool given the provided epoch.

## How to use

Compile the code as follows:

```
cd numbat-tools-go/trieTools/balancesExporter
go build .
```

Make sure you have a node database prepared (synchronized or downloaded) in advance.

Then, run the export command for an epoch of your choice:

```
./balancesExporter --log-save --db-path=db/1 --shard=0 --epoch=690 --format=plain-json
```

Furthermore, you can configure the filters applied to the exported accounts, as follows:

```
# include accounts with zero-balance
./balancesExporter [...] --with-zero
```

```
# include accounts that are smart contracts
./balancesExporter [...] --with-contracts
```

```
# exclude accounts that do not match the provided projected shard
./balancesExporter [...] --by-projected-shard=4
```

**Note:** the *projected shard of an account* is its containing shard, given a network with the maximum number of shards (256). In other words, the projected shard is given by the last byte of the public key.


### Export formats

When running the tool, you can specify the desired export format. The available formats are: 

`plain-text`:

```
drt1... 1000000000000000000 REWA
```

`plain-json`:

```
[
    {
        "address": "drt1...",
        "balance": "1000000000000000000"
    },
...
]
```

`rosetta-json`:

```
[
    {
        "account_identifier": {
            "address": "drt1..."
        },
        "currency": {
            "symbol": "REWA",
            "decimals": 18
        },
        "value": "1000000000000000000"
    },
...
]
```
