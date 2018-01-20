const logger = require('winston');
const Joi = require('joi');
const Request = require('../utils/request.js');
const Response = require('../utils/response.js');
const Validator = require('../utils/validator.js');
const Rx = require('rxjs/Rx');
const State = require('../utils/state');
const utils = require('../utils/utils');

module.exports = function ChainController(req, res, next) {
    return {
        getChain: async () => {
            try {
                await new Validator().Address(req.query.peer);
                await new Validator().Hash(req.query.lastBlockHash);
            } catch (err) {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                logger.error(err);
                return Response.InternalServerError().send(res);
            }
            try {
                let result = await new Request(req.query.peer).GetChain(req.query.lastBlockHash);
                return Response.OK(result.data).send(res);
            } catch (err) {
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                logger.error(err);
                return Response.InternalServerError(err).send(res);
            }
        }
    };
}