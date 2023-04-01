## Provides the filecoin wallet sign tool

[![License](https://img.shields.io/badge/license-MIT-blue)](https://raw.githubusercontent.com/ipdr/ipdr/master/LICENSE)
[![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)](https://github.com/emersion/stability-badges#stable)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

###
This tool can sign locally without importing the wallet to a remote node, preventing private key disclosure.
The following operations require access to the lotus-RPC endpoint ( default:https://api.node.glif.io/rpc/v1 ). you can configure environment variables to connect your endpoint.

```bash
export LOTUS_API=http://127.0.0.1:1234/rpc/v1
export LOTUS_API_TOKEN=eyJhbGcI.......BSiGNLrVVbdlDs
```

### checklist
- [x] send   
- [x] withdraw
- [x] change owner
- [x] sign and verify any string message
- [x] change the miner's beneficiary
- [ ] change worker
- multisig
  - [x] propose transfer
  - [x] propose withdraw
  - [x] inspect
  - [x] approve
  - [x] propose change owner
  - [x] propose change the miner's beneficiary
  - [ ] propose change worker

### usage:

```bash
$ ./filwallet-sign --help
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Filecoin wallet sign tools

Usage:
  wallet [command]

Available Commands:
  help                          Help about any command
  msig                          multisig address tools
  send                          transfer FIL
  sign                          sign any string message
  verify                        verify the signature of any string message
  withdraw                      withdraw from miner
  change-owner                  change owner
  change-beneficiary            change beneficiary
  confirm-change-beneficiary    confirm change beneficiary
```

#### transfer FIL

+ the private key belongs to sender address

```bash
$ ./filwallet-sign send t142e....4zfa 1                                             
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
send from t1og......exb7i to t142e....4zfa amount 1
...
message CID: bafy2bzacec55......a3rjg6dtyu  
```

#### withdraw from miner

+ the private key belongs to miner's owner

```bash
$ ./filwallet-sign withdraw t01234 6.6                             
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzacebe......a3rjg6dtyu

```

#### change miner's owner

+ first the private key belongs to old owner address
+ second the private key belongs to new owner address

```bash
$ ./filwallet-sign change-owner  t03..3 t01234 t03..4 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

#### change miner's beneficiary

+ the private key belongs to owner address

```bash
$ ./filwallet-sign change-beneficiary <beneficiaryAddress> <quota> <expiration> --actor <minerAddress> [--overwrite-pending-change]
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

#### confirm change miner's beneficiary

+ the private key belongs to new or existing beneficiary (depends on miner info)

```bash
$ ./filwallet-sign confirm-change-beneficiary <minerAddress> --new-beneficiary / --existing-beneficiary
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

### multisig

> Notice: private key must be signer one of  multisigAddress

### usage 

```bash
$ ./filwallet-sign msig               
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
multisig address tools

Usage:
  wallet msig [command]

Available Commands:
  approve     approve transaction of multisigAddress
  inspect     inspect multisigAddress 
  propose     make a proposal
Flags:
  -h, --help   help for msig
```

#### approve

```bash
$ ./filwallet-sign msig approve t03..3 6
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```

#### inspect

```bash
$ ./filwallet-sign msig inspect t03..3
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
#### propose

+ transfer FIL

```bash
 ./filwallet-signmsig propose t03..3 t3v.....marqq 1.2   
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
send from t3v.....marqq to t03..3 amount 1.2 
```

+ withdraw from miner 

```bash
$ ./filwallet-sign msig propose withdraw t03..3 t01234 99999 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
withdraw 99999 FIL from t01234 
```

+ change miner's owner

```bash
$ ./filwallet-sign msig propose change-owner  t03..3 t01234 t03..4 
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
change miner t01020  owner is t03..4
```

+ change miner's beneficiary

```bash
$ ./filwallet-sign msig propose change-beneficiary <beneficiaryAddress> <quota> <expiration> --msig-addr <msigAddress> --miner-addr <minerAddress> [--overwrite-pending-change]
LOTUS_API :  http://127.0.0.1:1234/rpc/v1
LOTUS_API_TOKEN :  Bearer eyJhbGcI.......BSiGNLrVVbdlDs
Please enter the private key: 7b225.......673d227d
...
message CID: bafy2bzaceah.....i4d5qkvs
```