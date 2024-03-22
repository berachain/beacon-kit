package e2e_test

func (s *BeaconKitE2ESuite) TestRandao() {
	err := s.WaitForFinalizedBlockNumber(2)
	s.Require().NoError(err)

	client := s.ConsensusClients()["cl-validator-beaconkit-0"]
	s.Require().NotNil(client)

	_, err = client.GetRandao(s.Ctx())
	s.Require().NoError(err)
}
