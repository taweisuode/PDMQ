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
 *  @desc:  消息发布
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
	if server.ctx.pdmqd.config.MsgMaxSize < len(topicMsg) {
		seelog.Error("pub msg size is bigger than MsgMaxSize,current length is %d", len(topicMsg))
		util.SendResult(c, common.MsgTooBig, common.RespMsg[common.MsgTooBig], nil)
		return
	}
	if len(topicMsg) == 0 {
		seelog.Error("pub msg size is empty ")
		util.SendResult(c, common.MsgTooBig, common.RespMsg[common.MsgTooBig], nil)
		return
	}

	topic := server.ctx.pdmqd.GetTopic(topicName)

	msg := CreateMessage(topic.GenerateID(), topicMsg)
	err = topic.PutMessage(msg)
	if err != nil {
		util.SendResult(c, common.TopicMsgError, common.RespMsg[common.TopicMsgError], nil)
	}
	sucJson := map[string]interface{}{}
	sucJson["topic"] = topicName
	sucJson["msg"] = string(topicMsg)

	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}
