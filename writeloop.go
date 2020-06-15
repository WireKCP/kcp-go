package kcp

import (
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/net/ipv4"
)

func (t *UDPTunnel) writeSingle(msgs []ipv4.Message) {
	nbytes := 0
	npkts := 0
	for k := range msgs {
		if n, err := t.conn.WriteTo(msgs[k].Buffers[0], msgs[k].Addr); err == nil {
			nbytes += n
			npkts++
		} else {
			t.notifyWriteError(errors.WithStack(err))
		}
	}

	atomic.AddUint64(&DefaultSnmp.OutPkts, uint64(npkts))
	atomic.AddUint64(&DefaultSnmp.OutBytes, uint64(nbytes))
}

func (t *UDPTunnel) defaultWriteLoop() {
	for {
		select {
		case <-t.die:
			return
		case <-t.chFlush:
		}

		msgss := t.popMsgss()
		for _, msgs := range msgss {
			t.writeSingle(msgs)
		}
		t.releaseMsgss(msgss)
	}
}
