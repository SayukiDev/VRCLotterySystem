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
	Id      string    `json:"id"`
	Date    time.Time `json:"date"`
	Showed  bool      `json:"showed"`
	Max     int       `json:"number"`
	Content `json:"-"`
}
type Content struct {
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
		Content: Content{
			Forms:     make(Inputs),
			BlackList: make(BlackList),
			StaffList: make(StaffList),
			Results:   make(Results),
		},
	}
}

func (d *Data) Save(path string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := d.BlackList.Save(path)
	if err != nil {
		return err
	}
	err = d.Forms.Save(path)
	if err != nil {
		return err
	}
	err = d.Results.Save(path)
	if err != nil {
		return err
	}
	err = save(filepath.Join(path, "data.json"), d.Content)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) Load(path string) error {
	err := d.BlackList.Load(path)
	if err != nil {
		return err
	}
	err = d.Forms.Load(path)
	if err != nil {
		return err
	}
	err = d.Results.Load(path)
	if err != nil {
		return err
	}
	err = save(filepath.Join(path, "data.json"), d.Content)
	if err != nil {
		return err
	}
	return nil
}

// Get will return a pointer to the Data object
// If you need to change the Data object, please using the Set method
// DO NOT modify the returned Data object
func (d *Data) Get() *Data {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d
}

func (d *Data) Set(handle func(data *Data)) {
	d.lock.Lock()
	defer d.lock.Unlock()
	handle(d)
}
