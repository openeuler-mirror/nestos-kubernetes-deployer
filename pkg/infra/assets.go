package infra

// path : contents
type Assets map[string][]byte

func (a Assets) ToDir(dirname string) error {
	return nil
}

func (a *Assets) Merge(b Assets) *Assets {
	return a
}

type AssetsGenerator interface {
	GenerateAssets() Assets
}