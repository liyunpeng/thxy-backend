package api

import (
	"net/http"
	"thxy/types"
	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, types.CommonRes{
		Code: types.ResSuccessCode,
		Msg:  msg,
		Data: data,
	})
}

func JSONError(c *gin.Context, msg string, err interface{}) {
	c.JSON(http.StatusOK, types.CommonRes{
		Code:  types.ResErrorCode,
		Msg:   msg,
		Error: err,
	})
}

func JSONExpire(c *gin.Context, msg string, err interface{}) {
	c.JSON(http.StatusOK, types.CommonRes{
		Code:  types.ResExpireCode,
		Msg:   msg,
		Error: err,
	})
}

func JSONMaintain(c *gin.Context, msg string, err interface{}) {
	c.JSON(http.StatusOK, types.CommonRes{
		Code:  types.ResMaintainCode,
		Msg:   msg,
		Error: err,
	})
}
