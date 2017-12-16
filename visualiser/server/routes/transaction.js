'use strict';

const Response = require('../utils/response.js');
const logger = require('winston');
const TransactionController = require('../controllers/transaction');

module.exports = function(router) {
  router.route('/')
    .post((req, res) => {
        return TransactionController(req, res).sendTransaction();
    });
};
