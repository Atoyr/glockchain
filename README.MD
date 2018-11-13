
	  /$$$$$$  /$$                     /$$  
	 /$$__  $$| $$                    | $$  
	| $$  \__/| $$  /$$$$$$   /$$$$$$$| $$   /$$
	| $$ /$$$$| $$ /$$__  $$ /$$_____/| $$  /$$/
	| $$|_  $$| $$| $$  \ $$| $$      | $$$$$$/ 
	| $$  \ $$| $$| $$  | $$| $$      | $$_  $$ 
	|  $$$$$$/| $$|  $$$$$$/|  $$$$$$$| $$ \  $$
	 \______/ |__/ \\______/  \_______/|__/  \__/
	
	      /$$$$$$  /$$                 /$$      
	     /$$__  $$| $$                |__/
	    | $$  \__/| $$$$$$$   /$$$$$$  /$$ /$$$$$$$
	    | $$      | $$__  $$ |____  $$| $$| $$__  $$
	    | $$      | $$  \ $$  /$$$$$$$| $$| $$  \ $$
	    | $$    $$| $$  | $$ /$$__  $$| $$| $$  | $$
	    |  $$$$$$/| $$  | $$|  $$$$$$$| $$| $$  | $$
	     \______/ |__/  |__/ \_______/|__/|__/  |__/

# GlockChain
This is golang blockchain application.
Not implement server function.

# Usage
## Install
Make sure you have a working Go environment.
To install it.
```
$ go get github.com/atoyr/glockchain
```

## Initialize
Generate private/public key pair and write into "wallet.dat".
Initialize blockchain with wallet key.
To initialize block, create coinbase transaction and genesis block.
Coinbase transaction value is 100.
Show your wallet address.

```
$ go run ./main.go initialize
```

## Create Wallet
Need other wallet key pair, execute command.
```
$ go run ./main.go wallet create
```

## Show Wallets
You show all wallets address.
```
$ go run ./main.go wallet list
```

## Get Balance
Show wallets balance.
```
$ go run ./main.go wallet balance
```

If you show one wallet balance, you execute it.
```
$ go run ./main.go wallet balance -a ADDRESS
```

## Send value
Send amount from wallet.
```
$ go run ./main.go transaction create -f FROM_ADDRESS -t TO_ADDRESS -am AMOUNT
```
CreatedTransaction is pooling.

## Show Transaction pool
Show transaciton pool
```
$ go run ./main.go transaciton list
```

## Mining
If pool transaction,you can execute mining.
```
$ go run ./main.go mining
```
Not implement mining incentive.

## Show Blockchain
Show blockchain.
```
$ go run ./main.go blockchain print
```


# Using library
## [Bolt DB](https://github.com/boltdb/bolt)
> Bolt DB
> The MIT License (MIT)
> Copyright (c) 2013 Ben Johnson

* ** Use It :** blockchain KeyValueDatabase
* ** License :**  The MIT License (MIT)
* ** License full text :** [licenses/BoltDB.txt](licenses/BoltDB.txt)

## [cli](https://github.com/urfave/cli)
> urfave/cli
> The MIT License (MIT)
> Copyright（c）2016 Jeremy Saenz＆Contributors

* ** Use It :** cli interface
* ** License :**  The MIT License (MIT)
* ** License full text :** [licenses/urfave_cli.txt](licenses/urfave_cli.txt)


# Refer
https://jeiwan.cc/posts/building-blockchain-in-go-part-1/

