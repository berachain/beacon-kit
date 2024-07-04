package collections

func query(store Store, storeKey, key []byte) ([]byte, error) {
	version, err := store.GetLatestVersion()
	if err != nil {
		return nil, err
	}
	resp, err := store.Query(storeKey, version, key, false)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}
