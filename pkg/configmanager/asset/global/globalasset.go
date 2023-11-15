package global

type GlobalAsset struct {
}

var (
	GlobalConfig GlobalAsset
	IsInitial    = false
)

// TODO: Init inits the global asset.
func (ga *GlobalAsset) Initial() error {
	// TODO: 将初始化的结果传给GlobalConfig
	IsInitial = true
	return nil
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
