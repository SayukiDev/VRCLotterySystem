package provider

import (
	"sync/atomic"

	"github.com/SayukiDev/VRCLotterySystem/config"
	"github.com/SayukiDev/VRCLotterySystem/internal/data"
)

type Provider struct {
	IdIndex atomic.Int32
	C       *config.Config
	Data    *data.Data
}

func NewProvider(c *config.Config) *Provider {
	p := &Provider{
		C:       c,
		Data:    data.NewData(),
		IdIndex: atomic.Int32{},
	}
	p.IdIndex.Store(-1)
	return p
}

func (p *Provider) Init() error {
	err := p.Data.Load(p.C.DataPath)
	if err != nil {
		return err
	}
	for i, v := range p.C.Form {
		if v.IsId {
			p.IdIndex.Store(int32(i))
			break
		}
	}
	return nil
}
