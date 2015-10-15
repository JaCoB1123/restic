package checker

import (
	"errors"
	"fmt"

	"github.com/restic/restic/backend"
	"github.com/restic/restic/crypto"
	"github.com/restic/restic/debug"
	"github.com/restic/restic/repository"
)

// Repacker extracts still used blobs from packs with unused blobs and creates
// new packs.
type Repacker struct {
	unusedBlobs []backend.ID
	repo        *repository.Repository
}

// NewRepacker returns a new repacker that (when Repack() in run) cleans up the
// repository and creates new packs and indexs so that all blobs in unusedBlobs
// aren't used any more.
func NewRepacker(repo *repository.Repository, unusedBlobs []backend.ID) *Repacker {
	return &Repacker{
		repo:        repo,
		unusedBlobs: unusedBlobs,
	}
}

// Repack runs the process of finding still used blobs in packs with unused
// blobs, extracts them and creates new packs with just the still-in-use blobs.
func (r *Repacker) Repack() error {
	debug.Log("Repacker.Repack", "searching packs for %v", r.unusedBlobs)
	packs, err := FindPacksforBlobs(r.repo, r.unusedBlobs)
	if err != nil {
		return err
	}

	debug.Log("Repacker.Repack", "found packs: %v", packs)

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

// RepackBlobs reads and packs all the blobs in blobIDs into new pack files.
func RepackBlobs(repo *repository.Repository, blobIDs backend.IDs) (packIDs map[backend.ID]struct{}, err error) {
	packIDs = make(map[backend.ID]struct{})

	for _, id := range blobIDs {
		_, tpe, _, length, err := repo.Index().Lookup(id)
		if err != nil {
			return nil, err
		}

		length -= crypto.Extension

		debug.Log("RepackBlobs", "repacking blob %v, len %v", id.Str(), length)

		buf := make([]byte, 0, length)
		buf, err = repo.LoadBlob(tpe, id, buf)
		if err != nil {
			return nil, err
		}

		if uint(len(buf)) != length {
			debug.Log("RepackBlobs", "repack blob %v: len(buf) isn't equal to length: %v = %v", id.Str(), len(buf), length)
			return nil, errors.New("LoadBlob returned wrong data, len() doesn't match")
		}

		_, err = repo.SaveAndEncrypt(tpe, buf, &id)
		if err != nil {
			return nil, err
		}
	}

	err = repo.Flush()
	if err != nil {
		return nil, err
	}

	err = repo.SaveFullIndex()
	if err != nil {
		return nil, err
	}

	for _, id := range blobIDs {
		packID, _, _, _, err := repo.Index().Lookup(id)
		if err != nil {
			return nil, err
		}

		if packID == nil {
			return nil, errors.New("packID is nil after Flush()")
		}

		debug.Log("RepackBlobs", "blob %v has been saved to pack %v", id.Str(), packID.Str())
		packIDs[*packID] = struct{}{}
	}

	return packIDs, nil
}
