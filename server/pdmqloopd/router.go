/**
 * @Time : 2019-07-09 17:24
 * @Author : zhuangjingpeng
 * @File : router
 * @Desc : file function description
 */
package pdmqloopd

import (
	"PDMQ/internal/common"
	"PDMQ/internal/util"
	"github.com/gin-gonic/gin"
)

func (server *httpServer) Ping(c *gin.Context) {
	sucJson := map[string]interface{}{}
	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}
