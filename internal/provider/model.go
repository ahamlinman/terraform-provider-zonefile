package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/miekg/dns"
	"github.com/samber/lo"
)

// RecordsModel represents the entire "zonefile_records" data source.
type RecordsModel struct {
	Content types.String `tfsdk:"content"`
	Origin  types.String `tfsdk:"origin"`

	Records []RecordsItemModel `tfsdk:"records"`
}

// RecordSetsModel represents the entire "zonefile_record_sets" data source.
type RecordSetsModel struct {
	Content types.String `tfsdk:"content"`
	Origin  types.String `tfsdk:"origin"`

	RRSets []RecordSetsItemModel `tfsdk:"rrsets"`
}

var schemaModelHead = map[string]schema.Attribute{
	"content": schema.StringAttribute{
		Required: true,
		Description: ("The entire zone file as a string. " +
			"You can read this from disk with the file(…) function or local_file data source."),
	},
	"origin": schema.StringAttribute{
		Optional: true,
		Description: ("The origin for relative record names in the file, " +
			"equivalent to an $ORIGIN directive at the top of the file. " +
			"If set, the provider will populate the \"name\" field of records. " +
			"Otherwise, only \"fqdn\" will be available even if the file includes an $ORIGIN directive."),
	},
}

var schemaRecordsModel = lo.Assign(
	schemaModelHead,
	map[string]schema.Attribute{
		"records": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{Attributes: schemaRecordsItemModel},
			Computed:     true,
			Description:  "The zone file's resource records.",
		},
	})

var schemaRecordSetsModel = lo.Assign(
	schemaModelHead,
	map[string]schema.Attribute{
		"rrsets": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{Attributes: schemaRecordSetsItemModel},
			Computed:     true,
			Description: ("The zone file's resource records grouped by name, class, and type. " +
				"Unlike the records data source, this data source will fail with an error if " +
				"any RRs in an RRSet have inconsistent TTLs (per RFC 2181 section 5.2)."),
		},
	})

// RecordsItemModel represents each element in the "records" list of the
// "zonefile_records" data source.
type RecordsItemModel struct {
	Name  types.String `tfsdk:"name"`
	FQDN  types.String `tfsdk:"fqdn"`
	Class types.String `tfsdk:"class"`
	Type  types.String `tfsdk:"type"`
	TTL   types.Int64  `tfsdk:"ttl"`

	Data types.String     `tfsdk:"data"`
	MX   *RecordsMXModel  `tfsdk:"mx"`
	SRV  *RecordsSRVModel `tfsdk:"srv"`
	TXT  types.String     `tfsdk:"txt"`
}

// RecordSetsItemModel represents each element in the "rrsets" list of the
// "zonefile_record_sets" data source.
type RecordSetsItemModel struct {
	Name  types.String `tfsdk:"name"`
	FQDN  types.String `tfsdk:"fqdn"`
	Class types.String `tfsdk:"class"`
	Type  types.String `tfsdk:"type"`
	TTL   types.Int64  `tfsdk:"ttl"`

	Data types.List `tfsdk:"data"`
	MX   types.List `tfsdk:"mx"`
	SRV  types.List `tfsdk:"srv"`
	TXT  types.List `tfsdk:"txt"`
}

var schemaItemModelHead = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Computed: true,
		Description: ("The record's name relative to the origin in the data source configuration. " +
			"This will be null for the zone apex (\"@\" in a zone file), or if the data source configuration " +
			"does not specify an origin (even if the zone file includes an $ORIGIN directive)."),
	},
	"fqdn": schema.StringAttribute{
		Computed: true,
		Description: ("The record's fully qualified name. " +
			"Unlike \"name\", this includes the effect of any $ORIGIN directives and ends with a trailing dot."),
	},
	"class": schema.StringAttribute{
		Computed:    true,
		Description: "The record's class, usually IN (Internet).",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "The record's type as an uppercase string: A, AAAA, CNAME, TXT, etc.",
	},
	"ttl": schema.Int64Attribute{
		Computed: true,
		Description: ("The record's TTL as an integer number of seconds. " +
			"This includes the effect of any $TTL directives in the zone file."),
	},
}

var schemaRecordsItemModel = lo.Assign(
	schemaItemModelHead,
	map[string]schema.Attribute{
		"data": schema.StringAttribute{
			Computed: true,
			Description: ("The record's data (RDATA) in its canonical presentation format " +
				"(that is, how you might write it in a zone file). " +
				"The provider parses the fields of select record types like MX and SRV, " +
				"which is more robust than pulling them out of the RDATA string."),
		},
		"mx": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "The parsed fields of an MX record, or null if this isn't an MX record.",
			Attributes:  schemaRecordsMXModel,
		},
		"srv": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "The parsed fields of an SRV record, or null if this isn't an SRV record.",
			Attributes:  schemaRecordsSRVModel,
		},
		"txt": schema.StringAttribute{
			Computed: true,
			Description: ("The concatenation of multiple strings in the TXT record, or null if this isn't a TXT record. " +
				"Individual strings in a TXT RDATA section must be 255 characters or less, " +
				"but a single TXT RR can define multiple logically concatenated strings. " +
				"Your provider may require special handling if this value is longer than 255 characters."),
		},
	},
)

var schemaRecordSetsItemModel = lo.Assign(
	schemaItemModelHead,
	map[string]schema.Attribute{
		"data": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: ("The record data (RDATA) for each RR in canonical presentation format " +
				"(that is, how you might write it in a zone file). " +
				"The provider parses the fields of select record types like MX and SRV, " +
				"which is more robust than pulling them out of the RDATA strings."),
		},
		"mx": schema.ListNestedAttribute{
			NestedObject: attributeObjectMXModel,
			Computed:     true,
			Description:  "The parsed fields of MX records, or null if this isn't an MX RRSet.",
		},
		"srv": schema.ListNestedAttribute{
			NestedObject: attributeObjectSRVModel,
			Computed:     true,
			Description:  "The parsed fields of SRV records, or null if this isn't an SRV RRSet.",
		},
		"txt": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: ("The concatenation of multiple strings in each TXT record, or null if this isn't a TXT RRSet. " +
				"Individual strings in a TXT RDATA section must be 255 characters or less, " +
				"but a single TXT RR can define multiple logically concatenated strings. " +
				"Your provider may require special handling if any of these values are longer than 255 characters."),
		},
	},
)

func nameModelValue(fqdn, origin string) types.String {
	if origin == "" {
		return types.StringNull()
	}
	name := strings.TrimSuffix(fqdn, dns.Fqdn(origin))
	name = strings.TrimSuffix(name, ".")
	if name == "" {
		return types.StringNull()
	}
	return types.StringValue(name)
}

func rdataModelValue(rr dns.RR) types.String {
	return types.StringValue(strings.TrimPrefix(rr.String(), rr.Header().String()))
}

func txtModelValue(rr dns.RR) types.String {
	if txt, ok := rr.(*dns.TXT); ok {
		return types.StringValue(strings.Join(txt.Txt, ""))
	}
	return types.StringNull()
}

// RecordsMXModel represents the parsed fields of MX records exposed through
// either data source.
type RecordsMXModel struct {
	Preference types.Int64  `tfsdk:"preference"`
	Exchange   types.String `tfsdk:"exchange"`
}

var (
	attributeObjectMXModel = schema.NestedAttributeObject{Attributes: schemaRecordsMXModel}
	schemaRecordsMXModel   = map[string]schema.Attribute{
		"preference": schema.Int64Attribute{
			Computed:    true,
			Description: "The preference given to this MX record among others at the same owner.",
		},
		"exchange": schema.StringAttribute{
			Computed:    true,
			Description: "The domain name of the host acting as a mail exchange for the owner name.",
		},
	}
)

func mxModelValue(rr dns.RR) *RecordsMXModel {
	if mx, ok := rr.(*dns.MX); ok {
		return &RecordsMXModel{
			Preference: types.Int64Value(int64(mx.Preference)),
			Exchange:   types.StringValue(mx.Mx),
		}
	}
	return nil
}

// RecordsSRVModel represents the parsed fields of SRV records exposed through
// either data source.
type RecordsSRVModel struct {
	Priority types.Int64  `tfsdk:"priority"`
	Weight   types.Int64  `tfsdk:"weight"`
	Port     types.Int64  `tfsdk:"port"`
	Target   types.String `tfsdk:"target"`
}

var (
	attributeObjectSRVModel = schema.NestedAttributeObject{Attributes: schemaRecordsSRVModel}
	schemaRecordsSRVModel   = map[string]schema.Attribute{
		"priority": schema.Int64Attribute{
			Computed:    true,
			Description: "The priority of the target host, with lower priorities taking precedence.",
		},
		"weight": schema.Int64Attribute{
			Computed:    true,
			Description: "The relative weight for a target with the same priority as another.",
		},
		"port": schema.Int64Attribute{
			Computed:    true,
			Description: "The port on this target host of this service.",
		},
		"target": schema.StringAttribute{
			Computed:    true,
			Description: "The domain name of the target host.",
		},
	}
)

func srvModelValue(rr dns.RR) *RecordsSRVModel {
	if srv, ok := rr.(*dns.SRV); ok {
		return &RecordsSRVModel{
			Priority: types.Int64Value(int64(srv.Priority)),
			Weight:   types.Int64Value(int64(srv.Weight)),
			Port:     types.Int64Value(int64(srv.Port)),
			Target:   types.StringValue(srv.Target),
		}
	}
	return nil
}
