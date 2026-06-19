package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	BlackListPath = "/blacklist.json"
	StaffListPath = "/stafflist.json"
	InputsPath    = "/inputs.json"
	ResultsPath   = "/results.json"
)

type Data struct {
	lock    sync.RWMutex
	content Content
}
type Content struct {
	Id         string    `json:"id"`
	Date       time.Time `json:"date"`
	Showed     bool      `json:"showed"`
	Max        int       `json:"number"`
	ExtContent `json:"-"`
}

type ExtContent struct {
	BlackList BlackList
	StaffList StaffList
	Forms     Inputs
	Results   Results
}

func save(path string, v any) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(v)
	if err != nil {
		return err
	}
	return nil
}

func load(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0755)
			return save(path, v)
		}
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

type BlackList map[string]struct{}

func (b *BlackList) Save(path string) error {
	return save(filepath.Join(path, BlackListPath), b)
}

func (b *BlackList) Load(path string) error {
	return load(filepath.Join(path, BlackListPath), b)
}

type StaffList map[string]struct{}

func (s *StaffList) Save(path string) error {
	return save(filepath.Join(path, StaffListPath), s)
}

func (s *StaffList) Load(path string) error {
	return load(filepath.Join(path, StaffListPath), s)
}

type Inputs map[string]Input
type Input struct {
	Content  map[int]string `json:"content"`
	Selected map[int][]int  `json:"selected"`
}

func (i *Inputs) Save(path string) error {
	return save(filepath.Join(path, InputsPath), i)
}

func (i *Inputs) Load(path string) error {
	return load(filepath.Join(path, InputsPath), i)
}

type Results map[string]struct{}

func (r *Results) Save(path string) error {
	return save(filepath.Join(path, ResultsPath), r)
}

func (r *Results) Load(path string) error {
	return load(filepath.Join(path, ResultsPath), r)
}

func NewData() *Data {
	return &Data{
		content: Content{
			ExtContent: ExtContent{
				BlackList: make(BlackList),
				StaffList: make(StaffList),
				Forms:     make(Inputs),
				Results:   make(Results),
			},
		},
	}
}

func (d *Data) Save(path string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := d.content.BlackList.Save(path)
	if err != nil {
		return err
	}
	err = d.content.Forms.Save(path)
	if err != nil {
		return err
	}
	err = d.content.Results.Save(path)
	if err != nil {
		return err
	}
	err = d.content.StaffList.Save(path)
	if err != nil {
		return err
	}
	err = save(filepath.Join(path, "data.json"), d.content)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) Load(path string) error {
	err := load(filepath.Join(path, "data.json"), &d.content)
	if err != nil {
		return err
	}
	err = d.content.BlackList.Load(path)
	if err != nil {
		return err
	}
	err = d.content.StaffList.Load(path)
	if err != nil {
		return err
	}
	err = d.content.Forms.Load(path)
	if err != nil {
		return err
	}
	err = d.content.Results.Load(path)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) Read(handle func(data *Content)) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	handle(&d.content)
}

func (d *Data) Write(handle func(data *Content)) {
	d.lock.Lock()
	defer d.lock.Unlock()
	handle(&d.content)
}
