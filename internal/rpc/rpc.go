package rpc

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  player.proto

import (
	"github.com/cswank/mcli/internal/schema"
)

func ResultFromPB(r *Result) schema.Result {
	return schema.Result{
		Service:   r.GetService(),
		Path:      r.GetPath(),
		PlayCount: int(r.GetPlaycount()),
		Track:     trackFromPB(r.GetTrack()),
		Album:     albumFromPB(r.GetAlbum()),
		Artist:    artistFromPB(r.GetArtist()),
		Playlist:  playlistFromPB(r.GetPlaylist()),
		Error:     r.GetError(),
	}
}

func trackFromPB(t *Track) schema.Track {
	if t == nil {
		return schema.Track{}
	}
	return schema.Track{
		ID:       t.GetId(),
		Title:    t.GetTitle(),
		Duration: int(t.GetDuration()),
	}
}

func playlistFromPB(p *Playlist) schema.Album {
	if p == nil {
		return schema.Album{}
	}
	return schema.Album{
		ID:    p.GetId(),
		Title: p.GetTitle(),
	}
}

func albumFromPB(a *Album) schema.Album {
	if a == nil {
		return schema.Album{}
	}
	return schema.Album{
		ID:    a.GetId(),
		Title: a.GetTitle(),
	}
}

func artistFromPB(a *Artist) schema.Artist {
	if a == nil {
		return schema.Artist{}
	}
	return schema.Artist{
		ID:   a.GetId(),
		Name: a.GetName(),
	}
}

func pbFromTrack(t schema.Track) *Track {
	return &Track{
		Id:       t.ID,
		Title:    t.Title,
		Duration: int64(t.Duration),
	}
}

func pbFromPlaylist(a schema.Album) *Playlist {
	return &Playlist{
		Id:    a.ID,
		Title: a.Title,
	}
}

func pbFromAlbum(a schema.Album) *Album {
	return &Album{
		Id:    a.ID,
		Title: a.Title,
	}
}

func pbFromArtist(a schema.Artist) *Artist {
	return &Artist{
		Id:   a.ID,
		Name: a.Name,
	}
}

func ResultsFromPB(r *Results) *schema.Results {
	pbRes := r.GetResults()
	res := make([]schema.Result, len(pbRes))
	for i := range pbRes {
		res[i] = ResultFromPB(pbRes[i])
	}
	return &schema.Results{
		Type:    r.GetType(),
		Header:  r.GetHeader(),
		Fmt:     r.GetFmt(),
		Results: res,
		Error:   r.GetError(),
	}
}

func PBFromResults(r *schema.Results) *Results {
	if r == nil {
		return &Results{}
	}
	out := make([]*Result, len(r.Results))
	for i := range r.Results {
		out[i] = PBFromResult(r.Results[i])
	}
	return &Results{
		Header:  r.Header,
		Type:    r.Type,
		Fmt:     r.Fmt,
		Results: out,
		Error:   r.Error,
	}
}

func PBFromResult(r schema.Result) *Result {
	return &Result{
		Service:   r.Service,
		Path:      r.Path,
		Playcount: int64(r.PlayCount),
		Track:     pbFromTrack(r.Track),
		Album:     pbFromAlbum(r.Album),
		Artist:    pbFromArtist(r.Artist),
		Playlist:  pbFromPlaylist(r.Playlist),
		Error:     r.Error,
	}
}

func PBFromProgress(p schema.Progress) *Progress {
	return &Progress{
		N:       int64(p.N),
		Total:   int64(p.Total),
		Payload: p.Payload,
	}
}

func ProgressFromPB(p *Progress) schema.Progress {
	return schema.Progress{
		N:       int(p.GetN()),
		Total:   int(p.GetTotal()),
		Payload: p.Payload,
	}
}
