'use strict';

const Response = require('../utils/response.js');
const logger = require('winston');
const BlockController = require('../controllers/block');

module.exports = function(router) {
  router.route('/')
    .post(async (req, res) => {
        return BlockController(req, res).sendBlock();
    });
};
