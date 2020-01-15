package BLC

type Chainparams struct {
	TargetBits int
	Subsidy    int
	Fee        int
}

func (params *Chainparams) init() {
	params.TargetBits = 10
	params.Subsidy = 10
	params.Fee = 1
}
