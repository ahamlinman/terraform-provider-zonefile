package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testOrigin   = "main.test."
	testZonefile = `
@ 1800 IN A 10.100.0.10
@ 1800 IN A 10.200.0.20

@ 3600 IN MX 10 mx1.mail.test.
@ 3600 IN MX 20 mx2.mail.test.

srv 1800 IN SRV 1 1 443 app1.app.test.
srv 1800 IN SRV 2 1 443 app2.app.test.

txt 300 IN TXT "first"
txt 300 IN TXT "and" " " "second"
`
)

var (
	eq   = resource.TestCheckResourceAttr
	null = resource.TestCheckNoResourceAttr
)

func TestZonefileDataSources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"zonefile": providerserver.NewProtocol6WithError(New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "zonefile_record_sets" "main" {
						origin  = %q
						content = %q
					}
					data "zonefile_records" "main" {
						origin  = %q
						content = %q
					}`,
					testOrigin, testZonefile,
					testOrigin, testZonefile),
				Check: resource.ComposeAggregateTestCheckFunc(
					eq("data.zonefile_record_sets.main", "rrsets.#", "4"),
					eq("data.zonefile_records.main", "records.#", "8"),

					eq("data.zonefile_record_sets.main", "rrsets.0.fqdn", "main.test."),
					eq("data.zonefile_record_sets.main", "rrsets.0.class", "IN"),
					eq("data.zonefile_record_sets.main", "rrsets.0.type", "A"),
					eq("data.zonefile_record_sets.main", "rrsets.0.ttl", "1800"),
					eq("data.zonefile_record_sets.main", "rrsets.0.data.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.0.data.0", "10.100.0.10"),
					eq("data.zonefile_record_sets.main", "rrsets.0.data.1", "10.200.0.20"),
					eq("data.zonefile_records.main", "records.0.fqdn", "main.test."),
					eq("data.zonefile_records.main", "records.0.class", "IN"),
					eq("data.zonefile_records.main", "records.0.type", "A"),
					eq("data.zonefile_records.main", "records.0.ttl", "1800"),
					eq("data.zonefile_records.main", "records.0.data", "10.100.0.10"),
					eq("data.zonefile_records.main", "records.1.fqdn", "main.test."),
					eq("data.zonefile_records.main", "records.1.class", "IN"),
					eq("data.zonefile_records.main", "records.1.type", "A"),
					eq("data.zonefile_records.main", "records.1.ttl", "1800"),
					eq("data.zonefile_records.main", "records.1.data", "10.200.0.20"),
					null("data.zonefile_record_sets.main", "rrsets.0.name"),
					null("data.zonefile_record_sets.main", "rrsets.0.mx"),
					null("data.zonefile_record_sets.main", "rrsets.0.srv"),
					null("data.zonefile_record_sets.main", "rrsets.0.txt"),
					null("data.zonefile_records.main", "records.0.name"),
					null("data.zonefile_records.main", "records.0.mx"),
					null("data.zonefile_records.main", "records.0.srv"),
					null("data.zonefile_records.main", "records.0.txt"),
					null("data.zonefile_records.main", "records.1.name"),
					null("data.zonefile_records.main", "records.1.mx"),
					null("data.zonefile_records.main", "records.1.srv"),
					null("data.zonefile_records.main", "records.1.txt"),

					eq("data.zonefile_record_sets.main", "rrsets.1.fqdn", "main.test."),
					eq("data.zonefile_record_sets.main", "rrsets.1.class", "IN"),
					eq("data.zonefile_record_sets.main", "rrsets.1.type", "MX"),
					eq("data.zonefile_record_sets.main", "rrsets.1.ttl", "3600"),
					eq("data.zonefile_record_sets.main", "rrsets.1.data.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.1.data.0", "10 mx1.mail.test."),
					eq("data.zonefile_record_sets.main", "rrsets.1.data.1", "20 mx2.mail.test."),
					eq("data.zonefile_record_sets.main", "rrsets.1.mx.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.1.mx.0.preference", "10"),
					eq("data.zonefile_record_sets.main", "rrsets.1.mx.0.exchange", "mx1.mail.test."),
					eq("data.zonefile_record_sets.main", "rrsets.1.mx.1.preference", "20"),
					eq("data.zonefile_record_sets.main", "rrsets.1.mx.1.exchange", "mx2.mail.test."),
					eq("data.zonefile_records.main", "records.2.fqdn", "main.test."),
					eq("data.zonefile_records.main", "records.2.class", "IN"),
					eq("data.zonefile_records.main", "records.2.type", "MX"),
					eq("data.zonefile_records.main", "records.2.ttl", "3600"),
					eq("data.zonefile_records.main", "records.2.data", "10 mx1.mail.test."),
					eq("data.zonefile_records.main", "records.2.mx.preference", "10"),
					eq("data.zonefile_records.main", "records.2.mx.exchange", "mx1.mail.test."),
					eq("data.zonefile_records.main", "records.3.fqdn", "main.test."),
					eq("data.zonefile_records.main", "records.3.class", "IN"),
					eq("data.zonefile_records.main", "records.3.type", "MX"),
					eq("data.zonefile_records.main", "records.3.ttl", "3600"),
					eq("data.zonefile_records.main", "records.3.data", "20 mx2.mail.test."),
					eq("data.zonefile_records.main", "records.3.mx.preference", "20"),
					eq("data.zonefile_records.main", "records.3.mx.exchange", "mx2.mail.test."),

					eq("data.zonefile_record_sets.main", "rrsets.2.fqdn", "srv.main.test."),
					eq("data.zonefile_record_sets.main", "rrsets.2.name", "srv"),
					eq("data.zonefile_record_sets.main", "rrsets.2.class", "IN"),
					eq("data.zonefile_record_sets.main", "rrsets.2.type", "SRV"),
					eq("data.zonefile_record_sets.main", "rrsets.2.ttl", "1800"),
					eq("data.zonefile_record_sets.main", "rrsets.2.data.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.2.data.0", "1 1 443 app1.app.test."),
					eq("data.zonefile_record_sets.main", "rrsets.2.data.1", "2 1 443 app2.app.test."),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.0.priority", "1"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.0.weight", "1"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.0.port", "443"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.0.target", "app1.app.test."),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.1.priority", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.1.weight", "1"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.1.port", "443"),
					eq("data.zonefile_record_sets.main", "rrsets.2.srv.1.target", "app2.app.test."),
					eq("data.zonefile_records.main", "records.4.fqdn", "srv.main.test."),
					eq("data.zonefile_records.main", "records.4.name", "srv"),
					eq("data.zonefile_records.main", "records.4.class", "IN"),
					eq("data.zonefile_records.main", "records.4.type", "SRV"),
					eq("data.zonefile_records.main", "records.4.data", "1 1 443 app1.app.test."),
					eq("data.zonefile_records.main", "records.4.srv.priority", "1"),
					eq("data.zonefile_records.main", "records.4.srv.weight", "1"),
					eq("data.zonefile_records.main", "records.4.srv.port", "443"),
					eq("data.zonefile_records.main", "records.4.srv.target", "app1.app.test."),
					eq("data.zonefile_records.main", "records.5.fqdn", "srv.main.test."),
					eq("data.zonefile_records.main", "records.5.name", "srv"),
					eq("data.zonefile_records.main", "records.5.class", "IN"),
					eq("data.zonefile_records.main", "records.5.type", "SRV"),
					eq("data.zonefile_records.main", "records.5.data", "2 1 443 app2.app.test."),
					eq("data.zonefile_records.main", "records.5.srv.priority", "2"),
					eq("data.zonefile_records.main", "records.5.srv.weight", "1"),
					eq("data.zonefile_records.main", "records.5.srv.port", "443"),
					eq("data.zonefile_records.main", "records.5.srv.target", "app2.app.test."),

					eq("data.zonefile_record_sets.main", "rrsets.3.fqdn", "txt.main.test."),
					eq("data.zonefile_record_sets.main", "rrsets.3.name", "txt"),
					eq("data.zonefile_record_sets.main", "rrsets.3.class", "IN"),
					eq("data.zonefile_record_sets.main", "rrsets.3.type", "TXT"),
					eq("data.zonefile_record_sets.main", "rrsets.3.ttl", "300"),
					eq("data.zonefile_record_sets.main", "rrsets.3.data.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.3.data.0", `"first"`),
					eq("data.zonefile_record_sets.main", "rrsets.3.data.1", `"and" " " "second"`),
					eq("data.zonefile_record_sets.main", "rrsets.3.txt.#", "2"),
					eq("data.zonefile_record_sets.main", "rrsets.3.txt.0", "first"),
					eq("data.zonefile_record_sets.main", "rrsets.3.txt.1", "and second"),
					eq("data.zonefile_records.main", "records.6.fqdn", "txt.main.test."),
					eq("data.zonefile_records.main", "records.6.name", "txt"),
					eq("data.zonefile_records.main", "records.6.class", "IN"),
					eq("data.zonefile_records.main", "records.6.type", "TXT"),
					eq("data.zonefile_records.main", "records.6.ttl", "300"),
					eq("data.zonefile_records.main", "records.6.data", `"first"`),
					eq("data.zonefile_records.main", "records.6.txt", "first"),
					eq("data.zonefile_records.main", "records.7.fqdn", "txt.main.test."),
					eq("data.zonefile_records.main", "records.7.name", "txt"),
					eq("data.zonefile_records.main", "records.7.class", "IN"),
					eq("data.zonefile_records.main", "records.7.type", "TXT"),
					eq("data.zonefile_records.main", "records.7.ttl", "300"),
					eq("data.zonefile_records.main", "records.7.data", `"and" " " "second"`),
					eq("data.zonefile_records.main", "records.7.txt", "and second"),
				),
			},
		},
	})
}
