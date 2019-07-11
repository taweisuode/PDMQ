/**
 * @Time : 2019-07-08 17:45
 * @Author : zhuangjingpeng
 * @File : wait_group_wrapper
 * @Desc : file function description
 */
package util

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(f func()) {
	w.Add(1)
	go func() {
		f()
		w.Done()
	}()
}
