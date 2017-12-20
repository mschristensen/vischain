'use strict';

const fs = require('fs');
const Validator = require('./validator');
const sm = require('./socket-manager');
const logger = require('winston');
const Request = require('./request');

class State {

    constructor() {
        this.state = {
            topology: [],       // defines network topology
            transactions: [],   // transactions currently in transit
            blocks: [],         // blocks currently in transit
            chains: [],         // local blockchain for each node
        };
    }

    // socket won't be ready in here!
    async Init() {
        try {
            await this.ConfigureNetworkTopology();
            for (let node of this.state.topology) {
                let response = await this.SyncNodeChain(node.address);
                this.state.chains.push({
                    address: node.address,
                    chain: response.data.payload
                });
            }
        } catch (err) {
            logger.error(err);
        }
        sm.Socket().on('connect', (socket) => {
            this.socket = socket;
            this.Emit();
        });
    }

    async ConfigureNetworkTopology() {
        return new Promise((resolve, reject) => {
            // read network configuration file and init state
            // (path relative to where `node` process was initiated, should be in `visualiser`)
            fs.readFile('../network.config', 'utf8', async (err, data) => {
                if (err) {
                    return reject(err);
                }
                let topology = [];
                data = data.split('\n');
                for (let line of data) {
                    line = line.split(' ');
                    topology.push({
                        address: line[0].trim(),
                        peers: line.slice(1).map(address => address.trim())
                    });
                }
                try {
                    await new Validator().Topology(topology);
                } catch (err) {
                    return reject(err);
                }
                this.state.topology = topology;
                resolve();
            });
        });
    }

    async SyncNodeChain(address) {
        try {
            await new Validator().Address(address);
        } catch (err) {
            return Promise.reject(err);
        }
        return new Request(address).GetChain();
    }

    Send(what, data) {
        if (this.state[what].indexOf(data) === -1) {
            this.state[what].push(data);
            this.Emit();
        }
    }
    
    Receive(what, data) {
        if (this.state[what].indexOf(data) > -1) {
            this.state[what].splice(this.state[what].indexOf(data), 1);
            this.Emit();
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