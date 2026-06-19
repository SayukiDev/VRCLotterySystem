package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) GetBlacklist() []string {
	d := p.Data.Get()
	resp := make([]string, 0, len(d.BlackList))
	for k := range d.BlackList {
		resp = append(resp, k)
	}
	return resp
}

func (p *Provider) AddToBlackList(id string) (err error) {
	p.Data.Set(func(d *data.Data) {
		d.BlackList[id] = struct{}{}
		err = d.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) RemoveFromBlackList(id string) (err error) {
	p.Data.Set(func(d *data.Data) {
		delete(d.BlackList, id)
		err = d.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) IsInBlackList(id string) bool {
	d := p.Data.Get()
	_, ok := d.BlackList[id]
	return ok
}
