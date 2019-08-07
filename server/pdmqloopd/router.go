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
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func (server *httpServer) Ping(c *gin.Context) {
	sucJson := map[string]interface{}{}
	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}

func (server *httpServer) Lookup(c *gin.Context) {
	topicName := c.GetString("topic")
	if topicName != "" {
		seelog.Error("get topic name is nil ")
		util.SendResult(c, common.LoopGetTopicError, common.RespMsg[common.LoopGetTopicError], make(map[string]interface{}))
		return
	}
	fmt.Printf("server.ctx.pdmqloopd.DB is [%+v]\n", server.ctx.pdmqloopd.DB)
	registration := server.ctx.pdmqloopd.DB.FindRegistrations("topic", topicName, "")

	fmt.Printf("registration is [%+v],len is [%d]\n", registration, len(registration))
	if len(registration) == 0 {
		seelog.Error("find topic is nil ,topicName is [%+v]\n", topicName)
		util.SendResult(c, common.LoopFindTopicError, common.RespMsg[common.LoopFindTopicError], make(map[string]interface{}))
		return
	}

	channels := server.ctx.pdmqloopd.DB.FindRegistrations("channel", topicName, "*").ChannelKeys()
	producers := server.ctx.pdmqloopd.DB.FindProducers("topic", topicName, "")
	//生产者需要根据 墓碑生存时间来（奇怪的名字） 进行过滤
	producers = producers.FilterByActive(server.ctx.pdmqloopd.config.TombstoneLifetime)
	sucJson := map[string]interface{}{
		"channels":  channels,
		"producers": producers.PeerInfo(),
	}
	util.SendResult(c, common.Success, common.RespMsg[common.Success], sucJson)
	return
}
