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
	id := p.GetId(input)
	var err error
	p.Data.Lock(func(d *data.Content) {
		if time.Now().After(d.Date) {
			err = errors.New("date expired")
			return
		}
		if _, ok := d.BlackList[id]; ok {
			return
		}
		if _, ok := d.Forms[id]; ok {
			err = errors.New("input already exists")
			return
		}
		d.Forms[id] = input
		err = d.Forms.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) ClearInputs() error {
	var err error
	p.Data.Lock(func(d *data.Content) {
		d.Forms = make(data.Inputs)
		err = d.Forms.Save(p.C.DataPath)
	})
	return err
}

func (p *Provider) GetInputList() []string {
	var resp []string
	p.Data.RLock(func(d *data.Content) {
		resp = make([]string, 0, len(d.Forms))
		for k := range d.Forms {
			resp = append(resp, k)
		}
	})
	return resp
}

func (p *Provider) GetInput(id string) (*data.Input, error) {
	var input data.Input
	var ok bool
	p.Data.RLock(func(d *data.Content) {
		input, ok = d.Forms[id]
	})
	if !ok {
		return nil, errors.New("input not found")
	}
	return &input, nil
}
