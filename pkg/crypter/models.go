package crypter

import "sync"

// Mapping structure represents relation between user files and saved encrypted file URIs

type StoredFile struct {
	fid     string
	URL     string `json:"url,omitempty"`
	Name    string `json:"name"`
	Error   string `json:"error,omitempty"`
	IsError bool   `json:"-"`
}

func (f *StoredFile) GetFid() string {
	return f.fid
}

func (f *StoredFile) GetUrlOrErr() string {
	if f.URL != "" {
		return f.URL
	}
	return f.Error
}

type Mapping struct {
	mapping []StoredFile
	m       *sync.Mutex
}

func (m *Mapping) Add(originName, encryptedURL, fid string) {
	m.m.Lock()
	defer m.m.Unlock()
	m.mapping = append(m.mapping, StoredFile{fid: fid, URL: encryptedURL, Name: originName})
}

func (m *Mapping) AddError(originName, err, fid string) {
	m.m.Lock()
	defer m.m.Unlock()
	m.mapping = append(m.mapping, StoredFile{fid: fid, Error: err, Name: originName, IsError: true})
}

func (m *Mapping) GetMapping() []StoredFile {
	return m.mapping
}

func NewFilesMapping() Mapping {
	return Mapping{mapping: make([]StoredFile, 0), m: &sync.Mutex{}}
}
