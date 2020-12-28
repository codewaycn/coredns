// Package log implements basic but useful request (access) logging plugin.
package log

import (
	"context"
	"time"

	"coredns/plugin"
	"coredns/plugin/metrics/vars"
	"coredns/plugin/pkg/dnstest"
	clog "coredns/plugin/pkg/log"
	"coredns/plugin/pkg/rcode"
	"coredns/plugin/pkg/replacer"
	"coredns/plugin/pkg/response"
	"coredns/request"

	"github.com/miekg/dns"
)

// Logger is a basic request logging plugin.
type Logger struct {
	Next      plugin.Handler
	Rules     []Rule
	ErrorFunc func(context.Context, dns.ResponseWriter, *dns.Msg, int) // failover error handler
}

// ServeDNS implements the plugin.Handler interface.
func (l Logger) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	for _, rule := range l.Rules {
		if !plugin.Name(rule.NameScope).Matches(state.Name()) {
			continue
		}

		rrw := dnstest.NewRecorder(w)
		rc, err := plugin.NextOrFailure(l.Name(), l.Next, ctx, rrw, r)

		if rc > 0 {
			// There was an error up the chain, but no response has been written yet.
			// The error must be handled here so the log entry will record the response size.
			if l.ErrorFunc != nil {
				l.ErrorFunc(ctx, rrw, r, rc)
			} else {
				answer := new(dns.Msg)
				answer.SetRcode(r, rc)
				state.SizeAndDo(answer)

				vars.Report(ctx, state, vars.Dropped, rcode.ToString(rc), answer.Len(), time.Now())

				w.WriteMsg(answer)
			}
			rc = 0
		}

		tpe, _ := response.Typify(rrw.Msg, time.Now().UTC())
		class := response.Classify(tpe)
		// If we don't set up a class in config, the default "all" will be added
		// and we shouldn't have an empty rule.Class.
		if rule.Class[response.All] || rule.Class[class] {
			rep := replacer.New(r, rrw, CommonLogEmptyValue)
			clog.Infof(rep.Replace(rule.Format))
		}

		return rc, err

	}
	return plugin.NextOrFailure(l.Name(), l.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (l Logger) Name() string { return "log" }

// Rule configures the logging plugin.
type Rule struct {
	NameScope string
	Class     map[response.Class]bool
	Format    string
}

const (
	// CommonLogFormat is the common log format.
	CommonLogFormat = `{remote}:{port} ` + CommonLogEmptyValue + ` {>id} "{type} {class} {name} {proto} {size} {>do} {>bufsize}" {rcode} {>rflags} {rsize} {duration}`
	// CommonLogEmptyValue is the common empty log value.
	CommonLogEmptyValue = "-"
	// CombinedLogFormat is the combined log format.
	CombinedLogFormat = CommonLogFormat + ` "{>opcode}"`
	// DefaultLogFormat is the default log format.
	DefaultLogFormat = CommonLogFormat
)
