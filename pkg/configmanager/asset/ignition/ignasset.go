package ignition

type IgnAsset struct {
}

// TODO: Init inits the ign asset.
func (ia *IgnAsset) Initial() ([]byte, error) {
	return nil, nil
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
