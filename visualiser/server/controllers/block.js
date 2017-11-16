const logger = require('winston');
const Request = require('../utils/request.js');
const Response = require('../utils/response.js');
const Validator = require('../utils/validator.js');
const Rx = require('rxjs/Rx');

module.exports = function BlockController(req, res, next) {
    return {
        receiveBlock: async (forwardTo) => {
            try {
                await new Validator().Block(req.body);
                let all = [];
                for (let peer of forwardTo) {
                    all.push(
                        Rx.Observable.fromPromise(
                            new Request(peer).SendBlock(req.body)
                        )
                    );
                }
                let responses = await Rx.Observable.forkJoin(...all).toPromise();
                result = {};
                for (let i in forwardTo) {
                    result[forwardTo[i]] = responses[i].data
                }
                return Response.OK(result).send(res);
            } catch (err) {
                if (err.isJoi) {
                    return Response.BadRequest().send(res);
                }
                if (err.code === 'ECONNREFUSED') {
                    return Response.BadGateway().send(res); 
                }
                logger.error(err);
                return Response.InternalServerError().send(res);
            };
        }
    };
}