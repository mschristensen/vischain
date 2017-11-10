'use strict';

const Response = require('../utils/response.js');
const logger = require('winston');
const BlockController = require('../controllers/block');

module.exports = function(router) {
  router.route('/')
    .post((req, res) => {
        // TODO validate peers
        return BlockController(req, res).receiveBlock(req.query.peers.split(','));
    });
};
