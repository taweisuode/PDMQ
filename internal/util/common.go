/**
 * @Time : 2019-07-08 19:55
 * @Author : zhuangjingpeng
 * @File : common
 * @Desc : file function description
 */
package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func PrintJson(desc string, data interface{}) {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(desc, string(res))
}

// SendResult send json result to client
func SendResult(c *gin.Context, errCode int, msg string, data interface{}) {
	result := map[string]interface{}{
		"code":        errCode,
		"message":     msg,
		"currentTime": time.Now().Unix(),
		"data":        data,
	}
	if errCode != 0 {
		c.Set("response", result)
		c.Set("errCode", errCode)
		c.Set("errMsg", msg)
	}
	c.JSON(200, result)
	return
}
