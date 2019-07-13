/**
 * @Time : 2019-07-13 17:01
 * @Author : zhuangjingpeng
 * @File : pdmqd_response
 * @Desc : file function description
 */
package common

const (
	Success = 0

	PubMsgReadError = 11000
	MsgTooBig       = 11001
	MsgEmpty        = 11002

	TopicMsgError = 12100
)

var RespMsg = map[int]string{
	Success:         "success",
	PubMsgReadError: "pub msg read error",
	MsgTooBig:       "pub msg is too big",
	MsgEmpty:        "pub msg is empty",
	TopicMsgError:   "topci pub msg error",
}
