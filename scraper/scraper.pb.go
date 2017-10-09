// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scraper.proto

/*
Package scraper is a generated protocol buffer package.

It is generated from these files:
	scraper.proto

It has these top-level messages:
	ScrapeReq
	ScrapeResp
*/
package scraper

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ScrapeReq struct {
	Id  string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Url string `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
}

func (m *ScrapeReq) Reset()                    { *m = ScrapeReq{} }
func (m *ScrapeReq) String() string            { return proto.CompactTextString(m) }
func (*ScrapeReq) ProtoMessage()               {}
func (*ScrapeReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ScrapeReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ScrapeReq) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type ScrapeResp struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ArchiveUrl string `protobuf:"bytes,2,opt,name=archive_url,json=archiveUrl" json:"archive_url,omitempty"`
}

func (m *ScrapeResp) Reset()                    { *m = ScrapeResp{} }
func (m *ScrapeResp) String() string            { return proto.CompactTextString(m) }
func (*ScrapeResp) ProtoMessage()               {}
func (*ScrapeResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ScrapeResp) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ScrapeResp) GetArchiveUrl() string {
	if m != nil {
		return m.ArchiveUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*ScrapeReq)(nil), "scraper.ScrapeReq")
	proto.RegisterType((*ScrapeResp)(nil), "scraper.ScrapeResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ScraperService service

type ScraperServiceClient interface {
	Scrape(ctx context.Context, in *ScrapeReq, opts ...grpc.CallOption) (*ScrapeResp, error)
}

type scraperServiceClient struct {
	cc *grpc.ClientConn
}

func NewScraperServiceClient(cc *grpc.ClientConn) ScraperServiceClient {
	return &scraperServiceClient{cc}
}

func (c *scraperServiceClient) Scrape(ctx context.Context, in *ScrapeReq, opts ...grpc.CallOption) (*ScrapeResp, error) {
	out := new(ScrapeResp)
	err := grpc.Invoke(ctx, "/scraper.ScraperService/Scrape", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ScraperService service

type ScraperServiceServer interface {
	Scrape(context.Context, *ScrapeReq) (*ScrapeResp, error)
}

func RegisterScraperServiceServer(s *grpc.Server, srv ScraperServiceServer) {
	s.RegisterService(&_ScraperService_serviceDesc, srv)
}

func _ScraperService_Scrape_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScrapeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScraperServiceServer).Scrape(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scraper.ScraperService/Scrape",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScraperServiceServer).Scrape(ctx, req.(*ScrapeReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _ScraperService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "scraper.ScraperService",
	HandlerType: (*ScraperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Scrape",
			Handler:    _ScraperService_Scrape_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scraper.proto",
}

func init() { proto.RegisterFile("scraper.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 151 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2e, 0x4a,
	0x2c, 0x48, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x74, 0xb9,
	0x38, 0x83, 0xc1, 0xcc, 0xa0, 0xd4, 0x42, 0x21, 0x3e, 0x2e, 0xa6, 0xcc, 0x14, 0x09, 0x46, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0xa6, 0xcc, 0x14, 0x21, 0x01, 0x2e, 0xe6, 0xd2, 0xa2, 0x1c, 0x09, 0x26,
	0xb0, 0x00, 0x88, 0xa9, 0x64, 0xcb, 0xc5, 0x05, 0x53, 0x5e, 0x5c, 0x80, 0xa1, 0x5e, 0x9e, 0x8b,
	0x3b, 0xb1, 0x28, 0x39, 0x23, 0xb3, 0x2c, 0x35, 0x1e, 0xa1, 0x8f, 0x0b, 0x2a, 0x14, 0x5a, 0x94,
	0x63, 0xe4, 0xca, 0xc5, 0x07, 0xd1, 0x5e, 0x14, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0x64,
	0xcc, 0xc5, 0x06, 0x11, 0x11, 0x12, 0xd2, 0x83, 0x39, 0x11, 0xee, 0x20, 0x29, 0x61, 0x0c, 0xb1,
	0xe2, 0x02, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x27, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x31,
	0x92, 0x5f, 0xaf, 0xd5, 0x00, 0x00, 0x00,
}
