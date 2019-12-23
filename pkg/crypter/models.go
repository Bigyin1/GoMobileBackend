package crypter

// Mapping structure represents relation between user files and saved encrypted file URIs

type Mapping []StoredFile

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

func (m *Mapping) Add(originName, encryptedURL, fid string) {
	*m = append(*m, StoredFile{fid: fid, URL: encryptedURL, Name: originName})
}

func (m *Mapping) AddError(originName, err, fid string) {
	*m = append(*m, StoredFile{fid: fid, Error: err, Name: originName, IsError: true})
}

func (m *Mapping) MergeWith(second Mapping) {
	*m = append(*m, second...)
}

func NewFilesMapping() Mapping {
	return make([]StoredFile, 0)
}
