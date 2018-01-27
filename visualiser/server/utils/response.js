'use strict';

const logger = require('winston');

module.exports = class Response {
    constructor(statusCode, headers, payload) {
        if (typeof statusCode !== 'number') {
            throw new TypeError('statusCode must be a number but was ' + (typeof statusCode));
        }
        if (headers === null) {
            throw new TypeError('headers cannot be null');
        }
        if (typeof headers !== 'object') {
            throw new TypeError('headers must be an object but was ' + (typeof headers));
        }
        this.statusCode = statusCode;
        this.headers = headers;
        this.payload = payload;
    }

    send(res) {
        let response = this.payload || null;
        res.statusCode = this.statusCode;
        for (let key in this.headers) {
            res.set(key.toLowerCase(), this.headers[key]);
        }
        return res.send(response);
    }

    // 2XX  SUCCESS
    static OK(payload) {
        return new Response(200, { 'Content-Type': 'application/json' }, payload);
    }
    static Created(payload) {
        return new Response(201, { 'Content-Type': 'application/json' }, payload);
    }

    // 3XX  REDIRECTION
    static MovedPermanently(payload) {
        return new Response(301, { 'Content-Type': 'application/json' }, payload);
    }

    // 4XX  CLIENT ERROR
    static BadRequest(payload) {
        return new Response(400, { 'Content-Type': 'application/json' }, payload);
    }
    static Forbidden(payload) {
        return new Response(403, { 'Content-Type': 'application/json' }, payload);
    }
    static NotFound(payload) {
        return new Response(404, { 'Content-Type': 'application/json' }, payload);
    }
    static MethodNotAllowed(payload) {
        return new Response(405, { 'Content-Type': 'application/json' }, payload);
    }
    static UnsupportedMediaType(payload) {
        return new Response(415, { 'Content-Type': 'application/json' }, payload);
    }
    static TooManyRequests(payload) {
        return new Response(429, { 'Content-Type': 'application/json' }, payload);
    }

    // 5XX SERVER ERROR
    static InternalServerError(payload) {
        return new Response(500, { 'Content-Type': 'application/json' }, payload);
    }
    static NotImplemented(payload) {
        return new Response(501, { 'Content-Type': 'application/json' }, payload);
    }
    static BadGateway(payload) {
        return new Response(502, { 'Content-Type': 'application/json' }, payload);
    }
    static ServiceUnavailable(payload) {
        return new Response(503, { 'Content-Type': 'application/json' }, payload);
    }

    // CUSTOM ERRORS
    static MongooseError(err) {
        // Transform common mongoose errors like ValidationError into a formatted response object
        if (err.name === 'ValidationError') {
            let invalids = [];
            for (let invalid in err.errors) {
                invalids.push(invalid);
            }
            return Response.BadRequest({ InvalidArguments: invalids });
        } else {
            logger.error('Mongoose Error:', err);
            return Response.InternalServerError();
        }
    }

    static InvalidArguments(args) {
        return Response.BadRequest({ InvalidArguments: args });
    }
};
