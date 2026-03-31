package mock

import "github.com/TerraDharitri/drt-go-chain-storage/types"

// DataMergerStub -
type DataMergerStub struct {
	MergeDBsCalled func(dest types.Persister, sources ...types.Persister) error
}

// MergeDBs -
func (stub *DataMergerStub) MergeDBs(dest types.Persister, sources ...types.Persister) error {
	if stub.MergeDBsCalled != nil {
		return stub.MergeDBsCalled(dest, sources...)
	}

	return nil
}

// IsInterfaceNil -
func (stub *DataMergerStub) IsInterfaceNil() bool {
	return stub == nil
}
