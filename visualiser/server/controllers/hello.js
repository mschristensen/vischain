const Response = require('../utils/response.js');

module.exports = function HelloController(req, res, next) {
  return {
    hello: (params) => Response.OK('Hello, world!').send(res)
  };
}