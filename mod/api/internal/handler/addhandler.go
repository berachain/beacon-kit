package handler

import (
	"github.com/berachain/beacon-kit/mod/api/internal/types"
	"github.com/berachain/beacon-kit/mod/generate-genesis/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"strconv"
)

func AddAccountAndPredeploy(c *gin.Context) {

	format := c.DefaultQuery("format", "geth")

	var genesisJSON []byte
	switch format {
	case "geth":
		var newAllocation types.Alloc

		if err := c.BindJSON(&newAllocation); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		g := genesis.NewGenesis()
		gen := g.ToGethGenesis()

		for _, account := range newAllocation.Accounts {
			balance, success := new(big.Int).SetString(account.Balance, 10)
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account balance"})
				return
			}
			g.AddAccount(gen, common.HexToAddress(account.Address), balance)

		}
		for _, predeploy := range newAllocation.Predeploys {
			balance, success := new(big.Int).SetString(predeploy.Balance, 10)
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid predeploy balance"})
				return
			}
			nonce, err := strconv.ParseUint(strconv.FormatUint(predeploy.Nonce, 10), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nonce"})
				return
			}
			g.AddPredeploy(gen, common.HexToAddress(predeploy.Address), common.FromHex(predeploy.Code), nil, balance, nonce)
		}
		var err error
		genesisJSON, err = g.WriteFileToJSON(gen, "genesis-eth-api.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "nethermind":
		var newAllocation types.Alloc

		if err := c.BindJSON(&newAllocation); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		g := genesis.NewGenesis()
		ng := g.ToNethermindGenesis()

		for _, account := range newAllocation.Accounts {
			balance, success := new(big.Int).SetString(account.Balance, 10)
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account balance"})
				return
			}
			g.AddAccountNethermind(ng, common.HexToAddress(account.Address), balance)
		}
		for _, predeploy := range newAllocation.Predeploys {
			balance, success := new(big.Int).SetString(predeploy.Balance, 10)
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid predeploy balance"})
				return
			}
			nonce, err := strconv.ParseUint(strconv.FormatUint(predeploy.Nonce, 10), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nonce"})
				return
			}
			g.AddPredeployNethermind(ng, common.HexToAddress(predeploy.Address), common.FromHex(predeploy.Code), balance, nonce)
		}
		var err error
		genesisJSON, err = g.WriteNethermindGenesisToJSON(ng, "genesis-nethermind.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}

	c.Data(http.StatusOK, "application/json", genesisJSON)

}
