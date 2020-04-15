package rpc

import (
	"fmt"
)

func (r *Rpc) Help(arg string, reply *string) error {
	*reply = fmt.Sprintln(`Usage:\n
		\tgetbalance -address ADDRESS - Get balance of ADDRESS\n /
		\tcreateblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS\n
		\tprintchain - Print all the blocks of the blockchain\n
		\tsendmany -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO\n
		\tlistaddress list all address from wallet\n
		\tcreateaddress create a new address\n
	`)
	return nil
}