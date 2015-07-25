package checker_test

import (
	"testing"

	"github.com/restic/restic/backend"
	"github.com/restic/restic/checker"

	. "github.com/restic/restic/test"
)

var findPackTests = []struct {
	blobIDs backend.IDs
	packIDs backend.IDSet
}{
	{
		backend.IDs{
			ParseID("534f211b4fc2cf5b362a24e8eba22db5372a75b7e974603ff9263f5a471760f4"),
			ParseID("51aa04744b518c6a85b4e7643cfa99d58789c2a6ca2a3fda831fa3032f28535c"),
			ParseID("454515bca5f4f60349a527bd814cc2681bc3625716460cc6310771c966d8a3bf"),
			ParseID("c01952de4d91da1b1b80bc6e06eaa4ec21523f4853b69dc8231708b9b7ec62d8"),
		},
		backend.IDSet{
			backend.ID{0x19, 0xa7, 0x31, 0xa5, 0x15, 0x61, 0x8e, 0xc8,
				0xb7, 0x5f, 0xc0, 0xff, 0x3b, 0x88, 0x7d, 0x8f,
				0xeb, 0x83, 0xae, 0xf1, 0x00, 0x1c, 0x98, 0x99,
				0xf6, 0x70, 0x27, 0x61, 0x14, 0x2e, 0xd0, 0x68}: struct{}{},
			backend.ID{0x65, 0x7f, 0x7f, 0xb6, 0x4f, 0x6a, 0x85, 0x4f,
				0xff, 0x6f, 0xe9, 0x27, 0x99, 0x98, 0xee, 0x09,
				0x03, 0x49, 0x01, 0xed, 0xed, 0x4e, 0x6d, 0xb9,
				0xbc, 0xee, 0x0e, 0x59, 0x74, 0x5b, 0xbc, 0xe6}: struct{}{},
		},
	},
}

func TestRepackerFindPacks(t *testing.T) {
	WithTestEnvironment(t, checkerTestData, func(repodir string) {
		repo := OpenLocalRepo(t, repodir)

		OK(t, repo.LoadIndex())

		for _, test := range findPackTests {
			packIDs, err := checker.FindPacksforBlobs(repo, test.blobIDs)
			OK(t, err)
			Equals(t, test.packIDs, packIDs)
		}
	})
}
