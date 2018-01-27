'use strict';

const changeCase = require('change-case');
const express = require('express');
const routes = require('require-dir')();  // requires all other files in this directory
const Response = require('../utils/response.js');

module.exports = function (app) {
    // Initialize all routes
    Object.keys(routes).forEach((routeName) => {
        let router = express.Router();

        // Initialize the route
        require('./' + routeName)(router);

        app.use((req, res, next) => { console.log(req.body); next(); })

        // Tie the router to it's url path
        app.use('/api/v1/' + changeCase.paramCase(routeName), router);
    });

    // Catch unknown API endpoints as 404
    app.all('/api/v1/*', function (req, res, next) {
        return Response.NotFound().send(res);
    });
};
