# github.com/labs3/filecoin-wallet-signing

## Provides the Filecoin wallet  tool

The following operations require access to the lotus-RPC endpoint ( default:https://api.node.glif.io/rpc/v1 ). you can configure environment variables to connect your endpoint

```bash
export LOTUS_API=http://127.0.0.1:1234/rpc/v1
export LOTUS_API_TOKEN=eyJhbGcI.......BSiGNLrVVbdlDs
```

### Checklist
- [x] send   
- [x] withdraw
- [x] change owner
- [ ] change worker
- multisig
  - [x] propose
  - [x] propose withdraw
  - [x] inspect
  - [x] approve
  - [x] change owner
  - [ ] change worker

usage

```bash
$ go run main.go   
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Filecoin wallet tool

Usage:
  wallet [command]

Available Commands:
  help        Help about any command
  msig        multisig address tool
  send        send
  withdraw    withdraw from miner

```

### send to address 

+ sender is address of private key

```bash
$ go run main.go send t142e....4zfa 1                                             
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
send from t1og......exb7i to t142e....4zfa amount 1
...
message CID: bafy2bzacec55......a3rjg6dtyu  
```

### withdraw from miner

+ address of private key  must be miner's owner

```bash
$ go run main.go withdraw t01234 6.6                             
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzacebe......a3rjg6dtyu

```

+ change miner's owner

```bash js
$ go run main.go  change-owner  t03..3 t01234 t03..4 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

### multisig

> Notice: private key must be signer one of  multisigAddress

usage 

```bash
$ go run main.go msig               
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
multisig address tool

Usage:
  wallet msig [command]

Available Commands:
  approve     approve  transaction of multisigAddress
  inspect     inspect multisigAddress 
  propose     make a proposal
Flags:
  -h, --help   help for msig

```

#### propose

```bash
 go run main.go msig propose t03..3 t3v.....marqq 1.2   
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
send from t3v.....marqq to t03..3 amount 1.2 
```

#### inspect

```bash
$ go run main.go msig inspect t03..3
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Address: t03..3, ID: t03..3
Number of signatories 3 threshold  2 
t03..4 : t1abbhj....s7exb74 
t03..5 : t1abbhj....s7exb75 
t03..6 : t1abbhj....s7exb76 
Pending transaction: 
pending id: 6 , to : t3v.....marqq , method: 0 , amount: 1.2 FIL, Params: , approved [t03..3], ps: send out  

```

#### approve

```bash
$ go run main.go msig approve t03..3 6
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

#### withdraw

+ withdraw from miner 

```bash js
$ go run main.go msig propose withdraw t03..3 t01234 99999 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
withdraw 99999 FIL from t01234 
```

+ change miner's owner

```bash js
$ go run main.go msig propose change-owner  t03..3 t01234 t03..4 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
change miner t01020  owner is t03..4
```