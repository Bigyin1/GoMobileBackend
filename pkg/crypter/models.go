package crypter

// Mapping structure represents relation between user files and saved encrypted file URIs
type Mapping struct {
	Mapping map[string]string `json:"mapping"`
}

func (m *Mapping) Add(originName, encryptedURI string) {
	m.Mapping[originName] = encryptedURI
}

func newFilesMapping() Mapping {
	return Mapping{Mapping: make(map[string]string)}
}
