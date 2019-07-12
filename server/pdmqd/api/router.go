/**
 * @Time : 2019-07-09 17:24
 * @Author : zhuangjingpeng
 * @File : router
 * @Desc : file function description
 */
package api

import (
	"PDMQ/internal/util"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	//fmt.Println("hello")
	sucJson := map[string]interface{}{}
	util.SendResult(c, 0, "成功", sucJson)
	return
}
