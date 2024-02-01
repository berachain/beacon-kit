package service

// WithName sets the name of the BaseService.
func (s BaseService) WithName(name string) BaseService {
	s.logger = s.logger.With("service", name)
	return s
}
