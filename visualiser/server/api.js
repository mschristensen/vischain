// TODO

'use strict';

const logger = require('winston');
const axios = require('axios');

class APIController {
    static SendTransaction(transaction) {
    
    }
}

module.exports = class Request {
  constructor(method, url, body, addr) {
    if (['POST', 'PUT', 'GET', 'DELETE'].indexOf(method) === -1) {
        throw new TypeError('unknown HTTP method');  
    }
    
    this.method = method;
    this.url = url;
    this.body = body;
    this.addr = addr;
  }

  send() {
    if (!this.addr) {
        return Promise.reject();
    }
    return axios({
        method: this.method,
        url: this.url,
        headers: {'Content-Type': 'application/json'},
        data: this.body
    });
  }

  static to(addr) {
      return new Request('POST', '/transaction', transaction, addr);
  }
};
