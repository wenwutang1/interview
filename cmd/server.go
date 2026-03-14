package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.POST("/wallets", createWallet())
	r.GET("/wallets/:wallet_id", getWalletOne())
	r.POST("/wallets/transfer", transferWallent())

	r.Run(":8000")
}
