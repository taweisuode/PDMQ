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
)

func PrintJson(desc string, data interface{}) {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(desc, string(res))
}
