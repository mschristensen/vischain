'use strict';

const logger = require('winston');
const Joi = require('joi');

module.exports = class Validator {

    constructor() {
        this.Schemas = {}

        // primitives
        this.Schemas.NumericString = Joi.string().regex(/^\d+$/);

        // custom
        this.Schemas.Address = this.Schemas.NumericString;
        this.Schemas.Transaction = Joi.object().keys({
            sender: this.Schemas.Address,
            recipient: this.Schemas.Address,
            amount: this.Schemas.NumericString
        });
        this.Schemas.Block = Joi.object().keys({
            index: this.Schemas.NumericString,
            timestamp: this.Schemas.NumericString,
            transactions: Joi.array().items(this.Schemas.Transaction),
            proof: Joi.string().length(8),
            prevHash: Joi.string().length(44)
        });
    }

    async Address(data) { return Joi.validate(data, this.Schemas.Address); }
    async Transaction(data) { return Joi.validate(data, this.Schemas.Transaction); }
    async Block(data) { return Joi.validate(data, this.Schemas.Block); }
};