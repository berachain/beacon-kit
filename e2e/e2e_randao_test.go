package e2e_test

func (s *BeaconKitE2ESuite) TestRandao() {
	err := s.WaitForFinalizedBlockNumber(6)
	s.Require().NoError(err)

	client := s.ConsensusClients()["cl-validator-beaconkit-0"]
	s.Require().NotNil(client)

	randaoResponse, err := client.GetRandaoMix()
	s.Require().NoError(err)

	err = s.WaitForFinalizedBlockNumber(7)
	s.Require().NoError(err)

	randaoResponse2, err := client.GetRandaoMix()
	s.Require().NoError(err)

	s.Require().NotEqual(randaoResponse.Data.Randao, randaoResponse2.Data.Randao)
	s.Require().NotEqual(randaoResponse.Data.Randao, "0x0000000000000000000000000000000000000000000000000000000000000000")
	s.Require().NotEqual(randaoResponse2.Data.Randao, "0x0000000000000000000000000000000000000000000000000000000000000000")
}
