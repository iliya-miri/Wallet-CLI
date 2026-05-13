package main

import (
	//"bufio"
	"context"
	"encoding/hex"
	"fmt"

	//"go/types"
	"math/big"

	//"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/eth/gasprice"
	"github.com/ethereum/go-ethereum/ethclient"

	//"go.mongodb.org/mongo-driver/v2/mongo/address"

	"github.com/gin-gonic/gin"
	//"github.com/go-delve/delve/pkg/dwarf/reader"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofika/bip39"
)

func main() {

	_ = bip39.NewMnemonic
	_ = common.HexToAddress
	_ = gin.Default

	var Choice int

	for Choice != 4 {
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println("\x1b[31m             GO WALLET - ETH               \x1b[0m")
		fmt.Println("\x1b[31m═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n What do you want to do?\n 1. Creatw a Wallet \n 2. Get Balance \n 3. Send Money \n 4. Exit")
		fmt.Print("Your Choice : ")
		fmt.Scan(&Choice)

		switch {
		case Choice == 1:
			craeteWallet()
		case Choice == 2:
			getBalanceCLI()
		case Choice == 3:
			sendTransactionCLI()
		}

	}

}

func craeteWallet() {

	var Choice int

	for Choice != 2 {
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n What do you want to do?\n 1. Creatw a Wallet \n 2. Back to main menu")
		fmt.Print("Your Choice : ")
		fmt.Scan(&Choice)

		if Choice == 1 {
			println("\n")
			entropy, _ := bip39.NewMnemonic()
			mnemonic, _ := entropy.GenerateMnemonic()
			seed := bip39.NewSeed(mnemonic)
			privateKey, _ := crypto.ToECDSA(seed[:32])
			address := crypto.PubkeyToAddress(privateKey.PublicKey)
			privateKeyHeX := hex.EncodeToString(privateKey.D.Bytes())
			fmt.Println("\x1b[32m Your Wallet is Ready!\x1b[0m \nSave all of them please\n")
			fmt.Println("\x1b[33m 12  Word : \x1b[0m ", mnemonic)
			fmt.Println("\x1b[33m Address : \x1b[0m ", address.Hex())
			fmt.Println("\x1b[33m Private Key : \x1b[0m", privateKeyHeX)
		}

	}
}

func getBalanceCLI() {
	var Choice string

	for Choice != "q" {
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n write your wallet address : \n *For back to the main menu write q* \n")
		fmt.Print("Your address : ")
		fmt.Scan(&Choice)

		if Choice != "q" {
			var addressInput string
			addressInput = Choice
			addressInput = strings.TrimSpace(addressInput)

			client, err := ethclient.Dial("http://127.0.0.1:8545")
			if err != nil {
				fmt.Println("Error in connect to ganache : ", err)
				return
			}
			defer client.Close()

			address := common.HexToAddress(addressInput)

			balanceWei, err := client.BalanceAt(context.Background(), address, nil)
			if err != nil {
				fmt.Println("Error in get balance : ", err)
				return
			}

			balanceEth := new(big.Float).SetInt(balanceWei)

			weiInEther := new(big.Float).SetFloat64(1000000000000000000)
			balanceEth.Quo(balanceEth, weiInEther)
			fmt.Printf("\x1b[32mYour balance for %s: %6f ETH\n \x1b[0m", addressInput, balanceEth)
		}

	}
}

func sendTransactionCLI() {
	var fromAddress, toAddress, privateKey string
	var amountEther float64

	fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
	fmt.Println(" \nFrom address : \n")
	fmt.Scanln()
	fmt.Scanln(&fromAddress)

	fmt.Println("To address : \n")
	fmt.Scanln(&toAddress)

	fmt.Println("\nAmount Ether \n")
	fmt.Scanln(&amountEther)

	fmt.Println("\nPrivate Key wallet: \n")
	fmt.Scanln(&privateKey)

	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}

	weiPerEther := new(big.Float).SetFloat64(1000000000000000000)
	amountWeiBig := new(big.Float).Mul(new(big.Float).SetFloat64(amountEther), weiPerEther)
	amountWei, _ := amountWeiBig.Int(nil)

	//fmt.Println("Connecting to Ganache on http://127.0.0.1:8545")
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println("Error in connect to ganache : ", err)
		return
	}
	defer client.Close()

	fromAddressCommon := common.HexToAddress(fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddressCommon)
	if err != nil {
		fmt.Println("error in get nonce : ", err)
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("error in get gas price : ", err)
		return
	}

	gasLimit := uint64(21000)

	toAddressCommon := common.HexToAddress(toAddress)
	tx := types.NewTransaction(
		nonce,
		toAddressCommon,
		amountWei,
		gasLimit,
		gasPrice,
		nil,
	)

	// chainId, err := client.NetworkID(context.Background())
	// fmt.Println("Chain ID:", chainId)
	// if err != nil {
	// 	fmt.Println("error in get chain ID : ", err)
	// 	return
	// }

	chainId := big.NewInt(1337)
	fmt.Println("Chain ID (manual):", chainId)

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("error in convert private key : ", err)
		return
	}

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKeyECDSA)
	if err != nil {
		fmt.Println("error in sign transaction : ", err)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("error in send transaction : ", err)
		return
	}

	fmt.Println("Done.\n")
	fmt.Printf("Transaction Hash : %s \n", signedTx.Hash().Hex())

}
