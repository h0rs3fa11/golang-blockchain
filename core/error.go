package core

func New(infor string) error {
	return &blockchainError{infor}
}

type blockchainError struct {
	infor string
}

func (be *blockchainError) Error() string {
	return be.infor
}