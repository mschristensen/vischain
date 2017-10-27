# Brocoin

Forget bro points. Real bros get brocoins.

Cmon bro.

## Bro?

Brocoin is a simple cryptocurrency which I have developed for the sole purpose of:

a) developing my understanding of blockchains
b) learning the Go programming language
c) give back to the bros

## Todo

- Collect and add transactions
- Once limit reached, keep collecting transactions anyway to mine later
- Mine + broadcast latest block
- When new block broadcast and verified, remove any transactions in pool which have already been mined
- PoW
- Longest chain wins
- Enforce fairness
   - when spending coins from one address, must sign transation w/ private key of that address
   - confirm transaction validity when adding to list
   - ensure all transactions in block are valid when new block received from network (https://en.bitcoin.it/wiki/Protocol_rules#.22tx.22_messages)
- Peer discovery
- Separate block header
   - Fixed size block header (only the header is hashed + mined, to ensure consistent difficulty regardless of # transactions in a block)
   - Transactions stored in merkel tree in block body
   - Store Merkel root in Block header only 
   - Build Merkle tree when collecting transactions
   - Verify transaction included in block using only necessary Merkel branches in transaction tree
   - Set block size limit
- Transaction fees
