package test

import (
	"testing"

	"coredns/plugin/test"

	"github.com/miekg/dns"

	// Load all managed plugins in coredns
	_ "coredns/core/plugin"
)

func benchmarkLookupBalanceRewriteCache(b *testing.B) {
	t := new(testing.T)
	name, rm, err := test.TempFile(".", exampleOrg)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()

	corefile := `example.org:0 {
    file ` + name + `
    rewrite type ANY HINFO
    loadbalance
}
`

	ex, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer ex.Stop()

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Exchange(m, udp)
	}
}
