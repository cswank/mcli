package download

import (
	"io"

	"bitbucket.org/cswank/mcli/internal/schema"
)

type Downloader interface {
	Download(id string, w io.Writer, f func(pg schema.Progress))
}

// func (r Remote) GetTrack(id string, w io.Writer, f func(pg schema.Progress)) {
// 	go func() {
// 		stream, err := r.client.GetTrack(context.Background(), &rpc.String{Value: id})
// 		if err != nil {
// 			log.Fatal("could not get stream for track", err)
// 		}
// 		for {
// 			p, err := stream.Recv()
// 			if err == io.EOF {
// 				time.Sleep(time.Second)
// 			} else if err != nil {
// 				log.Println(err)
// 			} else {
// 				_, err := w.Write(p.Payload)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 				f(rpc.ProgressFromPB(p))
// 			}
// 		}
// 	}()
// }
