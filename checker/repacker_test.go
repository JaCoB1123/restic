package checker_test

import (
	"fmt"
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
			ParseID("19a731a515618ec8b75fc0ff3b887d8feb83aef1001c9899f6702761142ed068"): struct{}{},
			ParseID("657f7fb64f6a854fff6fe9279998ee09034901eded4e6db9bcee0e59745bbce6"): struct{}{},
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

var repackBlobIDs = backend.IDs{
	ParseID("f41c2089a9d58a4b0bf39369fa37588e6578c928aea8e90a4490a6315b9905c1"),
	ParseID("04fdf6119abd8da279e5c25b0492704d1676043dc2cba1d4f8d40a260d61da65"),
	ParseID("db5ac30c70aaba7fef03db6be91e8d9438e1a417f759f417237efa3482e1f22b"),
	ParseID("356493f0b00a614d36c698591bbb2b1d801932d85328c1f508019550034549fc"),
	ParseID("08d0444e9987fa6e35ce4232b2b71473e1a8f66b2f9664cc44dc57aad3c5a63a"),
	ParseID("b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c"),
	ParseID("5249af22d3b2acd6da8048ac37b2a87fa346fabde55ed23bb866f7618843c9fe"),
	ParseID("b8b5e9c841e2a9d8ad7128ca65de0042a0aeb86abeab6c6400398b3beacb69cb"),
	ParseID("51aa04744b518c6a85b4e7643cfa99d58789c2a6ca2a3fda831fa3032f28535c"),
	ParseID("988a272ab9768182abfd1fe7d7a7b68967825f0b861d3b36156795832c772235"),
	ParseID("aa79d596dbd4c863e5400deaca869830888fe1ce9f51b4a983f532c77f16a596"),
	ParseID("016c84dc8b81eb996c7eb6f19e4302be16177bdbe00dde2352fa1bcdb06c6582"),
	ParseID("454515bca5f4f60349a527bd814cc2681bc3625716460cc6310771c966d8a3bf"),
	ParseID("2a6f01e5e92d8343c4c6b78b51c5a4dc9c39d42c04e26088c7614b13d8d0559d"),
	ParseID("18b51b327df9391732ba7aaf841a4885f350d8a557b2da8352c9acf8898e3f10"),
	ParseID("c01952de4d91da1b1b80bc6e06eaa4ec21523f4853b69dc8231708b9b7ec62d8"),
	ParseID("58c748bbe2929fdf30c73262bd8313fe828f8925b05d1d4a87fe109082acb849"),
	ParseID("b8a6bcdddef5c0f542b4648b2ef79bc0ed4377d4109755d2fb78aff11e042663"),
	ParseID("5714f7274a8aa69b1692916739dc3835d09aac5395946b8ec4f58e563947199a"),
	ParseID("b2396c92781307111accf2ebb1cd62b58134b744d90cb6f153ca456a98dc3e76"),
	ParseID("534f211b4fc2cf5b362a24e8eba22db5372a75b7e974603ff9263f5a471760f4"),
	ParseID("bec3a53d7dc737f9a9bee68b107ec9e8ad722019f649b34d474b9982c3a3fec7"),
}

func TestRepackBlobs(t *testing.T) {
	WithTestEnvironment(t, checkerTestData, func(repodir string) {
		repo := OpenLocalRepo(t, repodir)
		OK(t, repo.LoadIndex())

		packIDs, err := checker.RepackBlobs(repo, repackBlobIDs)
		OK(t, err)

		fmt.Printf("pack IDs: %v\n", packIDs)

		chkr := checker.New(repo)
		OK(t, chkr.LoadIndex())
		// OKs(t, checkPacks(chkr))
		OKs(t, checkStruct(chkr))
	})
}
