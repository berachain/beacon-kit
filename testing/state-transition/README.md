# Testing state transitions

This package, `statetransition`, contains code which is helpful to test state transitions.

Specifically, it contains code to instantiate a test state processor with a mocked bls signer and its own state.

This code can be reused across tests in multiple packages but **should not be used in production code**.