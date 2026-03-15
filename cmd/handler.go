package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createWallet() gin.HandlerFunc {
	return func(c *gin.Context) {

		wallentItem := NewWallet()
		GlobalWalletCollect.Add(wallentItem)

		c.JSON(http.StatusOK, map[string]any{
			"code": http.StatusOK,
			"msg":  "create success",
			"data": map[string]any{
				"wallent_id": wallentItem.WalletId,
				"balance":    wallentItem.Balance,
			},
		})
	}
}

func getWalletOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallentId := c.Param("wallet_id")

		wallentItem := GlobalWalletCollect.Get(wallentId)
		if wallentItem == nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": http.StatusNotFound,
				"msg":  "not found",
			})
			return
		}
		c.JSON(http.StatusOK, map[string]any{
			"code": http.StatusOK,
			"msg":  "query success",
			"data": map[string]any{
				"wallent_id": wallentItem.WalletId,
				"balance":    wallentItem.Balance,
			},
		})
	}
}

func transferWallent() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &TransferReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": http.StatusBadRequest,
				"msg":  "invalid request",
			})
			return
		}
		if req.WalletFromId == "" || req.WallentToId == "" || req.Amount <= 0 {
			c.JSON(http.StatusOK, map[string]any{
				"code": http.StatusBadRequest,
				"msg":  "invalid request",
			})
			return
		}

		err := GlobalWalletCollect.Transfer(req.WalletFromId, req.WallentToId, req.Amount)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": http.StatusTooManyRequests,
				"msg":  "server busy",
			})
			return
		}

		c.JSON(http.StatusOK, map[string]any{
			"code": http.StatusOK,
			"msg":  "success",
		})
	}
}
