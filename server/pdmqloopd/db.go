/**
 * @Time : 2019-07-31 17:10
 * @Author : zhuangjingpeng
 * @File : db
 * @Desc : file function description
 */
package pdmqloopd

import (
	"sync"
	"time"
)

type RegistrationDB struct {
	sync.RWMutex
	RegistrationMap map[Registration]ProducerMap
}

type Registration struct {
	Category    string
	TopicName   string
	ChannelName string
}

type Registrations []Registration
type PeerInfo struct {
	lastUpdate       int64
	id               string
	RemoteAddress    string `json:"remote_address"`
	Hostname         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TCPPort          int    `json:"tcp_port"`
	HTTPPort         int    `json:"http_port"`
	Version          string `json:"version"`
}

type Producer struct {
	peerInfo *PeerInfo

	tombstoned   bool
	tombstonedAt time.Time
}

type Producers []*Producer
type ProducerMap map[string]*Producer

func NewRegistrationDB() *RegistrationDB {
	return &RegistrationDB{
		RegistrationMap: make(map[Registration]ProducerMap),
	}
}
func (r *RegistrationDB) needFilter(topicName string, channelName string) bool {
	return topicName == "*" || channelName == "*"
}
func (r *RegistrationDB) FindRegistrations(category string, topicName string, channelName string) Registrations {
	r.RLock()
	defer r.RUnlock()
	if !r.needFilter(topicName, channelName) {
		k := Registration{category, topicName, channelName}
		if _, ok := r.RegistrationMap[k]; ok {
			return Registrations{k}
		}
		return Registrations{}
	}
	results := Registrations{}
	for k := range r.RegistrationMap {
		if !k.IsMatch(category, topicName, channelName) {
			continue
		}
		results = append(results, k)
	}
	return results
}

func (r *RegistrationDB) FindProducers(category string, topicName string, channelName string) Producers {
	r.RLock()
	defer r.RUnlock()
	if !r.needFilter(topicName, channelName) {
		registration := Registration{category, topicName, channelName}
		var producers Producers
		for _, producer := range r.RegistrationMap[registration] {
			producers = append(producers, producer)
		}
		return producers
	}

	results := make(map[string]struct{})
	var retProducers Producers
	for k, producers := range r.RegistrationMap {
		if !k.IsMatch(category, topicName, channelName) {
			continue
		}
		for _, producer := range producers {
			_, found := results[producer.peerInfo.id]
			if found == false {
				results[producer.peerInfo.id] = struct{}{}
				retProducers = append(retProducers, producer)
			}
		}
	}
	return retProducers
}
func (k *Registration) IsMatch(category string, topicName string, channelName string) bool {
	if category != k.Category {
		return false
	}
	if topicName != "*" && k.TopicName != topicName {
		return false
	}
	if channelName != "*" && k.ChannelName != channelName {
		return false
	}
	return true
}

func (rr Registrations) ChannelKeys() []string {
	subkeys := make([]string, len(rr))
	for i, k := range rr {
		subkeys[i] = k.ChannelName
	}
	return subkeys
}

//寻找租约时间内有效的生产者
func (pp Producers) FilterByActive(tombstoneLifetime time.Duration) Producers {
	results := Producers{}
	for _, p := range pp {
		if p.IsTombstoned(tombstoneLifetime) {
			continue
		}
		results = append(results, p)
	}
	return results
}

//生产者是否有效
func (p *Producer) IsTombstoned(lifetime time.Duration) bool {
	return p.tombstoned && time.Now().Sub(p.tombstonedAt) < lifetime
}

//返回生产者的所有信息（host port 等等）
func (pp Producers) PeerInfo() []*PeerInfo {
	results := []*PeerInfo{}
	for _, p := range pp {
		results = append(results, p.peerInfo)
	}
	return results
}
