package BLC

type Chainparams struct {
	TargetBits int
	Subsidy    int
	Fee        int
	Miner	string
}

func (params *Chainparams) init() {
	params.TargetBits = 10
	params.Subsidy = 10
	params.Fee = 1
	params.Miner = ""
}

func (params *Chainparams) setCoinbase() {
	wallets, _ := newWallets()
	for address, _ := range wallets.WalletsMap {
		params.Miner = address
	}
}