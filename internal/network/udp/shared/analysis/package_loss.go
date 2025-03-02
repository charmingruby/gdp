package analysis

import (
	"github.com/charmingruby/gdp/internal/storage"
)

type PackageLossData struct {
	PackagesReceived int `json:"packages_received"`
	PackagesLost     int `json:"packages_lost"`
}

func SavePackageLossData(path string, data PackageLossData) error {
	return storage.SaveRecord(path, storage.Record{
		Data: data,
	})
}
