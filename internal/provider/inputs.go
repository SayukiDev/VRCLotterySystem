package provider

import (
	"errors"
	"time"

	"github.com/SayukiDev/VRCLotterySystem/internal/data"
)

func (p *Provider) GetId(input data.Input) string {
	index := p.IdIndex.Load()
	if index < 0 {
		return ""
	}
	id, ok := input.Content[int(index)]
	if !ok {
		return ""
	}
	return id
}

func (p *Provider) AddInput(input data.Input) error {
	d := p.Data.Get()
	id := p.GetId(input)
	if time.Now().After(d.Date) {
		return errors.New("date expired")
	}
	if _, ok := d.BlackList[id]; ok {
		return nil
	}
	if _, ok := d.Forms[id]; ok {
		return errors.New("input already exists")
	}
	p.Data.Set(func(d *data.Data) {
		d.Forms[id] = input
	})
	return nil
}

func (p *Provider) ClearInputs() error {
	var err error
	p.Data.Set(func(d *data.Data) {
		d.Forms = make(data.Inputs)
		err = d.Forms.Save(p.C.DataPath)
	})
	return err
}

func (p *Provider) GetInputList() []string {
	resp := make([]string, 0, len(p.Data.Get().Forms))
	for k := range p.Data.Get().Forms {
		resp = append(resp, k)
	}
	return resp
}

func (p *Provider) GetInput(id string) (*data.Input, error) {
	input, ok := p.Data.Get().Forms[id]
	if !ok {
		return nil, errors.New("input not found")
	}
	return &input, nil
}
