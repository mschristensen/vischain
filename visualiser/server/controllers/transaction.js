'use strict';

const logger = require('winston');
const Response = require('../utils/response.js');
const Request = require('../utils/request');
const Validator = require('../utils/validator');
const State = require('../utils/state');

module.exports = function TransactionController(req, res, next) {
    return {
        sendTransaction: async (params) => {
            try {
                await new Validator().Transaction(req.body);
                await State.StartTransaction(req.body);
                await State.Emit();
                let result = await new Request(req.body.recipient).SendTransaction(req.body);
                State.StopTransaction(req.body);
                await State.Emit();
                return Response.OK(result.data).send(res);
            } catch (err) {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                logger.error(err);
                return Response.InternalServerError(err).send(res);
            };
        }
    };
}