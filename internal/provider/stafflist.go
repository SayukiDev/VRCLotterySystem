package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) GetStaffList() []string {
	d := p.Data.Get()
	resp := make([]string, 0, len(d.StaffList))
	for k := range d.StaffList {
		resp = append(resp, k)
	}
	return resp
}

func (p *Provider) AddStaff(id string) error {
	var err error
	p.Data.Set(func(data *data.Data) {
		data.StaffList[id] = struct{}{}
		err = data.StaffList.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) RemoveStaff(id string) error {
	var err error
	p.Data.Set(func(data *data.Data) {
		delete(data.StaffList, id)
		err = data.StaffList.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) IsStaff(id string) bool {
	d := p.Data.Get()
	_, ok := d.StaffList[id]
	return ok
}
