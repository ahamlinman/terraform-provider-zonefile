package provider

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func readZone(origin, content string) ([]dns.RR, error) {
	var rrs []dns.RR
	parser := dns.NewZoneParser(strings.NewReader(content), origin, "")
	for rr, ok := parser.Next(); ok; rr, ok = parser.Next() {
		rrs = append(rrs, rr)
	}
	return rrs, parser.Err()
}

type rrSet struct {
	Hdr dns.RR_Header
	RRs []dns.RR
}

func groupRRs(rrs []dns.RR) ([]rrSet, error) {
	type key struct {
		Name   string
		Class  uint16
		Rrtype uint16
	}

	// Care is taken to order RRSets based on the ordering of the original RRs.
	// This is friendly to the unit tests, and to users of count in Terraform
	// (though I'd recommend using for_each and inventing keys with
	// md5(jsonencode(rrset)) or similar).
	var rrSets []rrSet
	indices := make(map[key]int)
	for _, rr := range rrs {
		hdr := *rr.Header()
		k := key{hdr.Name, hdr.Class, hdr.Rrtype}
		if i, ok := indices[k]; ok {
			rrSets[i].RRs = append(rrSets[i].RRs, rr)
		} else {
			newRRSet := rrSet{
				Hdr: hdr, // hdr.Rdlength may be inconsistent, but we don't care about it.
				RRs: []dns.RR{rr},
			}
			rrSets = append(rrSets, newRRSet)
			indices[k] = len(rrSets) - 1
		}
	}

	for _, set := range rrSets {
		ttl := set.Hdr.Ttl
		for _, rr := range set.RRs {
			hdr := rr.Header()
			if hdr.Ttl != ttl {
				return nil, fmt.Errorf(
					"inconsistent TTLs between %s %s %s records (%d vs. %d); see RFC 2181 section 5.2",
					dns.ClassToString[hdr.Class],
					dns.TypeToString[hdr.Rrtype],
					hdr.Name,
					ttl, hdr.Ttl,
				)
			}
		}
	}

	return rrSets, nil
}
