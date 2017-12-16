'use strict';

const fs = require('fs');
const Validator = require('./validator');
const sm = require('./socket');
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
 *      }, ...]
 *  }
 */


class State {

    constructor() {
        this.state = {
            network: [],        // defines network topology
            transactions: [],   // transactions currently in transit
            blocks: []          // blocks currently in transit
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