# Vischain

A visualisation of a toy blockchain network.

## Todo `blockchain`

- Fix node crash, sometimes hangs
- If an incoming block is verified and added to the chain, remove any transactions in pool which have already been mined and start again
- If a longer chain is selected, ensure we pool any unmined transactions which were in our "lost" blocks to mine later
- Handle BadGateway errors (offline peers): manage peer connections centrally on API
- Choice of whether to be a miner or not
- Implement a block size limit
- Once block size limit reached, keep collecting transactions anyway to mine later
- Sign transations
- Confirm transaction validity when adding to list
- Ensure all transactions in block are valid when new block received from network (https://en.bitcoin.it/wiki/Protocol_rules#.22tx.22_messages)
- Peer discovery
- Separate block header
   - Fixed size block header (only the header is hashed + mined, to ensure consistent difficulty regardless of # transactions in a block)
   - Transactions stored in merkel tree in block body
   - Store Merkle root in Block header only
   - Build Merkle tree when collecting transactions
   - Verify transaction included in block using only necessary Merkle branches in transaction tree
   - Set block size limit
- Transaction fees
- Implement UXTO + balance system https://www.cryptocompare.com/mining/guides/bitcoin-transaction-inputs-and-outputs/

## Todo `visualiser`

- Visualise:
    - Block mined
    - Transaction rejected/accepted
    - Block rejected/accepted
    - Chain divergences
    - Balances
- Edit:
    - Difficulty
    - Nodes in network + connections
    - Which nodes are miners


## Protocol

## API

All requests proxied from node A to node B will forward the `X-Sender` header in the request to B.

### `POST /transaction`

Transaction will be sent to `recipient` and their response returned in `payload`.

Example request:

```
POST /transaction
{
	"sender": "8081",
	"recipient": "8080",
	"amount": "1"
}
```

Response example:

```
{
    "payload": {
        "Code": 1
    },
    "title": "OK"
}
```

### `POST /block`

Block in `data` will be sent to all `recipients` and responses returned in `payload`.

Example request:

```
POST /block
{
    originalSender: "8080",
    recipients: ["8081", "8082"],
    data: {
        "index": "5",
        "timestamp": "1510332444551936900",
        "transactions": [{
            "sender": "8081", "recipient": "8080", "amount": "1"
        },{
            "sender": "8081", "recipient": "8080", "amount": "1"
            
        }],
        "proof": "T5AEAA==",
        "prevHash": "tDw6oL/3BxXN+pTY6o/8M6eVBcKsUow3YTQgl88BscY="
    }
}
```

Response example:

```
{
    "payload": {
        "8081": {
            Code: 1
        },
        "8082": {
            Code: 1
        }
    },
    "title": "OK"
}
```

### `GET /chain`

Fetch a chain of blocks from a `peer`.

Example request:

```
GET /chain?peer=8081&lastBlockHash=tDw6oL/3BxXN+pTY6o/8M6eVBcKsUow3YTQgl88BscY=
```

`// TODO: Response example`


## Nodes

### `POST /transaction`
Example request:

```
POST /transaction
{
	"sender": "8081",
	"recipient": "8080",
	"amount": "1"
}
```

Response example:
```
{
    Code: 1
}
```

Response codes:
Code | Description
--- | ---
1 | The transaction was received successfully
2 | The request body could not be parsed
3 | The transaction could not be parsed from the body

### `POST /block`
Example request:

```
POST /block
{
	"index": "5",
	"timestamp": "1510332444551936900",
	"transactions": [{
		"sender": "8081", "recipient": "8080", "amount": "1"
	},{
    	"sender": "8081", "recipient": "8080", "amount": "1"
		
	}],
	"proof": "T5AEAA==",
	"prevHash": "tDw6oL/3BxXN+pTY6o/8M6eVBcKsUow3YTQgl88BscY="
}
```

Response example:
```
{
    Code: 1
}
```

Response codes:
Code | Description
--- | ---
1 | The block was received successfully
2 | The request body could not be parsed
3 | The block could not be parsed from the body

### `GET /chain`

Fetch a chain of blocks from this node. If a `lastBlockHash` is provided, return `Code: 1` with all the blocks on the chain after and not including the specified block.
If `lastBlockHash` is not found in the chain, the chain may have forked. Return `Code: 2` with the entire chain.

Example request:

```
GET /chain?lastBlockHash=tDw6oL%2F3BxXN%2BpTY6o%2F8M6eVBcKsUow3YTQgl88BscY%3D
```

Response example:
```
{
    Code: 1,
    Payload: [{
        "index": "2",
        "timestamp": "1510332444551936900",
        "transactions": [{
            "sender": "8081", "recipient": "8080", "amount": "1"
        },{
            "sender": "8081", "recipient": "8080", "amount": "1"
            
        }],
        "proof": "T5AEAA==",
        "prevHash": "tDw6oL/3BxXN+pTY6o/8M6eVBcKsUow3YTQgl88BscY="
    }, ...]
}
```

Response codes:
Code | Description
--- | ---
1 | The tail of the chain was returned
2 | The hash was not found, and the entire chain was returned
