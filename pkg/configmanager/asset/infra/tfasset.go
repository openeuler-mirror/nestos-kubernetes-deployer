package infra

type TFAsset struct {
}

// TODO: Init inits the tf asset.
func (tf *TFAsset) Initial() ([]byte, error) {
	return nil, nil
}

// TODO: Delete deletes the tf asset.
func (tf *TFAsset) Delete() error {
	return nil
}

// TODO: Persist persists the tf asset.
func (tf *TFAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	return nil
}
