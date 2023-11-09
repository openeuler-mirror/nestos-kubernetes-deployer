package ignition

type IgnAsset struct {
}

// TODO: Init inits the ign asset.
func (ia *IgnAsset) Init() ([]byte, error) {
	err := ia.validateAsset()
	if err != nil {
		return nil, err
	}
	ia.setAssetDefault()

	return nil, nil
}

func (ia *IgnAsset) validateAsset() error {
	return nil
}

func (ia *IgnAsset) setAssetDefault() {
}

// TODO: Delete deletes the ign asset.
func (ia *IgnAsset) Delete() error {
	return nil
}

// TODO: Persist persists the ign asset.
func (ia *IgnAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	return nil
}
