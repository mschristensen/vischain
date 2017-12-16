'use strict';

const logger = require('winston');
const Response = require('../utils/response.js');
const Request = require('../utils/request');
const Validator = require('../utils/validator');
const State = require('../utils/state');

module.exports = function TransactionController(req, res, next) {
    return {
        sendTransaction: async () => {
            try {
                await new Validator().Transaction(req.body);
            } catch (err) {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                logger.error(err);
                return Response.InternalServerError().send(res);
            }
            State.Send('transactions', req.body);
            try {
                let result = await new Request(req.body.recipient).SendTransaction(req.body);
                State.Receive('transactions', req.body);
                return Response.OK(result.data).send(res);
            } catch (err) {
                State.Receive(req.body);
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                logger.error(err);
                return Response.InternalServerError(err).send(res);
            }
        }
    };
}