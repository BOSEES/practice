package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultpracticeInstanceCost is initially set the same as in wasmd
	DefaultpracticeInstanceCost uint64 = 60_000
	// DefaultpracticeCompileCost set to a large number for testing
	DefaultpracticeCompileCost uint64 = 3
)

// practiceGasRegisterConfig is defaults plus a custom compile amount
func practiceGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultpracticeInstanceCost
	gasConfig.CompileCost = DefaultpracticeCompileCost

	return gasConfig
}

func NewpracticeWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(practiceGasRegisterConfig())
}
