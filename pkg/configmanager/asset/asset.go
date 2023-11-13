package asset

type Asset interface {
	Init() ([]byte, error)
	Delete() error
	Persist() error
}

func InitAsset(asset Asset) ([]byte, error) {
	assetData, err := asset.Init()
	if err != nil {
		return nil, err
	}

	return assetData, err
}

func DeleteAsset(asset Asset) error {
	err := asset.Delete()
	if err != nil {
		return err
	}

	return nil
}

func PersistAsset(asset Asset) error {
	err := asset.Persist()
	if err != nil {
		return err
	}

	return nil
}
