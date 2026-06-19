package provider

import "github.com/SayukiDev/VRCLotterySystem/internal/data"

func (p *Provider) GetStaffList() []string {
	var resp []string
	p.Data.RLock(func(d *data.Content) {
		resp = make([]string, 0, len(d.StaffList))
		for k := range d.StaffList {
			resp = append(resp, k)
		}
	})
	return resp
}

func (p *Provider) AddStaff(id string) error {
	var err error
	p.Data.Lock(func(d *data.Content) {
		d.StaffList[id] = struct{}{}
		err = d.StaffList.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) RemoveStaff(id string) error {
	var err error
	p.Data.Lock(func(d *data.Content) {
		delete(d.StaffList, id)
		err = d.StaffList.Save(p.C.DataPath)
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) IsStaff(id string) bool {
	var ok bool
	p.Data.RLock(func(d *data.Content) {
		_, ok = d.StaffList[id]
	})
	return ok
}
