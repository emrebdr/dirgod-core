package models

type Tree struct {
	TreeId   string
	Name     string
	Blobs    []string
	Trees    []string
	Checksum string
	Path     string
}

func (t *Tree) Print() {

}

func (t *Tree) GetName() string {
	return t.Name
}

func (t *Tree) GetType() string {
	return "tree"
}

func (t *Tree) GetChecksum() string {
	return t.Checksum
}
