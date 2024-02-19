package staking

// Option defines a function type that applies a configuration to the Service.
type Option func(*Service) error

// WithValsetChangeProvider returns an Option that sets the ValsetChangeProvider
// for the Service. This is used to inject the dependency that handles
// the application of changes to the validator set.
func WithValsetChangeProvider(vcp ValsetChangeProvider) Option {
	return func(s *Service) error {
		s.vcp = vcp
		return nil
	}
}
