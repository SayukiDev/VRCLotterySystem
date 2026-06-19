package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) AddResults(id []string) error {
	var err error
	p.Data.Set(func(d *data.Data) {
		for _, v := range id {
			d.Results[v] = struct{}{}
		}
		err = d.Results.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) IsInResults(id string) bool {
	_, ok := p.Data.Results[id]
	return ok
}

func (p *Provider) ClearResults() error {
	var err error
	p.Data.Set(func(d *data.Data) {
		d.Results = make(data.Results)
		err = d.Results.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GetResults() []string {
	resp := make([]string, 0, len(p.Data.Results))
	for k := range p.Data.Results {
		resp = append(resp, k)
	}
	return resp
}

func (p *Provider) DeleteResults(id string) error {
	var err error
	p.Data.Set(func(d *data.Data) {
		delete(d.Results, id)
		err = d.Results.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}
