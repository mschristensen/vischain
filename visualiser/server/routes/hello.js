'use strict';

const Response = require('../utils/response.js');
const logger = require('winston');
const HelloController = require('../controllers/hello');

module.exports = function(router) {
  router.route('/')
    .get((req, res) => {
      return HelloController(req, res).hello();
    })
    .post((req, res) => {
      console.log("Body", req.body);
      return Response.OK(req.body).send(res);
    });
};
