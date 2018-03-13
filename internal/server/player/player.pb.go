// Code generated by protoc-gen-go. DO NOT EDIT.
// source: player.proto

/*
Package player is a generated protocol buffer package.

It is generated from these files:
	player.proto

It has these top-level messages:
	Empty
	Page
	Progress
	Float
	Int
	Result
	Results
*/
package player

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

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Page struct {
	Page     int64 `protobuf:"varint,1,opt,name=page" json:"page,omitempty"`
	PageSize int64 `protobuf:"varint,2,opt,name=pageSize" json:"pageSize,omitempty"`
}

func (m *Page) Reset()                    { *m = Page{} }
func (m *Page) String() string            { return proto.CompactTextString(m) }
func (*Page) ProtoMessage()               {}
func (*Page) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Page) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *Page) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

type Progress struct {
	N     int64 `protobuf:"varint,1,opt,name=n" json:"n,omitempty"`
	Total int64 `protobuf:"varint,2,opt,name=total" json:"total,omitempty"`
}

func (m *Progress) Reset()                    { *m = Progress{} }
func (m *Progress) String() string            { return proto.CompactTextString(m) }
func (*Progress) ProtoMessage()               {}
func (*Progress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Progress) GetN() int64 {
	if m != nil {
		return m.N
	}
	return 0
}

func (m *Progress) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

type Float struct {
	Value float32 `protobuf:"fixed32,1,opt,name=value" json:"value,omitempty"`
}

func (m *Float) Reset()                    { *m = Float{} }
func (m *Float) String() string            { return proto.CompactTextString(m) }
func (*Float) ProtoMessage()               {}
func (*Float) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Float) GetValue() float32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type Int struct {
	Value int64 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *Int) Reset()                    { *m = Int{} }
func (m *Int) String() string            { return proto.CompactTextString(m) }
func (*Int) ProtoMessage()               {}
func (*Int) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Int) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type Result struct {
	Service  string           `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Path     string           `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Track    *Result_Track    `protobuf:"bytes,3,opt,name=track" json:"track,omitempty"`
	Album    *Result_Album    `protobuf:"bytes,4,opt,name=album" json:"album,omitempty"`
	Artist   *Result_Artist   `protobuf:"bytes,5,opt,name=artist" json:"artist,omitempty"`
	Playlist *Result_Playlist `protobuf:"bytes,6,opt,name=playlist" json:"playlist,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Result) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Result) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Result) GetTrack() *Result_Track {
	if m != nil {
		return m.Track
	}
	return nil
}

func (m *Result) GetAlbum() *Result_Album {
	if m != nil {
		return m.Album
	}
	return nil
}

func (m *Result) GetArtist() *Result_Artist {
	if m != nil {
		return m.Artist
	}
	return nil
}

func (m *Result) GetPlaylist() *Result_Playlist {
	if m != nil {
		return m.Playlist
	}
	return nil
}

type Result_Track struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title    string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Duration int64  `protobuf:"varint,3,opt,name=duration" json:"duration,omitempty"`
}

func (m *Result_Track) Reset()                    { *m = Result_Track{} }
func (m *Result_Track) String() string            { return proto.CompactTextString(m) }
func (*Result_Track) ProtoMessage()               {}
func (*Result_Track) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

func (m *Result_Track) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Result_Track) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Result_Track) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

type Result_Artist struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
}

func (m *Result_Artist) Reset()                    { *m = Result_Artist{} }
func (m *Result_Artist) String() string            { return proto.CompactTextString(m) }
func (*Result_Artist) ProtoMessage()               {}
func (*Result_Artist) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 1} }

func (m *Result_Artist) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Result_Artist) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Result_Album struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
}

func (m *Result_Album) Reset()                    { *m = Result_Album{} }
func (m *Result_Album) String() string            { return proto.CompactTextString(m) }
func (*Result_Album) ProtoMessage()               {}
func (*Result_Album) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 2} }

func (m *Result_Album) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Result_Album) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Result_Playlist struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
}

func (m *Result_Playlist) Reset()                    { *m = Result_Playlist{} }
func (m *Result_Playlist) String() string            { return proto.CompactTextString(m) }
func (*Result_Playlist) ProtoMessage()               {}
func (*Result_Playlist) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 3} }

func (m *Result_Playlist) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Result_Playlist) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Results struct {
	Results []*Result `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *Results) Reset()                    { *m = Results{} }
func (m *Results) String() string            { return proto.CompactTextString(m) }
func (*Results) ProtoMessage()               {}
func (*Results) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Results) GetResults() []*Result {
	if m != nil {
		return m.Results
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "player.Empty")
	proto.RegisterType((*Page)(nil), "player.Page")
	proto.RegisterType((*Progress)(nil), "player.Progress")
	proto.RegisterType((*Float)(nil), "player.Float")
	proto.RegisterType((*Int)(nil), "player.Int")
	proto.RegisterType((*Result)(nil), "player.Result")
	proto.RegisterType((*Result_Track)(nil), "player.Result.Track")
	proto.RegisterType((*Result_Artist)(nil), "player.Result.Artist")
	proto.RegisterType((*Result_Album)(nil), "player.Result.Album")
	proto.RegisterType((*Result_Playlist)(nil), "player.Result.Playlist")
	proto.RegisterType((*Results)(nil), "player.Results")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Player service

type PlayerClient interface {
	Play(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Empty, error)
	PlayAlbum(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Empty, error)
	Volume(ctx context.Context, in *Float, opts ...grpc.CallOption) (*Empty, error)
	Pause(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	FastForward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Queue(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Results, error)
	RemoveFromQueue(ctx context.Context, in *Int, opts ...grpc.CallOption) (*Results, error)
	NextSong(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_NextSongClient, error)
	PlayProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_PlayProgressClient, error)
	DownloadProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_DownloadProgressClient, error)
	History(ctx context.Context, in *Page, opts ...grpc.CallOption) (*Results, error)
}

type playerClient struct {
	cc *grpc.ClientConn
}

func NewPlayerClient(cc *grpc.ClientConn) PlayerClient {
	return &playerClient{cc}
}

func (c *playerClient) Play(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/Play", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) PlayAlbum(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/PlayAlbum", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) Volume(ctx context.Context, in *Float, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/Volume", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) Pause(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/Pause", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) FastForward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/FastForward", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) Queue(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Results, error) {
	out := new(Results)
	err := grpc.Invoke(ctx, "/player.Player/Queue", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) RemoveFromQueue(ctx context.Context, in *Int, opts ...grpc.CallOption) (*Results, error) {
	out := new(Results)
	err := grpc.Invoke(ctx, "/player.Player/RemoveFromQueue", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) NextSong(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_NextSongClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Player_serviceDesc.Streams[0], c.cc, "/player.Player/NextSong", opts...)
	if err != nil {
		return nil, err
	}
	x := &playerNextSongClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Player_NextSongClient interface {
	Recv() (*Result, error)
	grpc.ClientStream
}

type playerNextSongClient struct {
	grpc.ClientStream
}

func (x *playerNextSongClient) Recv() (*Result, error) {
	m := new(Result)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playerClient) PlayProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_PlayProgressClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Player_serviceDesc.Streams[1], c.cc, "/player.Player/PlayProgress", opts...)
	if err != nil {
		return nil, err
	}
	x := &playerPlayProgressClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Player_PlayProgressClient interface {
	Recv() (*Progress, error)
	grpc.ClientStream
}

type playerPlayProgressClient struct {
	grpc.ClientStream
}

func (x *playerPlayProgressClient) Recv() (*Progress, error) {
	m := new(Progress)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playerClient) DownloadProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_DownloadProgressClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Player_serviceDesc.Streams[2], c.cc, "/player.Player/DownloadProgress", opts...)
	if err != nil {
		return nil, err
	}
	x := &playerDownloadProgressClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Player_DownloadProgressClient interface {
	Recv() (*Progress, error)
	grpc.ClientStream
}

type playerDownloadProgressClient struct {
	grpc.ClientStream
}

func (x *playerDownloadProgressClient) Recv() (*Progress, error) {
	m := new(Progress)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playerClient) History(ctx context.Context, in *Page, opts ...grpc.CallOption) (*Results, error) {
	out := new(Results)
	err := grpc.Invoke(ctx, "/player.Player/History", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Player service

type PlayerServer interface {
	Play(context.Context, *Result) (*Empty, error)
	PlayAlbum(context.Context, *Result) (*Empty, error)
	Volume(context.Context, *Float) (*Empty, error)
	Pause(context.Context, *Empty) (*Empty, error)
	FastForward(context.Context, *Empty) (*Empty, error)
	Queue(context.Context, *Empty) (*Results, error)
	RemoveFromQueue(context.Context, *Int) (*Results, error)
	NextSong(*Empty, Player_NextSongServer) error
	PlayProgress(*Empty, Player_PlayProgressServer) error
	DownloadProgress(*Empty, Player_DownloadProgressServer) error
	History(context.Context, *Page) (*Results, error)
}

func RegisterPlayerServer(s *grpc.Server, srv PlayerServer) {
	s.RegisterService(&_Player_serviceDesc, srv)
}

func _Player_Play_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Result)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Play(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Play",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Play(ctx, req.(*Result))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_PlayAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Result)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).PlayAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/PlayAlbum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).PlayAlbum(ctx, req.(*Result))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_Volume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Float)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Volume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Volume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Volume(ctx, req.(*Float))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_Pause_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Pause(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Pause",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Pause(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_FastForward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).FastForward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/FastForward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).FastForward(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_Queue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Queue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Queue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Queue(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_RemoveFromQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).RemoveFromQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/RemoveFromQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).RemoveFromQueue(ctx, req.(*Int))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_NextSong_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlayerServer).NextSong(m, &playerNextSongServer{stream})
}

type Player_NextSongServer interface {
	Send(*Result) error
	grpc.ServerStream
}

type playerNextSongServer struct {
	grpc.ServerStream
}

func (x *playerNextSongServer) Send(m *Result) error {
	return x.ServerStream.SendMsg(m)
}

func _Player_PlayProgress_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlayerServer).PlayProgress(m, &playerPlayProgressServer{stream})
}

type Player_PlayProgressServer interface {
	Send(*Progress) error
	grpc.ServerStream
}

type playerPlayProgressServer struct {
	grpc.ServerStream
}

func (x *playerPlayProgressServer) Send(m *Progress) error {
	return x.ServerStream.SendMsg(m)
}

func _Player_DownloadProgress_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlayerServer).DownloadProgress(m, &playerDownloadProgressServer{stream})
}

type Player_DownloadProgressServer interface {
	Send(*Progress) error
	grpc.ServerStream
}

type playerDownloadProgressServer struct {
	grpc.ServerStream
}

func (x *playerDownloadProgressServer) Send(m *Progress) error {
	return x.ServerStream.SendMsg(m)
}

func _Player_History_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Page)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).History(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/History",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).History(ctx, req.(*Page))
	}
	return interceptor(ctx, in, info, handler)
}

var _Player_serviceDesc = grpc.ServiceDesc{
	ServiceName: "player.Player",
	HandlerType: (*PlayerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Play",
			Handler:    _Player_Play_Handler,
		},
		{
			MethodName: "PlayAlbum",
			Handler:    _Player_PlayAlbum_Handler,
		},
		{
			MethodName: "Volume",
			Handler:    _Player_Volume_Handler,
		},
		{
			MethodName: "Pause",
			Handler:    _Player_Pause_Handler,
		},
		{
			MethodName: "FastForward",
			Handler:    _Player_FastForward_Handler,
		},
		{
			MethodName: "Queue",
			Handler:    _Player_Queue_Handler,
		},
		{
			MethodName: "RemoveFromQueue",
			Handler:    _Player_RemoveFromQueue_Handler,
		},
		{
			MethodName: "History",
			Handler:    _Player_History_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "NextSong",
			Handler:       _Player_NextSong_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PlayProgress",
			Handler:       _Player_PlayProgress_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadProgress",
			Handler:       _Player_DownloadProgress_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "player.proto",
}

func init() { proto.RegisterFile("player.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xd1, 0x6e, 0xd3, 0x3c,
	0x1c, 0xc5, 0x9b, 0xa6, 0x49, 0xdb, 0x7f, 0xfb, 0x6d, 0x93, 0xb5, 0x4f, 0x44, 0x45, 0x48, 0x55,
	0x6e, 0x16, 0x26, 0x1a, 0x8d, 0x56, 0x82, 0x6b, 0x24, 0xa8, 0xe8, 0x0d, 0x0a, 0x1e, 0xe2, 0xde,
	0x5b, 0xac, 0x10, 0xe1, 0xc4, 0x95, 0xed, 0x74, 0x94, 0x27, 0xe0, 0xe9, 0x78, 0x26, 0x64, 0x3b,
	0x89, 0x58, 0x5a, 0x69, 0xe5, 0xaa, 0x3e, 0x3d, 0xbf, 0x63, 0xfb, 0xfc, 0x9b, 0x06, 0xa6, 0x5b,
	0x46, 0xf6, 0x54, 0xc4, 0x5b, 0xc1, 0x15, 0x47, 0xbe, 0x55, 0xe1, 0x10, 0xbc, 0x0f, 0xc5, 0x56,
	0xed, 0xc3, 0x37, 0x30, 0x48, 0x48, 0x46, 0x11, 0x82, 0xc1, 0x96, 0x64, 0x34, 0x70, 0xe6, 0x4e,
	0xe4, 0x62, 0xb3, 0x46, 0x33, 0x18, 0xe9, 0xcf, 0xdb, 0xfc, 0x27, 0x0d, 0xfa, 0xe6, 0xfb, 0x56,
	0x87, 0x31, 0x8c, 0x12, 0xc1, 0x33, 0x41, 0xa5, 0x44, 0x53, 0x70, 0xca, 0x3a, 0xe8, 0x94, 0xe8,
	0x12, 0x3c, 0xc5, 0x15, 0x61, 0x75, 0xc4, 0x8a, 0xf0, 0x05, 0x78, 0x6b, 0xc6, 0x89, 0xd2, 0xf6,
	0x8e, 0xb0, 0xca, 0x9e, 0xd4, 0xc7, 0x56, 0x84, 0xcf, 0xc1, 0xdd, 0x94, 0x1d, 0xd3, 0x6d, 0xcc,
	0xdf, 0x2e, 0xf8, 0x98, 0xca, 0x8a, 0x29, 0x14, 0xc0, 0x50, 0x52, 0xb1, 0xcb, 0xef, 0x2d, 0x32,
	0xc6, 0x8d, 0xb4, 0x05, 0xd4, 0x37, 0x73, 0xea, 0x18, 0x9b, 0x35, 0xba, 0x06, 0x4f, 0x09, 0x72,
	0xff, 0x3d, 0x70, 0xe7, 0x4e, 0x34, 0x59, 0x5e, 0xc6, 0xf5, 0x2c, 0xec, 0x66, 0xf1, 0x17, 0xed,
	0x61, 0x8b, 0x68, 0x96, 0xb0, 0xbb, 0xaa, 0x08, 0x06, 0x47, 0xd9, 0x77, 0xda, 0xc3, 0x16, 0x41,
	0x0b, 0xf0, 0x89, 0x50, 0xb9, 0x54, 0x81, 0x67, 0xe0, 0xff, 0xbb, 0xb0, 0x31, 0x71, 0x0d, 0xa1,
	0x15, 0x8c, 0xb4, 0xcf, 0x74, 0xc0, 0x37, 0x81, 0x67, 0x9d, 0x40, 0x52, 0xdb, 0xb8, 0x05, 0x67,
	0x1b, 0xf0, 0xcc, 0xfd, 0xd0, 0x19, 0xf4, 0xf3, 0xb4, 0x6e, 0xdb, 0xcf, 0x53, 0x33, 0xdf, 0x5c,
	0x31, 0x5a, 0x37, 0xb5, 0x42, 0xff, 0x56, 0x69, 0x25, 0x88, 0xca, 0x79, 0x69, 0xda, 0xba, 0xb8,
	0xd5, 0xb3, 0x18, 0x7c, 0x7b, 0xa3, 0xd3, 0xf6, 0x9a, 0x2d, 0xc0, 0x33, 0x75, 0x4f, 0xc4, 0x6f,
	0x60, 0xd4, 0xdc, 0xff, 0xb4, 0x44, 0xb8, 0x82, 0xa1, 0x2d, 0x2e, 0x51, 0x04, 0x43, 0x61, 0x97,
	0x81, 0x33, 0x77, 0xa3, 0xc9, 0xf2, 0xec, 0xf1, 0x68, 0x70, 0x63, 0x2f, 0x7f, 0x0d, 0xc0, 0x4f,
	0x8c, 0x85, 0xae, 0x60, 0xa0, 0x57, 0xa8, 0xc3, 0xce, 0xfe, 0x6b, 0xb4, 0x7d, 0xb6, 0x7b, 0xe8,
	0x15, 0x8c, 0x35, 0x58, 0xb7, 0x79, 0x8a, 0x8e, 0xc0, 0xff, 0xca, 0x59, 0x55, 0x50, 0xd4, 0x5a,
	0xe6, 0x99, 0x3d, 0x24, 0xaf, 0xc0, 0x4b, 0x48, 0x25, 0xff, 0x02, 0x8d, 0x73, 0x08, 0x2e, 0x60,
	0xb2, 0x26, 0x52, 0xad, 0xb9, 0x78, 0x20, 0x22, 0x7d, 0x12, 0x7f, 0x09, 0xde, 0xe7, 0x8a, 0x56,
	0x07, 0xfb, 0x9e, 0x3f, 0xbe, 0xba, 0x0c, 0x7b, 0xe8, 0x35, 0x9c, 0x63, 0x5a, 0xf0, 0x1d, 0x5d,
	0x0b, 0x5e, 0xd8, 0xd0, 0xa4, 0xa1, 0x36, 0xa5, 0x3a, 0x16, 0x59, 0xc0, 0xe8, 0x13, 0xfd, 0xa1,
	0x6e, 0x79, 0x99, 0x75, 0x0f, 0xe8, 0xcc, 0x26, 0xec, 0xdd, 0x38, 0x68, 0x05, 0x53, 0x3d, 0xbc,
	0xf6, 0x6f, 0xde, 0x89, 0x5c, 0x34, 0xb2, 0x01, 0x4c, 0xe8, 0x2d, 0x5c, 0xbc, 0xe7, 0x0f, 0x25,
	0xe3, 0x24, 0xfd, 0xb7, 0xe0, 0x35, 0x0c, 0x3f, 0xe6, 0x52, 0x71, 0xb1, 0x47, 0xd3, 0x16, 0x20,
	0x19, 0x3d, 0x52, 0xe4, 0xce, 0x37, 0x2f, 0xb3, 0xd5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x24,
	0x44, 0x2b, 0x3c, 0xdc, 0x04, 0x00, 0x00,
}
