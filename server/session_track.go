package main

import (
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
)

// Logic to track session -> connection mapping via websocket messages due to
// arguments coming from the slash commands does not bring the connection ID
// to the plugin.

func (p *Plugin) initializeSessionTracking() {
	p.sessionToConn = make(map[string]string)
	p.sessionToConnMu = sync.RWMutex{}
}

func (p *Plugin) WebSocketMessageHasBeenPosted(webConnID, userID string, req *model.WebSocketRequest) {
	if req.Session.Id != "" {
		p.sessionToConnMu.Lock()
		p.sessionToConn[req.Session.Id] = webConnID
		p.sessionToConnMu.Unlock()
	}
}

func (p *Plugin) OnWebSocketDisconnect(webConnID, userID string) {
	p.sessionToConnMu.Lock()
	defer p.sessionToConnMu.Unlock()

	for sid, cid := range p.sessionToConn {
		if cid == webConnID {
			delete(p.sessionToConn, sid)
		}
	}
}

func (p *Plugin) GetConnectionIDForSession(sessionID string) (string, bool) {
	p.sessionToConnMu.RLock()
	defer p.sessionToConnMu.RUnlock()

	connID, ok := p.sessionToConn[sessionID]
	return connID, ok
}
