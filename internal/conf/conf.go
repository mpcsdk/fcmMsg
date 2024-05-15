package conf

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

type Server struct {
	Address string `json:"address" v:"required"`
	WorkId  int    `json:"workId" v:"required|min:1"`
	Name    string `json:"name" v:"required"`
}

type Nrpcfg struct {
	NatsUrl string `json:"natsUrl" v:"required"`
}

type Cache struct {
	Duration          int    `json:"duration" v:"required|min:10000"`
	RetentionDataTime string `json:"retentionDataTime" v:"required"`
}

// /
type Jaeger struct {
	Enable bool   `json:"enable" v:"required"`
	Url    string `json:"url" v:"required"`
}

// //
type Cfg struct {
	Server *Server `json:"server" v:"required"`
	Jaeger *Jaeger `json:"jaeger" v:"required"`
	Nrpc   *Nrpcfg `json:"nrpc" v:"required"`
	Cache  *Cache  `json:"cache" v:"required"`
}

var Config = &Cfg{}

func init() {
	ctx := gctx.GetInitCtx()
	cfg := gcfg.Instance()
	v, err := cfg.Data(ctx)
	if err != nil {
		panic(err)
	}
	val := gvar.New(v)
	err = val.Structs(Config)
	if err != nil {
		panic(err)
	}
	if err := g.Validator().Data(Config).Run(ctx); err != nil {
		panic(err)
	}
}
