package main

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(133700000)}
	sim := backends.NewSimulatedBackend(alloc)

	// deploy contract
	addr, _, contract, err := DeployTwoD(auth, sim)
	if err != nil {
		log.Fatalf("could not deploy contract: %v", err)
	}
	log.Printf("contract deployed to %s\n", addr.String())

	// interact with contract
	
	value, err := contract.GetValue(nil, 1, 1)
	if err != nil {
		log.Printf("could not get value: %v", err)
	}
	log.Printf("pre-mining value of [1][1]: %d\n", value)

	log.Println("Mining...")
	// simulate mining
	sim.Commit()

	value, err = contract.GetValue(nil, 1, 1)
	if err != nil {
		log.Fatalf("could not get value: %v", err)
	}
	log.Printf("post-mining value of [1][1]: %d\n", value)

	// instantiate deployed contract
	log.Printf("Instantiating contract at address %s...\n", auth.From.String())
	instContract, err := NewTwoD(addr, sim)
	if err != nil {
		log.Fatalf("could not instantiate contract: %v", err)
	}

	value, err = instContract.GetValue(nil, 2, 2)
	if err != nil {
		log.Fatalf("could not get value: %v", err)
	}
	log.Printf("instantiated contract: value of [2][2]: %d\n", value)
}
