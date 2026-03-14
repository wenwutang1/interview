package main

type TransferReq struct {
	WalletFromId string `json:"wallet_from_id"`
	WallentToId  string `json:"wallet_to_id"`
	Amount       int64  `json:"amount"`
}
