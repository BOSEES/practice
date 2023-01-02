package containers

// ImageConfig contains all images and their respective tags
// needed for running e2e tests.
type ImageConfig struct {
	InitRepository string
	InitTag        string

	practiceRepository string
	practiceTag        string

	RelayerRepository string
	RelayerTag        string
}

//nolint:deadcode
const (
	// Current Git branch practice repo/version. It is meant to be built locally.
	// It is used when skipping upgrade by setting practice_E2E_SKIP_UPGRADE to true).
	// This image should be pre-built with `make docker-build-debug` either in CI or locally.
	CurrentBranchRepository = "practice"
	CurrentBranchTag        = "debug"
	// Pre-upgrade practice repo/tag to pull.
	// It should be uploaded to Docker Hub. practice_E2E_SKIP_UPGRADE should be unset
	// for this functionality to be used.
	previousVersionRepository = "ghcr.io/BOSEES/practice"
	previousVersionTag        = "11.0.0-e2e"
	// Pre-upgrade repo/tag for practice initialization (this should be one version below upgradeVersion)
	previousVersionInitRepository = "ghcr.io/BOSEES/practice"
	previousVersionInitTag        = "11.0.0-e2e-init-chain"
	// Hermes repo/version for relayer
	relayerRepository = "ghcr.io/cosmoscontracts/hermes"
	relayerTag        = "0.13.0"
)

// Returns ImageConfig needed for running e2e test.
// If isUpgrade is true, returns images for running the upgrade
// If isFork is true, utilizes provided fork height to initiate fork logic
func NewImageConfig(isUpgrade, isFork bool) ImageConfig {
	config := ImageConfig{
		RelayerRepository: relayerRepository,
		RelayerTag:        relayerTag,
	}

	if !isUpgrade {
		// If upgrade is not tested, we do not need InitRepository and InitTag
		// because we directly call the initialization logic without
		// the need for Docker.
		config.practiceRepository = CurrentBranchRepository
		config.practiceTag = CurrentBranchTag
		return config
	}

	// If upgrade is tested, we need to utilize InitRepository and InitTag
	// to initialize older state with Docker
	config.InitRepository = previousVersionInitRepository
	config.InitTag = previousVersionInitTag

	if isFork {
		// Forks are state compatible with earlier versions before fork height.
		// Normally, validators switch the binaries pre-fork height
		// Then, once the fork height is reached, the state breaking-logic
		// is run.
		config.practiceRepository = CurrentBranchRepository
		config.practiceTag = CurrentBranchTag
	} else {
		// Upgrades are run at the time when upgrade height is reached
		// and are submitted via a governance proposal. Thefore, we
		// must start running the previous practice version. Then, the node
		// should auto-upgrade, at which point we can restart the updated
		// practice validator container.
		config.practiceRepository = previousVersionRepository
		config.practiceTag = previousVersionTag
	}

	return config
}
