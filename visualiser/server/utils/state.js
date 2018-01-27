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
                await this.SyncNodeChain(node.address);
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
        let response;
        try {
            await new Validator().Address(address);
            response = await new Request(address).GetChain();
        } catch (err) {
            return Promise.reject(err);
        }
        let idx = this.state.chains.map(chain => chain.address).indexOf(address);
        let newChain = {
            address,
            chain: response.data.payload
        };
        if (idx > -1) {
            this.state.chains[idx] = newChain;
        } else {
            this.state.chains.push(newChain);
        }
        console.log('UPDATED ', address, response.data.payload)
        return Promise.resolve();
    }

    Send(what, data) {
        if (this.state[what].indexOf(data) === -1) {
            this.state[what].push(data);
            this.Emit();
        }
    }

    Receive(what, data, recipientAddr) {
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