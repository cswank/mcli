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
	Page     int64  `protobuf:"varint,1,opt,name=page" json:"page,omitempty"`
	PageSize int64  `protobuf:"varint,2,opt,name=pageSize" json:"pageSize,omitempty"`
	Sort     string `protobuf:"bytes,3,opt,name=sort" json:"sort,omitempty"`
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

func (m *Page) GetSort() string {
	if m != nil {
		return m.Sort
	}
	return ""
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
	Service   string           `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Path      string           `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Playcount int64            `protobuf:"varint,3,opt,name=playcount" json:"playcount,omitempty"`
	Track     *Result_Track    `protobuf:"bytes,4,opt,name=track" json:"track,omitempty"`
	Album     *Result_Album    `protobuf:"bytes,5,opt,name=album" json:"album,omitempty"`
	Artist    *Result_Artist   `protobuf:"bytes,6,opt,name=artist" json:"artist,omitempty"`
	Playlist  *Result_Playlist `protobuf:"bytes,7,opt,name=playlist" json:"playlist,omitempty"`
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

func (m *Result) GetPlaycount() int64 {
	if m != nil {
		return m.Playcount
	}
	return 0
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
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
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

func (m *Result_Artist) GetName() string {
	if m != nil {
		return m.Name
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
	Type    string    `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Header  string    `protobuf:"bytes,2,opt,name=header" json:"header,omitempty"`
	Fmt     string    `protobuf:"bytes,3,opt,name=fmt" json:"fmt,omitempty"`
	Results []*Result `protobuf:"bytes,4,rep,name=results" json:"results,omitempty"`
}

func (m *Results) Reset()                    { *m = Results{} }
func (m *Results) String() string            { return proto.CompactTextString(m) }
func (*Results) ProtoMessage()               {}
func (*Results) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Results) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Results) GetHeader() string {
	if m != nil {
		return m.Header
	}
	return ""
}

func (m *Results) GetFmt() string {
	if m != nil {
		return m.Fmt
	}
	return ""
}

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
	PlayAlbum(ctx context.Context, in *Results, opts ...grpc.CallOption) (*Empty, error)
	Volume(ctx context.Context, in *Float, opts ...grpc.CallOption) (*Float, error)
	Pause(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	FastForward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Rewind(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Queue(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Results, error)
	RemoveFromQueue(ctx context.Context, in *Int, opts ...grpc.CallOption) (*Results, error)
	NextSong(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_NextSongClient, error)
	PlayProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_PlayProgressClient, error)
	DownloadProgress(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Player_DownloadProgressClient, error)
	History(ctx context.Context, in *Page, opts ...grpc.CallOption) (*Results, error)
	Done(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
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

func (c *playerClient) PlayAlbum(ctx context.Context, in *Results, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/PlayAlbum", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) Volume(ctx context.Context, in *Float, opts ...grpc.CallOption) (*Float, error) {
	out := new(Float)
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

func (c *playerClient) Rewind(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/Rewind", in, out, c.cc, opts...)
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

func (c *playerClient) Done(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/player.Player/Done", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Player service

type PlayerServer interface {
	Play(context.Context, *Result) (*Empty, error)
	PlayAlbum(context.Context, *Results) (*Empty, error)
	Volume(context.Context, *Float) (*Float, error)
	Pause(context.Context, *Empty) (*Empty, error)
	FastForward(context.Context, *Empty) (*Empty, error)
	Rewind(context.Context, *Empty) (*Empty, error)
	Queue(context.Context, *Empty) (*Results, error)
	RemoveFromQueue(context.Context, *Int) (*Results, error)
	NextSong(*Empty, Player_NextSongServer) error
	PlayProgress(*Empty, Player_PlayProgressServer) error
	DownloadProgress(*Empty, Player_DownloadProgressServer) error
	History(context.Context, *Page) (*Results, error)
	Done(context.Context, *Empty) (*Empty, error)
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
	in := new(Results)
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
		return srv.(PlayerServer).PlayAlbum(ctx, req.(*Results))
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

func _Player_Rewind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Rewind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Rewind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Rewind(ctx, req.(*Empty))
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

func _Player_Done_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).Done(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/player.Player/Done",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).Done(ctx, req.(*Empty))
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
			MethodName: "Rewind",
			Handler:    _Player_Rewind_Handler,
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
		{
			MethodName: "Done",
			Handler:    _Player_Done_Handler,
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
	// 599 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0x86, 0xe3, 0xf8, 0x2f, 0x39, 0xc9, 0xd7, 0x56, 0x47, 0xfd, 0xc0, 0x32, 0x20, 0x45, 0x5e,
	0x50, 0x53, 0x91, 0xa8, 0xb4, 0x0b, 0xd6, 0x48, 0x25, 0xa2, 0x2c, 0x50, 0x98, 0x22, 0xf6, 0xd3,
	0x7a, 0x48, 0x2d, 0x6c, 0x4f, 0x18, 0x8f, 0x5b, 0xc2, 0xad, 0x72, 0x07, 0x5c, 0x05, 0x9a, 0x19,
	0xdb, 0x50, 0x27, 0x52, 0xca, 0x2a, 0xe7, 0xf5, 0xfb, 0x9c, 0x9f, 0x39, 0xce, 0x18, 0xc6, 0xab,
	0x8c, 0xae, 0x99, 0x98, 0xad, 0x04, 0x97, 0x1c, 0x3d, 0xa3, 0x22, 0x1f, 0xdc, 0xb7, 0xf9, 0x4a,
	0xae, 0xa3, 0xf7, 0xe0, 0x2c, 0xe8, 0x92, 0x21, 0x82, 0xb3, 0xa2, 0x4b, 0x16, 0x58, 0x13, 0x2b,
	0xb6, 0x89, 0x8e, 0x31, 0x84, 0x81, 0xfa, 0xbd, 0x4c, 0x7f, 0xb0, 0xa0, 0xaf, 0x9f, 0xb7, 0x5a,
	0xf1, 0x25, 0x17, 0x32, 0xb0, 0x27, 0x56, 0x3c, 0x24, 0x3a, 0x8e, 0x66, 0x30, 0x58, 0x08, 0xbe,
	0x14, 0xac, 0x2c, 0x71, 0x0c, 0x56, 0x51, 0x17, 0xb3, 0x0a, 0x3c, 0x04, 0x57, 0x72, 0x49, 0xb3,
	0xba, 0x8c, 0x11, 0xd1, 0x33, 0x70, 0xe7, 0x19, 0xa7, 0x52, 0xd9, 0xb7, 0x34, 0xab, 0x4c, 0xf7,
	0x3e, 0x31, 0x22, 0x7a, 0x02, 0xf6, 0x45, 0xd1, 0x31, 0xed, 0xc6, 0xfc, 0x65, 0x83, 0x47, 0x58,
	0x59, 0x65, 0x12, 0x03, 0xf0, 0x4b, 0x26, 0x6e, 0xd3, 0x6b, 0x83, 0x0c, 0x49, 0x23, 0xcd, 0xa1,
	0xe4, 0x8d, 0xee, 0x3a, 0x24, 0x3a, 0xc6, 0xa7, 0x30, 0x54, 0x3b, 0xb8, 0xe6, 0x55, 0x61, 0xa6,
	0xb7, 0xc9, 0x9f, 0x07, 0x78, 0x0c, 0xae, 0x14, 0xf4, 0xfa, 0x6b, 0xe0, 0x4c, 0xac, 0x78, 0x74,
	0x7a, 0x38, 0xab, 0xb7, 0x67, 0x5a, 0xcd, 0x3e, 0x29, 0x8f, 0x18, 0x44, 0xb1, 0x34, 0xbb, 0xaa,
	0xf2, 0xc0, 0xdd, 0xca, 0xbe, 0x51, 0x1e, 0x31, 0x08, 0x4e, 0xc1, 0xa3, 0x42, 0xa6, 0xa5, 0x0c,
	0x3c, 0x0d, 0xff, 0xdf, 0x85, 0xb5, 0x49, 0x6a, 0x08, 0xcf, 0x60, 0xa0, 0xfc, 0x4c, 0x25, 0xf8,
	0x3a, 0xe1, 0x71, 0x27, 0x61, 0x51, 0xdb, 0xa4, 0x05, 0xc3, 0x0b, 0x70, 0xf5, 0x7c, 0xb8, 0x07,
	0xfd, 0x34, 0xa9, 0x77, 0xd1, 0x4f, 0x13, 0xbd, 0xfd, 0x54, 0x66, 0xac, 0xde, 0x83, 0x11, 0xea,
	0xed, 0x26, 0x95, 0xa0, 0x32, 0xe5, 0x45, 0xbd, 0x87, 0x56, 0x87, 0x2f, 0xc1, 0x33, 0x13, 0x6d,
	0xd4, 0x42, 0x70, 0x0a, 0x9a, 0x37, 0xa5, 0x74, 0x1c, 0x4e, 0xc1, 0xd5, 0x87, 0x7d, 0x58, 0xe3,
	0xf0, 0x04, 0x06, 0xcd, 0xf4, 0x0f, 0xcb, 0x88, 0xbe, 0x81, 0x6f, 0x8e, 0x5d, 0xaa, 0xfe, 0x72,
	0xbd, 0x6a, 0xde, 0xb4, 0x8e, 0xf1, 0x11, 0x78, 0x37, 0x8c, 0x26, 0x4c, 0xd4, 0x59, 0xb5, 0xc2,
	0x03, 0xb0, 0xbf, 0xe4, 0xcd, 0x5f, 0x54, 0x85, 0x18, 0x83, 0x2f, 0x4c, 0xa1, 0xc0, 0x99, 0xd8,
	0xf1, 0xe8, 0x74, 0xef, 0xfe, 0x5a, 0x49, 0x63, 0x9f, 0xfe, 0x74, 0xc0, 0x5b, 0x68, 0x0b, 0x8f,
	0xc0, 0x51, 0x11, 0x76, 0xd8, 0xf0, 0xbf, 0x46, 0x9b, 0x9b, 0xd4, 0xc3, 0x29, 0x0c, 0x15, 0x68,
	0x76, 0xb1, 0x7f, 0x9f, 0x2e, 0x37, 0xf1, 0x18, 0xbc, 0xcf, 0x3c, 0xab, 0x72, 0x86, 0xad, 0xa5,
	0xaf, 0x43, 0x78, 0x5f, 0x46, 0x3d, 0x3c, 0x02, 0x77, 0x41, 0xab, 0xf2, 0x2f, 0x50, 0xd7, 0xd8,
	0x36, 0xc1, 0x68, 0x4e, 0x4b, 0x39, 0xe7, 0xe2, 0x8e, 0x8a, 0x64, 0x27, 0x1e, 0xab, 0x3b, 0x74,
	0x97, 0x16, 0xbb, 0xc9, 0x17, 0xe0, 0x7e, 0xac, 0x58, 0xb5, 0x31, 0x41, 0xf7, 0x94, 0x51, 0x0f,
	0x5f, 0xc1, 0x3e, 0x61, 0x39, 0xbf, 0x65, 0x73, 0xc1, 0x73, 0x93, 0x34, 0x6a, 0xa8, 0x8b, 0x42,
	0x6e, 0x4b, 0x99, 0xc2, 0xe0, 0x03, 0xfb, 0x2e, 0x2f, 0x79, 0xb1, 0xec, 0x36, 0xe8, 0x2c, 0x3d,
	0xea, 0x9d, 0x58, 0x78, 0x06, 0x63, 0xb5, 0xe7, 0xf6, 0x5b, 0xd3, 0x49, 0x39, 0x68, 0x64, 0x03,
	0xe8, 0xa4, 0xd7, 0x70, 0x70, 0xce, 0xef, 0x8a, 0x8c, 0xd3, 0xe4, 0xdf, 0x12, 0x8f, 0xc1, 0x7f,
	0x97, 0x96, 0x92, 0x8b, 0x35, 0x8e, 0x5b, 0x80, 0x2e, 0xd9, 0xb6, 0x83, 0x3c, 0x07, 0xe7, 0x9c,
	0x17, 0x3b, 0xdf, 0xd3, 0x95, 0xa7, 0xbf, 0xc6, 0x67, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xea,
	0x25, 0xaa, 0x06, 0x9d, 0x05, 0x00, 0x00,
}
