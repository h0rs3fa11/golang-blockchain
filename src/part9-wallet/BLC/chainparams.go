package BLC

type Chainparams struct {
	TargetBits int
	Subsidy    int
	Fee        int
	Miner string
}

func (params *Chainparams) init() {
	params.TargetBits = 10
	params.Subsidy = 10
	params.Fee = 1
<<<<<<< HEAD
	Miner = nil
=======
	params.Miner = ""
>>>>>>> 9d204a21856777a3477c2aa964f463d33e45bc5c
}
