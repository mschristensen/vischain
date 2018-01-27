'use strict';

const Response = require('../utils/response.js');
const logger = require('winston');
const ChainController = require('../controllers/chain');

module.exports = function(router) {
  router.route('/')
    .get(async (req, res) => {
        return ChainController(req, res).getChain();
    })
    .post(async (req, res) => {
        return ChainController(req, res).updateChain();
    });
};
