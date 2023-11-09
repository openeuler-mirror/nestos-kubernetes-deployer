package infra

type TFAsset struct {
}

// TODO: Init inits the tf asset.
func (tf *TFAsset) Init() ([]byte, error) {
	err := tf.validateAsset()
	if err != nil {
		return nil, err
	}
	tf.setAssetDefault()

	return nil, nil
}

func (tf *TFAsset) validateAsset() error {
	return nil
}

func (tf *TFAsset) setAssetDefault() {
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
