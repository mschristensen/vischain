'use strict';
const server = require('./server/config/initialisers/server');
const logger = require('winston');

// Load env vars from .env file into process.env
require('dotenv').load();

logger.info('[APP] Starting server initialization...');

// Start the server
module.exports = new Promise(function (resolve, reject) {
    server().then((app) => {
        logger.info('[APP] initialized SUCCESSFULLY');
        resolve(app);
    }).catch((err) => {
        logger.error('[APP] initialization failed', err);
        reject(err);
    });
});