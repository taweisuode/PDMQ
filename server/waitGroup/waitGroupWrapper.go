package waitGroup

import "sync"

type WaitGroupWrapper struct {
	sync.WaitGroup
}

/**
 * @desc 它能够一直等到所有的goroutine执行完成，并且阻塞主线程的执行，直到所有的goroutine执行完成。
	如果不加这个，会出现这样的情况，主线程 执行完毕了 但是goroutine 并没有执行完。
	类似于在主进程与goroutine 中加锁，使得先执行所有goroutine的内容，然后再执行主进程
 * @param (query map[string]string)
 * @return (id int64, err error)
*/
func (w *WaitGroupWrapper) Wrap(f func()) {
	w.Add(1)
	go func() {
		f()
		w.Done()
	}()
}
