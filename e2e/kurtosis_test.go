package e2e_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

func TestXxx(t *testing.T) {
	ctx := context.Background()
	kCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		t.Fatalf("Error instantiating Kurtosis context: %v", err)
	}

	enclave, err := kCtx.CreateEnclave(ctx, "e2e-test-enclave")
	if err != nil {
		t.Fatalf("Error creating enclave: %v", err)
	}

	defer func() {
		if err := kCtx.StopEnclave(ctx, "e2e-test-enclave"); err != nil {
			t.Fatalf("Error stopping enclave: %v", err)
		}
	}()

	resp, cancel, err := enclave.RunStarlarkPackage(
		ctx,
		"../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(),
	)
	defer cancel()

	if err != nil {
		t.Fatalf("Error running Starlark package: %v", err)
	}

	fmt.Println(resp)
	time.Sleep(30 * time.Second)
}
