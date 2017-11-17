'use strict';

const fs = require('fs');
const Validator = require('./validator');

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
        this.state = {};
    }

    async Init() {
        await this.ConfigureNetwork();
    }

    async ConfigureNetwork() {
        return new Promise((resolve, reject) => {
            // read network configuration file and init state
            fs.readFile('network.config', 'utf8', async (err, data) => {
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
                    await new Validator().Network(network)
                } catch (err) {
                    return reject(err);
                }
                this.state.network = network;
                resolve();
            });
        });
    }
};

const state = new State();
state.Init();
module.exports = state;