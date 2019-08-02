/**
 * @Time : 2019-08-01 11:58
 * @Author : zhuangjingpeng
 * @File : loop
 * @Desc : pdmqd 对 loopd 的监听
 */
package pdmqd

import (
	"GO_PDMQ"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"time"
)

func (pdmqd *PDMQD) loop() {
	hostname, err := os.Hostname()
	if err != nil {
		seelog.Errorf("failed to get hostname - %s", err)
		os.Exit(1)
	}
	ticker := time.NewTicker(15 * time.Second)
	loopPeers := make([]*loopPeer, 0)
	lookupAddrs := make([]string, 0)

	fmt.Println(1111, pdmqd.config.LoopTCPAddresses)
	for {
		for _, address := range pdmqd.config.LoopTCPAddresses {
			loopPeer := newLookupPeer(address, pdmqd.config.MsgMaxSize,
				connectCallback(pdmqd, hostname))
			loopPeer.Command(nil)
			loopPeers = append(loopPeers, loopPeer)
			lookupAddrs = append(lookupAddrs, address)
		}

		fmt.Printf("[PDMQLOOP] [%+v] loopPeers: [%+v],lookupAddrs: [%+v]\n", time.Now().Format("2006-01-02 15:04:05"), loopPeers, lookupAddrs)
		select {
		case <-ticker.C:
			//没15秒探活
			for _, loopPeer := range loopPeers {
				cmd := pdmq.Ping()
				_, err := loopPeer.Command(cmd)
				if err != nil {
					seelog.Errorf("LOOKUPD(%s): %s - %s", loopPeer, cmd, err)
				}
			}
		}
	}
}

func connectCallback(pdmqd *PDMQD, hostname string) func(*loopPeer) {
	return func(lp *loopPeer) {
		ci := make(map[string]interface{})
		ci["version"] = pdmqd.config.Version
		ci["tcp_port"] = pdmqd.RealTCPAddr().Port
		ci["http_port"] = pdmqd.RealHTTPAddr().Port
		ci["hostname"] = hostname
		ci["broadcast_address"] = pdmqd.config.BroadcastAddress

		cmd, err := pdmq.Identify(ci)
		if err != nil {
			lp.Close()
			return
		}
		resp, err := lp.Command(cmd)
		if err != nil {
			seelog.Errorf("LOOKUPD(%s): %s - %s", lp, cmd, err)
			return
		} else if bytes.Equal(resp, []byte("E_INVALID")) {
			seelog.Errorf("LOOKUPD(%s): lookupd returned %s", lp, resp)
			lp.Close()
			return
		} else {
			err = json.Unmarshal(resp, &lp.Info)
			if err != nil {
				seelog.Errorf("LOOKUPD(%s): parsing response - %s", lp, resp)
				lp.Close()
				return
			} else {
				seelog.Infof("LOOKUPD(%s): peer info %+v", lp, lp.Info)
				if lp.Info.BroadcastAddress == "" {
					seelog.Errorf("LOOKUPD(%s): no broadcast address", lp)
				}
			}
		}

		// build all the commands first so we exit the lock(s) as fast as possible
		var commands []*pdmq.Command
		pdmqd.RLock()
		for _, topic := range pdmqd.topicMap {
			topic.RLock()
			if len(topic.channelMap) == 0 {
				commands = append(commands, pdmq.Register(topic.topicName, ""))
			} else {
				for _, channel := range topic.channelMap {
					commands = append(commands, pdmq.Register(channel.topicName, channel.ChannelName))
				}
			}
			topic.RUnlock()
		}
		pdmqd.RUnlock()

		for _, cmd := range commands {
			seelog.Infof("LOOKUPD(%s): %s", lp, cmd)
			_, err := lp.Command(cmd)
			if err != nil {
				seelog.Errorf("LOOKUPD(%s): %s - %s", lp, cmd, err)
				return
			}
		}
	}
}
