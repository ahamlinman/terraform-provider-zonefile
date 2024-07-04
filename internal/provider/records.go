package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/miekg/dns"
	"github.com/samber/lo"
)

var _ datasource.DataSource = &RecordsDataSource{}

type RecordsDataSource struct{}

func NewRecordsDataSource() datasource.DataSource {
	return &RecordsDataSource{}
}

func (d *RecordsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_records"
}

func (d *RecordsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read a DNS zone file and return a flat list of resource records.",
		Attributes:  schemaRecordsModel,
	}
}

func (d *RecordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	origin := data.Origin.ValueString()
	rrs, err := readZone(origin, data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Invalid zone file", err.Error()))
		return
	}

	data.Records = lo.Map(rrs, func(rr dns.RR, _ int) RecordsItemModel {
		hdr := rr.Header()
		return RecordsItemModel{
			Name:  nameModelValue(hdr.Name, origin),
			FQDN:  types.StringValue(hdr.Name),
			Class: types.StringValue(dns.ClassToString[hdr.Class]),
			Type:  types.StringValue(dns.TypeToString[hdr.Rrtype]),
			TTL:   types.Int64Value(int64(hdr.Ttl)),

			Data: rdataModelValue(rr),
			MX:   mxModelValue(rr),
			SRV:  srvModelValue(rr),
			TXT:  txtModelValue(rr),
		}
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
