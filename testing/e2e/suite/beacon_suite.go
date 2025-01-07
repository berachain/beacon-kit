package suite

// BeaconKitE2ESuite is a suite of tests simulating a fully functional beacon-kit network.
type BeaconKitE2ESuite struct {
	KurtosisE2ESuite
}

func NewBeaconKitE2ESuite() *BeaconKitE2ESuite {
	return &BeaconKitE2ESuite{
		KurtosisE2ESuite: KurtosisE2ESuite{
			networks:  make(map[string]*NetworkInstance),
			testSpecs: make(map[string]string),
		},
	}
}

func (s *BeaconKitE2ESuite) SetupSuite() {
	s.KurtosisE2ESuite.SetupSuite()
}
