/**
 * @Time : 2019-07-09 17:24
 * @Author : zhuangjingpeng
 * @File : router
 * @Desc : file function description
 */
package pdmqd

import (
	"PDMQ/internal/common"
	"PDMQ/internal/util"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func (server *httpServer) Ping(c *gin.Context) {
	sucJson := map[string]interface{}{}
	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}

/**
 *  @desc:  description
 *  @input: data String
 *  @resp:  err resp
 *
**/

func (server *httpServer) Pub(c *gin.Context) {
	topicName := c.Query("topic")
	topicMsg, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		seelog.Errorf("read topic msg error is %v", err.Error())
		util.SendResult(c, common.PubMsgReadError, common.RespMsg[common.PubMsgReadError], nil)
		return
	}
	fmt.Println(server.ctx.pdmqd.config.MsgMaxSize, len(topicMsg))
	if server.ctx.pdmqd.config.MsgMaxSize < len(topicMsg) {
		seelog.Error("pub msg size is bigger than MsgMaxSize,current length is %d", len(topicMsg))
		util.SendResult(c, common.MsgTooBig, common.RespMsg[common.MsgTooBig], nil)
		return
	}
	sucJson := map[string]interface{}{}
	sucJson["topic"] = topicName
	sucJson["msg"] = string(topicMsg)

	util.PrintJson("response data is ", sucJson)
	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}
