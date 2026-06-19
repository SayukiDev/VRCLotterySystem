package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) AddResults(id []string) error {
	var err error
	p.Data.Lock(func(d *data.Content) {
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
	var ok bool
	p.Data.RLock(func(d *data.Content) {
		_, ok = d.Results[id]
	})
	return ok
}

func (p *Provider) ClearResults() error {
	var err error
	p.Data.Lock(func(d *data.Content) {
		d.Results = make(data.Results)
		err = d.Results.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GetResults() []string {
	var resp []string
	p.Data.RLock(func(d *data.Content) {
		resp = make([]string, 0, len(d.Results))
		for k := range d.Results {
			resp = append(resp, k)
		}
	})
	return resp
}

func (p *Provider) DeleteResults(id string) error {
	var err error
	p.Data.Lock(func(d *data.Content) {
		delete(d.Results, id)
		err = d.Results.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}
