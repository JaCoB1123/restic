package checker

import (
	"fmt"

	"github.com/restic/restic/backend"
	"github.com/restic/restic/repository"
)

// Repacker extracts still used blobs from packs with unused blobs and creates
// new packs.
type Repacker struct {
	UnusedBlobs []backend.ID
}

// Repack runs the process of finding still used blobs in packs with unused
// blobs, extracts them and creates new packs with just the still-in-use blobs.
func (r *Repacker) Repack() error {

	return nil
}

// FindPacksforBlobs returns the set of packs that contain the blobs.
func FindPacksforBlobs(repo *repository.Repository, blobs []backend.ID) (backend.IDSet, error) {
	packs := backend.NewIDSet()
	idx := repo.Index()
	for _, id := range blobs {
		packID, _, _, _, err := idx.Lookup(id)
		if err != nil {
			return nil, err
		}

		if packID == nil {
			return nil, fmt.Errorf("blob id %v not found in any index", id)
		}

		packs.Insert(*packID)
	}

	return packs, nil
}
