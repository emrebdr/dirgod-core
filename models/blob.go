package models

type Blob struct {
	BlobId   string
	Name     string
	Content  []byte
	Checksum string
	Path     string
}

func (b *Blob) Print() {

}

func (b *Blob) GetName() string {
	return b.Name
}

func (b *Blob) GetType() string {
	return "blob"
}

func (b *Blob) GetChecksum() string {
	return b.Checksum
}
