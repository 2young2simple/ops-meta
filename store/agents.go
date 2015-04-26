package store

import (
	"gitcafe.com/ops/common/model"
	"sync"
)

type AgentsMap struct {
	sync.RWMutex
	M map[string]*model.RealAgent
}

func NewAgentsMap() *AgentsMap {
	return &AgentsMap{M: make(map[string]*model.RealAgent)}
}

func (this *AgentsMap) Get(agentName string) (*model.RealAgent, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[agentName]
	return val, exists
}

func (this *AgentsMap) Put(agentName string, realAgent *model.RealAgent) {
	this.Lock()
	defer this.Unlock()
	this.M[agentName] = realAgent
}

type HostAgentsMap struct {
	sync.RWMutex
	M map[string]*AgentsMap
}

func NewHostAgentsMap() *HostAgentsMap {
	return &HostAgentsMap{M: make(map[string]*AgentsMap)}
}

var HostAgents = NewHostAgentsMap()

func (this *HostAgentsMap) Get(hostname string) (*AgentsMap, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[hostname]
	return val, exists
}

func (this *HostAgentsMap) Put(hostname string, am *AgentsMap) {
	this.Lock()
	defer this.Unlock()
	this.M[hostname] = am
}

func ParseHeartbeatRequest(req *model.HeartbeatRequest) {
	agentsMap, exists := HostAgents.Get(req.Hostname)
	if exists {
		for _, a := range req.RealAgents {
			agentsMap.Put(a.Name, a)
		}
	} else {
		am := NewAgentsMap()
		for _, a := range req.RealAgents {
			am.Put(a.Name, a)
		}
		HostAgents.Put(req.Hostname, am)
	}
}
