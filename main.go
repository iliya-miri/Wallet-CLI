package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofika/bip39"
)

// ================================ MAIN ================================

func main() {
	var Choice int

	// Main loop - runs until user chooses exit (4)
	for Choice != 4 {
		// Print main menu
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println("\x1b[31m             GO WALLET - ETH               \x1b[0m")
		fmt.Println("\x1b[31m═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n What do you want to do?\n 1. Creatw a Wallet \n 2. Get Balance \n 3. Send Money \n 4. Exit")
		fmt.Print("Your Choice : ")
		fmt.Scan(&Choice)

		// Call function based on user choice
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

// ================================ CREATE WALLET ================================

func craeteWallet() {
	var Choice int

	// Loop until user goes back to main menu
	for Choice != 2 {
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n Which one do you perfer?\n 1. Creatw a Wallet \n 2. Back to main menu")
		fmt.Print("Your Choice : ")
		fmt.Scan(&Choice)

		if Choice == 1 {
			println("\n")

			// Generate 12-word mnemonic phrase
			entropy, _ := bip39.NewMnemonic()
			mnemonic, _ := entropy.GenerateMnemonic()
			// Convert mnemonic to seed
			seed := bip39.NewSeed(mnemonic)
			// Create private key from first 32 bytes of seed
			privateKey, _ := crypto.ToECDSA(seed[:32])
			// Derive public address from private key
			address := crypto.PubkeyToAddress(privateKey.PublicKey)
			// Convert private key to hex string
			privateKeyHeX := hex.EncodeToString(privateKey.D.Bytes())

			// Display wallet information to user
			fmt.Println("\x1b[32m Your Wallet is Ready!\x1b[0m \nSave all of them please\n")
			fmt.Println("\x1b[33m 12  Word : \x1b[0m ", mnemonic)
			fmt.Println("\x1b[33m Address : \x1b[0m ", address.Hex())
			fmt.Println("\x1b[33m Private Key : \x1b[0m", privateKeyHeX)
		}
	}
}

// ================================ GET BALANCE ================================

func getBalanceCLI() {
	var Choice string

	// Loop until user types "quit"
	for Choice != "quit" {
		fmt.Println("\x1b[31m\n═══════════════════════════════════════════\x1b[0m")
		fmt.Println(" \n write your wallet address : \n *For back to the main menu write quit* \n")
		fmt.Print("Your address : ")
		fmt.Scan(&Choice)

		if Choice != "quit" {
			var addressInput string
			addressInput = Choice
			addressInput = strings.TrimSpace(addressInput)

			// Connect to Ganache local blockchain
			client, err := ethclient.Dial("http://127.0.0.1:8545")
			if err != nil {
				fmt.Println("Error in connect to ganache : ", err)
				return
			}
			defer client.Close()

			// Convert string address to common.Address type
			address := common.HexToAddress(addressInput)

			// Get balance in Wei (smallest unit)
			balanceWei, err := client.BalanceAt(context.Background(), address, nil)
			if err != nil {
				fmt.Println("Error in get balance : ", err)
				return
			}

			// Convert Wei to Ether (1 ETH = 10^18 Wei)
			balanceEth := new(big.Float).SetInt(balanceWei)
			weiInEther := new(big.Float).SetFloat64(1000000000000000000)
			balanceEth.Quo(balanceEth, weiInEther)

			fmt.Printf("\x1b[32mYour balance for %s: %6f ETH\n \x1b[0m", addressInput, balanceEth)
		}
	}
}

// ================================ SEND TRANSACTION ================================

func sendTransactionCLI() {
	var fromAddress, toAddress, privateKey string
	var amountEther float64

	// Get transaction details from user
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

	// Remove "0x" prefix if present (HexToECDSA requires raw hex)
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = privateKey[2:]
	}

	// Convert Ether amount to Wei
	weiPerEther := new(big.Float).SetFloat64(1000000000000000000)
	amountWeiBig := new(big.Float).Mul(new(big.Float).SetFloat64(amountEther), weiPerEther)
	amountWei, _ := amountWeiBig.Int(nil)

	// Connect to Ganache
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println("Error in connect to ganache : ", err)
		return
	}
	defer client.Close()

	// Get nonce (transaction count for this address)
	fromAddressCommon := common.HexToAddress(fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddressCommon)
	if err != nil {
		fmt.Println("error in get nonce : ", err)
		return
	}

	// Get suggested gas price from network
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("error in get gas price : ", err)
		return
	}

	// Set gas limit for simple ETH transfer
	gasLimit := uint64(21000)

	// Create unsigned transaction
	toAddressCommon := common.HexToAddress(toAddress)
	tx := types.NewTransaction(
		nonce, toAddressCommon, amountWei, gasLimit, gasPrice, nil,
	)

	// Set Chain ID for Ganache (1337 = default testnet)
	chainId := big.NewInt(1337)
	fmt.Println("Chain ID (manual):", chainId)

	// Convert hex private key to ECDSA format
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("error in convert private key : ", err)
		return
	}

	// Sign transaction with private key
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKeyECDSA)
	if err != nil {
		fmt.Println("error in sign transaction : ", err)
		return
	}

	// Send signed transaction to network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("error in send transaction : ", err)
		return
	}

	// Display success message and transaction hash
	fmt.Println("Done.\n")
	fmt.Printf("Transaction Hash : %s \n", signedTx.Hash().Hex())
}
