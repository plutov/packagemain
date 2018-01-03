Ethereum is an open blockchain platform that lets anyone build and use decentralized applications that run on blockchain technology. Like Bitcoin, no one controls or owns Ethereum – it is an open-source project built by many people around the world. But unlike the Bitcoin protocol, Ethereum was designed to be adaptable and flexible. It is easy to create new applications on the Ethereum platform, and with the Homestead release, it is now safe for anyone to use those applications.

Ethereum is a programmable blockchain. Rather than give users a set of pre-defined operations (e.g. bitcoin transactions), Ethereum allows users to create their own operations of any complexity they wish. In this way, it serves as a platform for many different types of decentralized blockchain applications, including but not limited to cryptocurrencies.

Whereas the Bitcoin blockchain was purely a list of transactions, Ethereum’s basic unit is the account. The Ethereum blockchain tracks the state of every account, and all state transitions on the Ethereum blockchain are transfers of value and information between accounts.

First and most importantly, you need the solc Solidity compiler.
http://solidity.readthedocs.io/en/develop/installing-solidity.html

Alright - with solc and geth devtools in place, we can start by generating a Go-version of the contract.sol file, which holds our smart contract.

Install VS Code solidity.

Install go-ethereum.

```
./abigen --sol=contract.sol --pkg=main --out=contract.go
```

We will be using the SimulatedBackend as our target blockchain for simplicity.

Resources:
https://zupzup.org/smart-contract-solidity/
https://zupzup.org/eth-smart-contracts-go/