'use strict';

const Joi = require('joi');
const Response = require('../utils/response.js');
const logger = require('winston');
const BlockController = require('../controllers/block');

module.exports = function(router) {
  router.route('/')
    .post((req, res) => {
        let peers = [];
        const validate = () => {
            if (!req.query.peers) {
                return Promise.reject();
            }
            peers = req.query.peers.split(',');
            return Joi.validate(peers, Joi.array().min(1).unique().items(Joi.string()));
        }

        validate().then(() => {
            return BlockController(req, res).receiveBlock(peers);
        }).catch(err => {
            return Response.BadRequest().send(res);
        });
    });
};
