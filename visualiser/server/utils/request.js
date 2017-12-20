'use strict';

const axios = require('axios');

const DOMAIN = `http://localhost:`;

module.exports = class Request {
    constructor(addr) {
        this.addr = addr;
    }

    send(method, url, data) {
        if (
            !method ||
            ['POST', 'PUT', 'GET', 'DELETE'].indexOf(method.toUpperCase()) === -1 ||
            !this.addr
        ) {
            return Promise.reject();
        }

        url = DOMAIN + this.addr + url;

        return axios({
            method,
            url,
            headers: { 'Content-Type': 'application/json' },
            data
        });
    }

    SendTransaction(transaction) { return this.send('POST', '/transaction', transaction); }
    SendBlock(block) { return this.send('POST', '/block', block); }
    GetChain(lastBlockHash) {
        let query = '';
        if (lastBlockHash) {
            lastBlockHash = encodeURIComponent(lastBlockHash);
            query = `?lastBlockHash=${lastBlockHash}`;
        }
        return this.send('GET', `/chain${query}`);
    }
};