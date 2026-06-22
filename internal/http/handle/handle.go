package handle

import (
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"
)

type Handle struct {
	p *provider.Provider
}

func NewHandle(p *provider.Provider) *Handle {
	return &Handle{p: p}
}
