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
	p.Data.Set(func(d *data.Data) {
		d.Id = id
		d.Date, err = time.Parse("2006-01-02 15:04", date)
		d.Showed = false
		d.Max = max
		d.Forms = make(data.Inputs)
		d.Results = make(data.Results)
		err = d.Save(p.C.DataPath)
	})
	return nil
}

func (p *Provider) Drawing() ([]string, error) {
	d := p.Data.Get()
	c := eth.NewClient(ethUrl)
	formList := make([]string, 0, len(d.Forms))
	for k := range d.Forms {
		formList = append(formList, k)
	}
	ids := make([]string, 0, d.Max)
	indexs, err := c.RandomsInRange(0, int64(len(d.Forms)-1), min(d.Max, len(d.Forms)))
	if err != nil {
		return nil, err
	}
	for _, i := range indexs {
		id := formList[i]
		ids = append(ids, id)
	}
	return ids, nil
}
