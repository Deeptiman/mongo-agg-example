package main

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"sync/atomic"
)

var (

	// mongo conn pool
	connPoolCreated uint64
	connSuccess     uint64
	connClosed      uint64
	connReturned    uint64
	connPoolCleared uint64

	// cmd Monitor
	cmdMonitorStart   = []*event.CommandStartedEvent{}
	cmdMonitorSucceed = []*event.CommandSucceededEvent{}
	cmdMonitorFailed  = []*event.CommandFailedEvent{}

	// server Heartbeat
	serverHeartbeatStart     = []*event.ServerHeartbeatStartedEvent{}
	serverHeartbeatOpenning  = []*event.ServerOpeningEvent{}
	serverHeartbeatSucceeded = []*event.ServerHeartbeatSucceededEvent{}
	serverHeartbeatFailed    = []*event.ServerHeartbeatFailedEvent{}
	serverHeartbeatClosed    = []*event.ServerClosedEvent{}
)

type MongoHealthCheck struct {
	// mongo conn pool
	ConnPoolCreated uint64 `json:"connPoolCreated"`
	ConnSuccess     uint64 `json:"connSuccess"`
	ConnClosed      uint64 `json:"connClosed"`
	ConnReturned    uint64 `json:"connReturned"`
	ConnPoolCleared uint64 `json:"connPoolCleared"`

	// server Heartbeat
	ServerHeartbeatStart     []*event.ServerHeartbeatStartedEvent   `json:"serverHeartbeatStart"`
	ServerHeartbeatOpenning  []*event.ServerOpeningEvent            `json:"serverHeartbeatOpenning"`
	ServerHeartbeatSucceeded []*event.ServerHeartbeatSucceededEvent `json:"serverHeartbeatSucceeded"`
	ServerHeartbeatFailed    []*event.ServerHeartbeatFailedEvent    `json:"serverHeartbeatFailed"`
	ServerHeartbeatClosed    []*event.ServerClosedEvent             `json:"serverHeartbeatClosed"`

	// mongo commands
	CmdMonitorStart   []*event.CommandStartedEvent   `json:"cmdMonitorStart"`
	CmdMonitorSucceed []*event.CommandSucceededEvent `json:"cmdMonitorSucceed"`
	CmdMonitorFailed  []*event.CommandFailedEvent    `json:"cmdMonitorFailed"`
}

func MongoMonitors() (*event.PoolMonitor, *event.CommandMonitor, *event.ServerMonitor) {

	poolMonitor := &event.PoolMonitor{
		Event: func(e *event.PoolEvent) {
			switch e.Type {
			case event.PoolCreated:
				atomic.AddUint64(&connPoolCreated, 1)
				log.Info("PoolMonitor", "PoolCreated", getConnPoolCreated())
			case event.GetSucceeded:
				atomic.AddUint64(&connSuccess, 1)
				log.Info("PoolMonitor", "ConnSucces", getConnPoolSuccess())
			case event.ConnectionClosed:
				atomic.AddUint64(&connClosed, 1)
				log.Info("PoolMonitor", "ConnClosed", getConnPoolClosed())
			case event.ConnectionReturned:
				atomic.AddUint64(&connReturned, 1)
				log.Info("PoolMonitor", "ConnReturned", getConnPoolReturned())
			case event.PoolCleared:
				atomic.AddUint64(&connPoolCleared, 1)
				log.Info("PoolMonitor", "ConnPoolCleared", getConnPoolCleared())
			}
		},
	}

	// cmd Monitor
	cmdMonitorStart = []*event.CommandStartedEvent{}
	cmdMonitorSucceed = []*event.CommandSucceededEvent{}
	cmdMonitorFailed = []*event.CommandFailedEvent{}

	// server Heartbeat
	serverHeartbeatStart = []*event.ServerHeartbeatStartedEvent{}
	serverHeartbeatOpenning = []*event.ServerOpeningEvent{}
	serverHeartbeatSucceeded = []*event.ServerHeartbeatSucceededEvent{}
	serverHeartbeatFailed = []*event.ServerHeartbeatFailedEvent{}
	serverHeartbeatClosed = []*event.ServerClosedEvent{}

	cmdMonitor := &event.CommandMonitor{

		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			cmdMonitorStart = append(cmdMonitorStart, evt)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			cmdMonitorSucceed = append(cmdMonitorSucceed, evt)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			cmdMonitorFailed = append(cmdMonitorFailed, evt)
		},
	}

	serverMonitor := &event.ServerMonitor{
		ServerHeartbeatStarted: func(evt *event.ServerHeartbeatStartedEvent) {
			serverHeartbeatStart = append(serverHeartbeatStart, evt)
		},
		ServerOpening: func(evt *event.ServerOpeningEvent) {
			serverHeartbeatOpenning = append(serverHeartbeatOpenning, evt)
		},
		ServerHeartbeatSucceeded: func(evt *event.ServerHeartbeatSucceededEvent) {
			serverHeartbeatSucceeded = append(serverHeartbeatSucceeded, evt)
		},
		ServerHeartbeatFailed: func(evt *event.ServerHeartbeatFailedEvent) {
			serverHeartbeatFailed = append(serverHeartbeatFailed, evt)
		},
		ServerClosed: func(evt *event.ServerClosedEvent) {
			serverHeartbeatClosed = append(serverHeartbeatClosed, evt)
		},
	}

	return poolMonitor, cmdMonitor, serverMonitor
}

func ReadMongoHealthMonitor() *MongoHealthCheck {

	return &MongoHealthCheck{
		ConnPoolCreated:          getConnPoolCreated(),
		ConnSuccess:              getConnPoolSuccess(),
		ConnClosed:               getConnPoolClosed(),
		ConnReturned:             getConnPoolReturned(),
		ConnPoolCleared:          getConnPoolCleared(),
		ServerHeartbeatStart:     serverHeartbeatStart,
		ServerHeartbeatOpenning:  serverHeartbeatOpenning,
		ServerHeartbeatSucceeded: serverHeartbeatSucceeded,
		ServerHeartbeatFailed:    serverHeartbeatFailed,
		ServerHeartbeatClosed:    serverHeartbeatClosed,
		CmdMonitorStart:          cmdMonitorStart,
		CmdMonitorSucceed:        cmdMonitorSucceed,
		CmdMonitorFailed:         cmdMonitorFailed,
	}
}

func getConnPoolCreated() uint64 {
	return atomic.LoadUint64(&connPoolCreated)
}

func getConnPoolSuccess() uint64 {
	return atomic.LoadUint64(&connSuccess)
}

func getConnPoolReturned() uint64 {
	return atomic.LoadUint64(&connReturned)
}

func getConnPoolClosed() uint64 {
	return atomic.LoadUint64(&connClosed)
}

func getConnPoolCleared() uint64 {
	return atomic.LoadUint64(&connPoolCleared)
}
