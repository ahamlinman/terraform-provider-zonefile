package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/miekg/dns"
	"github.com/samber/lo"
)

var _ datasource.DataSource = &RecordSetsDataSource{}

type RecordSetsDataSource struct{}

func NewRecordSetsDataSource() datasource.DataSource {
	return &RecordSetsDataSource{}
}

func (d *RecordSetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record_sets"
}

func (d *RecordSetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read a DNS zone file and return RRSets: resource records grouped by name, class, and type.",
		Attributes:  schemaRecordSetsModel,
	}
}

func (d *RecordSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordSetsModel
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

	rrSets, err := groupRRs(rrs)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Can't group some RRs into RRSets", err.Error()))
		return
	}

	tryList := func(list basetypes.ListValue, diag diag.Diagnostics) basetypes.ListValue {
		resp.Diagnostics.Append(diag...)
		return list
	}

	data.RRSets = lo.Map(rrSets, func(set rrSet, _ int) RecordSetsItemModel {
		hdr := set.Hdr
		return RecordSetsItemModel{
			Name:  nameModelValue(hdr.Name, origin),
			FQDN:  types.StringValue(hdr.Name),
			Class: types.StringValue(dns.ClassToString[hdr.Class]),
			Type:  types.StringValue(dns.TypeToString[hdr.Rrtype]),
			TTL:   types.Int64Value(int64(hdr.Ttl)),

			Data: tryList(types.ListValue(types.StringType,
				lo.Map(set.RRs, func(rr dns.RR, _ int) attr.Value {
					return rdataModelValue(rr)
				}))),

			MX: lo.Ternary(
				hdr.Rrtype != dns.TypeMX,
				types.ListNull(attributeObjectMXModel.Type()),
				tryList(types.ListValueFrom(ctx, attributeObjectMXModel.Type(),
					lo.Map(set.RRs, func(rr dns.RR, _ int) *RecordsMXModel {
						return mxModelValue(rr)
					})))),

			SRV: lo.Ternary(
				hdr.Rrtype != dns.TypeSRV,
				types.ListNull(attributeObjectSRVModel.Type()),
				tryList(types.ListValueFrom(ctx, attributeObjectSRVModel.Type(),
					lo.Map(set.RRs, func(rr dns.RR, _ int) *RecordsSRVModel {
						return srvModelValue(rr)
					})))),

			TXT: lo.Ternary(
				hdr.Rrtype != dns.TypeTXT,
				types.ListNull(types.StringType),
				tryList(types.ListValue(types.StringType,
					lo.Map(set.RRs, func(rr dns.RR, _ int) attr.Value {
						return txtModelValue(rr)
					})))),
		}
	})

	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}
