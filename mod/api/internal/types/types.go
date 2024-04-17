// mod/api/types/types.go
package types

type Account struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type Predeploy struct {
	Address string `json:"predeployAddress"`
	Code    string `json:"code"`
	Balance string `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

type Alloc struct {
	Accounts   []Account   `json:"accounts"`
	Predeploys []Predeploy `json:"predeploys"`
}
