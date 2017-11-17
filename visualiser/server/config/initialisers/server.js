'use strict';

const express = require('express');
const bodyParser = require('body-parser');
const morgan = require('morgan');
const logger = require('winston');
const helmet = require('helmet');
let app;

module.exports = function() {
  return new Promise((resolve, reject) => {
    app = express();

    // Security
    app.use(helmet());
    
    // Log HTTP requests
    app.use(morgan('common'));

    // Parse request bodies
    app.use(bodyParser.urlencoded({ extended: true }));
    app.use(bodyParser.json({ type: '*/*' }));

    // CROSS ORIGIN RESOURCE SHARING
    app.use((req, res, next) => {
      const allowedOrigins = ['http://localhost:3000'];
      const origin = req.headers.origin;
      if (allowedOrigins.indexOf(origin) > -1) {
        res.setHeader('Access-Control-Allow-Origin', origin);
      }
      res.header("Access-Control-Allow-Credentials", "true");
      res.header('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE,PATCH,OPTIONS');
      res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Authorization, Accept");
      next();
    });

    logger.info('[SERVER] Initializing routes');
    require('../../routes/index')(app);

    const server = app.listen(process.env.PORT || 3001);
    logger.info('[SERVER] Listening on port ' + (process.env.PORT || 3001));

    // initialisations
    require('../../utils/socket').Init(server);
    require('../../utils/state').Init();

    return resolve(app);
  });
};