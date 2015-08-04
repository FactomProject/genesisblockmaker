
Factom Genesis Block
=============

Important files:
- [**genesis.csv**](https://github.com/FactomProject/genesisblockmaker/blob/master/genesis.csv) - This is a list of data scraped from the bitcoin blockchain.  It will constitute the Software Sale fraction of the initial Factoid distribution.
- [**genesis_block_generator.py**](https://github.com/FactomProject/genesisblockmaker/blob/master/genesis_block_generator.py) - This is a script which downloads data from the Blockchain.info webpage and saves a list of the purchases.  
- [**scanblocks.py**](https://github.com/FactomProject/genesisblockmaker/blob/master/scanblocks.py) - This is a script which connects to a local btcd full node.  It finds transactions sent to the Factom Multisig address and saves a list of the purchases.


Note: Some of the transactions did not specify a pubkey and some specified the wrong pubkey.  Koinify is providing the correct ed25519 pubkeys which will be updated in genesis.csv.


# How to build genesis block

Go into /GenesisCSVParser . Compile the code, run it. The genesis block will be printed into Genesis.txt . The test genesis block (with extra balances taken from testing.csv) will be printed into TestGenesis.txt .