---
page_title: "Provider: Zonefile"
description: |-
  The Zonefile provider parses the contents of DNS zone files (RFC 1035).
---

# Zonefile Provider

The Zonefile provider parses the contents of DNS [**zone files**][zone file], as
defined by [RFC 1035][RFC 1035], and provides two data sources to read their
contents:

1. `zonefile_records` provides a flat list of the **resource records** (RRs) in
   the zone file, exactly as they were originally defined.
2. `zonefile_record_sets` groups the RRs into **resource record sets** (RRSets)
   by name, class, and type.

[zone file]: https://en.wikipedia.org/wiki/Zone_file
[RFC 1035]: https://datatracker.ietf.org/doc/html/rfc1035

## What is a zone file?

Zone files are text files listing the records for a DNS zone:

```
$ORIGIN example.com.  ; Set the default "base domain" for subsequent records
$TTL 3600             ; Set the default TTL for subsequent records

; A proper zone file (for use with a server like BIND) must include SOA and NS
; records for the zone. However, the kinds of DNS hosts you manage in Terraform
; probably define these records for you, and don't let you control them. The
; zonefile provider won't complain if you omit them.

; Use "@" to represent the apex of the zone (not a subdomain).
@ IN A    10.100.0.10
@ IN A    10.100.0.20
@ IN AAAA fdb6:733c:8b38::100:10
@ IN AAAA fdb6:733c:8b38::100:20

; Some records have multiple fields, like "preference" and "exchange" for MX
; (mail server) records. Note that "IN" is the default class, and any domain not
; fully qualified (with a "." at the end) is assumed to be relative to $ORIGIN.
@ MX 10 mx1.mail.example.
@ MX 20 mx2.mail.example.

; Write dot-separated labels to create subdomains. You can also override $TTL
; for individual records, or specify $TTL again to cover the rest of the file.
staging 60 IN A    10.200.0.10
staging 60 IN AAAA fdb6:733c:8b38::200:10
```

## Why would you use this?

- You find zone files easier to read and write than Terraform resource blocks.
- You're "importing" a zone from your own nameservers into a cloud provider.
- You use multiple DNS providers for business or availability reasons, and want
  to keep your Terraform [DRY][DRY] without inventing a bespoke record schema in
  your Terraform variables.

[DRY]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself

## Why _wouldn't_ you use this?

- You're more comfortable with Terraform syntax than zone files, especially if
  you're not used to things like fully qualifying external hostnames.
- You find dynamic Terraform constructs like `for_each` harder to debug than a
  static list of resource blocks.
- Your provider supports constructs like "alias records" that can't be defined
  in a standard zone file. (However, you could combine dynamic resources driven
  by zone file data with static resource blocks for those special cases.)

## Should you use `records` or `record_sets`?

This depends on the structure of the resource blocks that manage DNS records in
your provider(s) of choice. For example, the following providers use resource
types that manage RRSets rather than individual RRs:

- `hashicorp/aws`
- `hashicorp/azurerm`
- `hashicorp/google`
- `hashicorp/dns`

In contrast, the following providers use resource types that are better driven
by a flat list of RRs:

- `dnsimple/dnsimple`
- `namecheap/namecheap`

If your provider represents DNS record data in a single resource block as a list
rather than a single string, it probably manages an RRSet rather than an
individual RR.
