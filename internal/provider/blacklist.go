package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) GetBlacklist() []string {
	var resp []string
	p.Data.Read(func(d *data.Content) {
		resp = make([]string, 0, len(d.BlackList))
		for k := range d.BlackList {
			resp = append(resp, k)
		}
	})
	return resp
}

func (p *Provider) AddToBlackList(id string) (err error) {
	p.Data.Write(func(d *data.Content) {
		d.BlackList[id] = struct{}{}
	})
	err = p.Data.Save(p.C.DataPath)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) RemoveFromBlackList(id string) (err error) {
	p.Data.Write(func(d *data.Content) {
		delete(d.BlackList, id)
	})
	err = p.Data.Save(p.C.DataPath)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) IsInBlackList(id string) bool {
	var ok bool
	p.Data.Read(func(d *data.Content) {
		_, ok = d.BlackList[id]
	})
	return ok
}
