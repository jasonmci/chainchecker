package main

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

// AddChain adds a new chain with default configuration to the TestConfig
func (tc *TestConfig) AddChain(chainName string) {
    if tc.LaneConfigs == nil {
        tc.LaneConfigs = make(map[string]LaneConfig)
    }
    tc.LaneConfigs[chainName] = LaneConfig{
        IsNativeFeeToken:  true,
        FeeToken:          "",
        BridgeTokens:      []string{},
        BridgeTokensPools: []string{},
        Arm:               "",
        Router:            "",
        PriceRegistry:     "",
        WrappedNative:     "",
        SrcContracts:      make(map[string]SrcContract),
        DestContracts:     make(map[string]DestContract),
    }
}

// AddSrcContract adds a new source contract for a specified chain
func (tc *TestConfig) AddSrcContract(chainName, srcChainName, onRamp string, deployedAt int64) {
    
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.SrcContracts[srcChainName] = SrcContract{
            OnRamp:     onRamp,
            DeployedAt: deployedAt,
        }
        tc.LaneConfigs[chainName] = lane
    }
}

// AddDestContract adds a new destination contract for a specified chain
func (tc *TestConfig) AddDestContract(chainName, destChainName, offRamp, commitStore, receiverDapp string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.DestContracts[destChainName] = DestContract{
            OffRamp:      offRamp,
            CommitStore:  commitStore,
            ReceiverDapp: receiverDapp,
        }
        tc.LaneConfigs[chainName] = lane
    }
}

// UpdateFeeToken updates the fee token for a specified chain
func (tc *TestConfig) UpdateFeeToken(chainName, feeToken string, isNative bool) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.FeeToken = feeToken
        lane.IsNativeFeeToken = isNative
        tc.LaneConfigs[chainName] = lane
    }
}

// AddBridgeToken adds a new bridge token and its corresponding pool
func (tc *TestConfig) AddBridgeToken(chainName, token, pool string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.BridgeTokens = append(lane.BridgeTokens, token)
        lane.BridgeTokensPools = append(lane.BridgeTokensPools, pool)
        tc.LaneConfigs[chainName] = lane
    }
}

func (tc *TestConfig) UpdateArm(chainName, arm string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.Arm = arm
        tc.LaneConfigs[chainName] = lane
    }
}
// UpdateRouter updates the router for a specified chain
func (tc *TestConfig) UpdateRouter(chainName, router string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.Router = router
        tc.LaneConfigs[chainName] = lane
    }
}

// UpdatePriceRegistry updates the price registry for a specified chain
func (tc *TestConfig) UpdatePriceRegistry(chainName, priceRegistry string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.PriceRegistry = priceRegistry
        tc.LaneConfigs[chainName] = lane
    }
}

// UpdateWrappedNative updates the wrapped native token address for a specified chain
func (tc *TestConfig) UpdateWrappedNative(chainName, wrappedNative string) {
    if lane, ok := tc.LaneConfigs[chainName]; ok {
        lane.WrappedNative = wrappedNative
        tc.LaneConfigs[chainName] = lane
    }
}

func populateLaneConfigFromResponse( genConfig GenConfig, laneResponse *CCIPLaneResponse, chainAResponse *CCIPChainResponse, chainBResponse *CCIPChainResponse,
    config *TestConfig) {

    // var chainAName string
    // var chainBName string

    config.AddChain(chainAResponse.Data.CCIP.Chain.Network.Name)

    // if chainAname is GnosisChain Mainnet rename it to Gnosis Mainnet
    if chainAResponse.Data.CCIP.Chain.Network.Name == "Gnosis Chain Mainnet" {
        chainAResponse.Data.CCIP.Chain.Network.Name = "Gnosis Mainnet"
    }

    // if chainAname is GnosisChain Mainnet rename it to Gnosis Mainnet
    if chainBResponse.Data.CCIP.Chain.Network.Name == "Gnosis Chain Mainnet" {
        chainBResponse.Data.CCIP.Chain.Network.Name = "Gnosis Mainnet"
    }

    config.AddSrcContract(chainAResponse.Data.CCIP.Chain.Network.Name, chainBResponse.Data.CCIP.Chain.Network.Name, "0x1", 111111)
    config.AddDestContract(chainAResponse.Data.CCIP.Chain.Network.Name, chainBResponse.Data.CCIP.Chain.Network.Name, "0x2", "0x3", "0x4")
    config.UpdateWrappedNative(chainAResponse.Data.CCIP.Chain.Network.Name, genConfig.NativeFeeTokenMap[chainAResponse.Data.CCIP.Chain.Network.Name])

    config.AddChain(chainBResponse.Data.CCIP.Chain.Network.Name)
    config.AddSrcContract(chainBResponse.Data.CCIP.Chain.Network.Name, chainAResponse.Data.CCIP.Chain.Network.Name, "0x5", 111111)
    config.AddDestContract(chainBResponse.Data.CCIP.Chain.Network.Name, chainAResponse.Data.CCIP.Chain.Network.Name, "0x6", "0x7", "0x8")
    config.UpdateWrappedNative(chainBResponse.Data.CCIP.Chain.Network.Name, genConfig.NativeFeeTokenMap[chainBResponse.Data.CCIP.Chain.Network.Name])

 }
