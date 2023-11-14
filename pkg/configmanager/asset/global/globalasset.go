package global

type GlobalAsset struct {
}

// TODO: Init inits the global asset.
func (ga *GlobalAsset) Initial() ([]byte, error) {
	return nil, nil
}

// TODO: Delete deletes the global asset.
func (ga *GlobalAsset) Delete() error {
	return nil
}

// TODO: Persist persists the global asset.
func (ga *GlobalAsset) Persist() error {
	// TODO: Serialize the global asset to json or yaml.
	return nil
}
