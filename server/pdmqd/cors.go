/**
 * @Time : 2019-07-11 20:55
 * @Author : zhuangjingpeng
 * @File : cors
 * @Desc : file function description
 */
package pdmqd

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
)

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Next()
}

func AddTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceid := UniqueId()
		c.Set("traceid", traceid)
		c.Writer.Header().Set("X-REQUEST-ID", traceid)
	}
}

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
