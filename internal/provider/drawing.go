package provider

import (
	"time"

	"github.com/SayukiDev/VRCLotterySystem/internal/data"
	"github.com/SayukiDev/VRCLotterySystem/pkg/eth"
)

const ethUrl = "https://ethereum-rpc.publicnode.com"

// SetDrawing Date format: 2006-01-02 15:04
func (p *Provider) SetDrawing(max int, date string) error {
	id, err := eth.NewClient(ethUrl).RandomString(8)
	if err != nil {
		return err
	}
	p.Data.Lock(func(d *data.Content) {
		d.Id = id
		d.Date, err = time.Parse("2006-01-02 15:04", date)
		d.Showed = false
		d.Max = max
		d.Forms = make(data.Inputs)
		d.Results = make(data.Results)
	})
	if err != nil {
		return err
	}
	err = p.Data.Save(p.C.DataPath)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) Drawing() ([]string, error) {
	c := eth.NewClient(ethUrl)
	var formList []string
	var max int
	p.Data.RLock(func(d *data.Content) {
		formList = make([]string, 0, len(d.Forms))
		for k := range d.Forms {
			formList = append(formList, k)
		}
		max = d.Max
	})
	if len(formList) == 0 {
		return []string{}, nil
	}
	indexs, err := c.RandomIndexInRange(0, len(formList)-1, len(formList))
	if err != nil {
		return nil, err
	}
	take := min(max, len(formList))
	ids := make([]string, 0, take)
	for _, i := range indexs[:take] {
		ids = append(ids, formList[i])
	}
	p.Data.Lock(func(d *data.Content) {
		d.Showed = true
	})
	err = p.Data.Save(p.C.DataPath)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
