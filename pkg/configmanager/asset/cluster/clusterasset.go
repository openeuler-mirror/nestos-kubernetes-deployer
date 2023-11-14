package cluster

type ClusterAsset struct {
}

// TODO: Init inits the cluster asset.
func (ca *ClusterAsset) Initial() ([]byte, error) {
	return nil, nil
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
