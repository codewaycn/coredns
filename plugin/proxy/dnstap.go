package proxy

import (
	"context"
	"time"

	"coredns/plugin/dnstap"
	"coredns/plugin/dnstap/msg"
	"coredns/request"

	tap "github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
)

func toDnstap(ctx context.Context, host string, ex Exchanger, state request.Request, reply *dns.Msg, start time.Time) error {
	tapper := dnstap.TapperFromContext(ctx)
	if tapper == nil {
		return nil
	}

	// Query
	b := msg.New().Time(start).HostPort(host)

	t := ex.Transport()
	if t == "" {
		t = state.Proto()
	}
	if t == "tcp" {
		b.SocketProto = tap.SocketProtocol_TCP
	} else {
		b.SocketProto = tap.SocketProtocol_UDP
	}

	if tapper.Pack() {
		b.Msg(state.Req)
	}
	m, err := b.ToOutsideQuery(tap.Message_FORWARDER_QUERY)
	if err != nil {
		return err
	}
	tapper.TapMessage(m)

	// Response
	if reply != nil {
		if tapper.Pack() {
			b.Msg(reply)
		}
		m, err := b.Time(time.Now()).ToOutsideResponse(tap.Message_FORWARDER_RESPONSE)
		if err != nil {
			return err
		}
		tapper.TapMessage(m)
	}

	return nil
}
