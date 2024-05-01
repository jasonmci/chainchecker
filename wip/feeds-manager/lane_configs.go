package main

import (

)

type TestConfig struct {
    LaneConfigs map[string]LaneConfig `json:"lane_configs"`
}

type LaneConfig struct {
    IsNativeFeeToken   bool                  `json:"is_native_fee_token"`
    FeeToken           string                `json:"fee_token"`
    BridgeTokens       []string              `json:"bridge_tokens"`
    BridgeTokensPools  []string              `json:"bridge_tokens_pools"`
    Arm                string                `json:"arm"`
    Router             string                `json:"router"`
    PriceRegistry      string                `json:"price_registry"`
    WrappedNative      string                `json:"wrapped_native"`
    SrcContracts       map[string]SrcContract `json:"src_contracts"`
    DestContracts      map[string]DestContract `json:"dest_contracts"`
}

type SrcContract struct {
    OnRamp     string `json:"on_ramp"`
    DeployedAt int64  `json:"deployed_at"`
}

type DestContract struct {
    OffRamp       string `json:"off_ramp"`
    CommitStore   string `json:"commit_store"`
    ReceiverDapp  string `json:"receiver_dapp"`
}

