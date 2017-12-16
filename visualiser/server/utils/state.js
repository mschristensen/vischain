'use strict';

const fs = require('fs');
const Validator = require('./validator');
const sm = require('./socket-manager');
const logger = require('winston');
/**
 *  State schema:
 *  {
 *      network: [{
 *          address: "8080",
 *          peers: ["8081", "8082"]
 *      }, {
 *          address: "8081",
 *          peers: ["8080", "8082"]
 *      }, ...],
 *      broadcasts: [{
 *          type: 'transaction',
 *          data: {
 *              sender: "8080",
 *              recipient: "8081",
 *              amount: 1
 *          }
 *      }]
 *  }
 */


class State {

    constructor() {
        this.state = {
            network: [],
            transactions: []
        };
    }

    // socket won't be ready in here!
    async Init() {
        await this.ConfigureNetwork();
        sm.Socket().on('connect', (socket) => {
            this.socket = socket;
            this.Emit();
        });
    }

    async ConfigureNetwork() {
        return new Promise((resolve, reject) => {
            // read network configuration file and init state
            // (path relative to where `node` process was initiated, should be in `visualiser`)
            fs.readFile('../network.config', 'utf8', async (err, data) => {
                if (err) {
                    return reject(err);
                }
                let network = [];
                data = data.split('\n');
                for (let line of data) {
                    line = line.split(' ');
                    network.push({
                        address: line[0].trim(),
                        peers: line.slice(1).map(address => address.trim())
                    });
                }
                try {
                    await new Validator().Network(network);
                } catch (err) {
                    return reject(err);
                }
                this.state.network = network;
                resolve();
            });
        });
    }

    async StartTransaction(transaction) {
        return new Promise(async (resolve, reject) => {
            try {
                await new Validator().Transaction(transaction);
            } catch (err) {
                return reject(err);
            }
            if (this.state.transactions.indexOf(transaction) === -1) {
                this.state.transactions.push(transaction);
            }
            return resolve();
        });
    }

    StopTransaction(transaction) {
        if (this.state.transactions.indexOf(transaction) > -1) {
            this.state.transactions.splice(this.state.transactions.indexOf(transaction), 1);
        }
    }

    Emit() {
        if (this.socket) {
            this.socket.emit('stateUpdate', this.state);
        } else {
            logger.error('Tried to use socket before it was initialised!');
        }
    }
};

const state = new State();
module.exports = state;