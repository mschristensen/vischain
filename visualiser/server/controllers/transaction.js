'use strict';

const logger = require('winston');
const Response = require('../utils/response.js');
const Request = require('../utils/request');
const Validator = require('../utils/validator');

module.exports = function TransactionController(req, res, next) {
    return {
        receiveTransaction: (params) => {
            new Validator().Transaction(req.body).then(() => {
                return new Request(req.body.recipient).SendTransaction(req.body);
            }).then(result => {
                return Response.OK(result.data).send(res);
            }).catch(err => {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                return Response.InternalServerError(err).send(res);
            });
        }
    };
}