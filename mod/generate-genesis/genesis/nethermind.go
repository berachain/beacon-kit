package genesis

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
)

type NethermindGenesis struct {
	Name    string `json:"name"`
	DataDir string `json:"dataDir"`
	Engine  struct {
		Ethash struct {
			Params struct {
				MinimumDifficulty      string   `json:"minimumDifficulty"`
				DifficultyBoundDivisor string   `json:"difficultyBoundDivisor"`
				DurationLimit          string   `json:"durationLimit"`
				BlockReward            struct{} `json:"blockReward"`
				HomesteadTransition    string   `json:"homesteadTransition"`
				DaoHardforkTransition  string   `json:"daoHardforkTransition"`
				Eip100bTransition      string   `json:"eip100bTransition"`
				DifficultyBombDelays   struct{} `json:"difficultyBombDelays"`
			} `json:"params"`
		} `json:"Ethash"`
	} `json:"engine"`
	Params struct {
		GasLimitBoundDivisor          string `json:"gasLimitBoundDivisor"`
		Registrar                     string `json:"registrar"`
		AccountStartNonce             string `json:"accountStartNonce"`
		MaximumExtraDataSize          string `json:"maximumExtraDataSize"`
		MinGasLimit                   string `json:"minGasLimit"`
		NetworkID                     string `json:"networkID"`
		ForkBlock                     string `json:"forkBlock"`
		MaxCodeSize                   string `json:"maxCodeSize"`
		MaxCodeSizeTransition         string `json:"maxCodeSizeTransition"`
		Eip150Transition              string `json:"eip150Transition"`
		Eip160Transition              string `json:"eip160Transition"`
		Eip161abcTransition           string `json:"eip161abcTransition"`
		Eip161dTransition             string `json:"eip161dTransition"`
		Eip155Transition              string `json:"eip155Transition"`
		Eip140Transition              string `json:"eip140Transition"`
		Eip211Transition              string `json:"eip211Transition"`
		Eip214Transition              string `json:"eip214Transition"`
		Eip658Transition              string `json:"eip658Transition"`
		Eip145Transition              string `json:"eip145Transition"`
		Eip1014Transition             string `json:"eip1014Transition"`
		Eip1052Transition             string `json:"eip1052Transition"`
		Eip1283Transition             string `json:"eip1283Transition"`
		Eip1283DisableTransition      string `json:"eip1283DisableTransition"`
		Eip152Transition              string `json:"eip152Transition"`
		Eip1108Transition             string `json:"eip1108Transition"`
		Eip1344Transition             string `json:"eip1344Transition"`
		Eip1884Transition             string `json:"eip1884Transition"`
		Eip2028Transition             string `json:"eip2028Transition"`
		Eip2200Transition             string `json:"eip2200Transition"`
		Eip2565Transition             string `json:"eip2565Transition"`
		Eip2929Transition             string `json:"eip2929Transition"`
		Eip2930Transition             string `json:"eip2930Transition"`
		Eip1559Transition             string `json:"eip1559Transition"`
		Eip3198Transition             string `json:"eip3198Transition"`
		Eip3529Transition             string `json:"eip3529Transition"`
		Eip3541Transition             string `json:"eip3541Transition"`
		Eip4895TransitionTimestamp    string `json:"eip4895TransitionTimestamp"`
		Eip3855TransitionTimestamp    string `json:"eip3855TransitionTimestamp"`
		Eip3651TransitionTimestamp    string `json:"eip3651TransitionTimestamp"`
		Eip3860TransitionTimestamp    string `json:"eip3860TransitionTimestamp"`
		Eip1153TransitionTimestamp    string `json:"eip1153TransitionTimestamp"`
		Eip4788TransitionTimestamp    string `json:"eip4788TransitionTimestamp"`
		Eip4844TransitionTimestamp    string `json:"eip4844TransitionTimestamp"`
		Eip5656TransitionTimestamp    string `json:"eip5656TransitionTimestamp"`
		Eip6780TransitionTimestamp    string `json:"eip6780TransitionTimestamp"`
		TerminalTotalDifficulty       string `json:"terminalTotalDifficulty"`
		TerminalTotalDifficultyPassed bool   `json:"terminalTotalDifficultyPassed"`
	} `json:"params"`
	Genesis struct {
		Coinbase   string `json:"coinbase"`
		Difficulty string `json:"difficulty"`
		ExtraData  string `json:"extraData"`
		GasLimit   string `json:"gasLimit"`
		Nonce      string `json:"nonce"`
		Timestamp  string `json:"timestamp"`
	} `json:"genesis"`
	Accounts map[string]struct {
		Balance string `json:"balance"`
		Nonce   string `json:"nonce"`
		Code    string `json:"code"`
	} `json:"accounts"`
}

func (n *NethermindGenesis) AddAccount(address common.Address, balance *big.Int) {
	n.Accounts[address.Hex()] = struct {
		Balance string `json:"balance"`
		Nonce   string `json:"nonce"`
		Code    string `json:"code"`
	}{
		Balance: "0x" + balance.Text(16), // Convert balance to hexadecimal

	}

}

func (n *NethermindGenesis) AddPredeploy(address common.Address, code []byte, balance *big.Int, nonce uint64) {
	n.Accounts[address.Hex()] = struct {
		Balance string `json:"balance"`
		Nonce   string `json:"nonce"`
		Code    string `json:"code"`
	}{
		Balance: "0x" + balance.Text(16),              // Convert balance to hexadecimal
		Nonce:   "0x" + strconv.FormatUint(nonce, 16), // Convert nonce to hexadecimal
		Code:    "0x" + common.Bytes2Hex(code),        // Convert code to hexadecimal
	}
}

func (n *NethermindGenesis) ToJSON(filename string) error {
	_, err := WriteGenesisToJSON(n, filename)
	return err
}

func (n *NethermindGenesis) ToNethermindGenesis() *NethermindGenesis {
	ng := n.initializeNethermindGenesis()
	n.populateEIPTransitions(ng)
	return ng
}

func (n *NethermindGenesis) initializeNethermindGenesis() *NethermindGenesis {
	ng := &NethermindGenesis{}
	// Populate the NethermindGenesis struct with the necessary data
	ng.Name = "Ethereum"
	ng.DataDir = "ethereum"
	ng.Engine.Ethash.Params.MinimumDifficulty = "0x0"
	ng.Engine.Ethash.Params.DifficultyBoundDivisor = "0x0"
	ng.Engine.Ethash.Params.DurationLimit = "0x0"
	ng.Params.GasLimitBoundDivisor = "0x400"
	ng.Params.Registrar = "0xe3389675d0338462dC76C6f9A3e432550c36A142"
	ng.Params.AccountStartNonce = "0x0"
	ng.Params.MaximumExtraDataSize = "0x20"
	ng.Params.MinGasLimit = "0x1c9c380"
	ng.Params.NetworkID = "0x138d7"
	ng.Params.ForkBlock = "0x0"
	ng.Params.MaxCodeSize = "0x6000"
	ng.Params.MaxCodeSizeTransition = "0x0"
	ng.Genesis.Coinbase = "0x0000000000000000000000000000000000000000"
	ng.Genesis.Difficulty = "0x0"
	ng.Genesis.ExtraData = "0x0000000000000000000000000000000000000000000000000000000000000000658bdf435d810c91414ec09147daa6db624063790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	ng.Genesis.GasLimit = "0x1c9c380"
	ng.Genesis.Nonce = "0x0000000000000000"
	ng.Genesis.Timestamp = "0x0"
	ng.Accounts = make(map[string]struct {
		Balance string `json:"balance"`
		Nonce   string `json:"nonce"`
		Code    string `json:"code"`
	})

	return ng
}

func (n *NethermindGenesis) populateEIPTransitions(ng *NethermindGenesis) {
	ng.Params.Eip150Transition = "0x0"
	ng.Params.Eip160Transition = "0x0"
	ng.Params.Eip161abcTransition = "0x0"
	ng.Params.Eip161dTransition = "0x0"
	ng.Params.Eip155Transition = "0x0"
	ng.Params.Eip140Transition = "0x0"
	ng.Params.Eip211Transition = "0x0"
	ng.Params.Eip214Transition = "0x0"
	ng.Params.Eip658Transition = "0x0"
	ng.Params.Eip145Transition = "0x0"
	ng.Params.Eip1014Transition = "0x0"
	ng.Params.Eip1052Transition = "0x0"
	ng.Params.Eip1283Transition = "0x0"
	ng.Params.Eip1283DisableTransition = "0x0"
	ng.Params.Eip152Transition = "0x0"
	ng.Params.Eip1108Transition = "0x0"
	ng.Params.Eip1344Transition = "0x0"
	ng.Params.Eip1884Transition = "0x0"
	ng.Params.Eip2028Transition = "0x0"
	ng.Params.Eip2200Transition = "0x0"
	ng.Params.Eip2565Transition = "0x0"
	ng.Params.Eip2929Transition = "0x0"
	ng.Params.Eip2930Transition = "0x0"
	ng.Params.Eip1559Transition = "0x0"
	ng.Params.Eip3198Transition = "0x0"
	ng.Params.Eip3529Transition = "0x0"
	ng.Params.Eip3541Transition = "0x0"
	ng.Params.Eip4895TransitionTimestamp = "0x0"
	ng.Params.Eip3855TransitionTimestamp = "0x0"
	ng.Params.Eip3651TransitionTimestamp = "0x0"
	ng.Params.Eip3860TransitionTimestamp = "0x0"
	ng.Params.Eip1153TransitionTimestamp = "0x0"
	ng.Params.Eip4788TransitionTimestamp = "0x0"
	ng.Params.Eip4844TransitionTimestamp = "0x0"
	ng.Params.Eip5656TransitionTimestamp = "0x0"
	ng.Params.Eip6780TransitionTimestamp = "0x0"
	ng.Params.TerminalTotalDifficulty = "0"
	ng.Params.TerminalTotalDifficultyPassed = true
}
