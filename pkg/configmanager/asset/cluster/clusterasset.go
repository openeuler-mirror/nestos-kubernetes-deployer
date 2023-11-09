package cluster

type ClusterAsset struct {
}

// TODO: Init inits the cluster asset.
func (ca *ClusterAsset) Init() ([]byte, error) {
	err := ca.validateAsset()
	if err != nil {
		return nil, err
	}
	ca.setAssetDefault()

	return nil, nil
}

func (ca *ClusterAsset) validateAsset() error {
	return nil
}

func (ca *ClusterAsset) setAssetDefault() {
}

// TODO: Delete deletes the cluster asset.
func (ca *ClusterAsset) Delete() error {
	return nil
}

// TODO: Persist persists the cluster asset.
func (ca *ClusterAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	return nil
}
