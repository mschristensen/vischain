const logger = require('winston');
const Joi = require('joi');
const Request = require('../utils/request.js');
const Response = require('../utils/response.js');
const Validator = require('../utils/validator.js');
const Rx = require('rxjs/Rx');
const State = require('../utils/state');
const utils = require('../utils/utils');

module.exports = function BlockController(req, res, next) {
    return {
        sendBlock: async () => {
            try {
                await new Validator().Address(req.body.originalSender);
                await new Validator().Recipients(req.body.recipients);
                await new Validator().Block(req.body.data);
            } catch (err) {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                logger.error(err);
                return Response.InternalServerError().send(res);
            }
            let all = [];
            for (let recipient of req.body.recipients) {
                all.push(
                    Rx.Observable.fromPromise(
                        new Request(recipient).SendBlock(req.body.data)
                    )
                );
            }
            State.Send('blocks', req.body);
            await State.SyncNodeChain(utils.getSenderAddressFromRequest(req));    // sender's chain has updated
            try {
                let responses = await Rx.Observable.forkJoin(...all).toPromise();
                result = {};
                for (let i in req.body.recipients) {
                    result[req.body.recipients[i]] = responses[i].data
                    await State.SyncNodeChain(req.body.recipients[i]);    // recipient's chain may have updated
                }
                State.Receive('blocks', req.body);
                return Response.OK(result).send(res);
            } catch (err) {
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                logger.error(err);
                return Response.InternalServerError().send(res);
            }
        }
    };
}