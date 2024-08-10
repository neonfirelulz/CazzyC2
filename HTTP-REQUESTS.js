var http = require('http')
  , parseURL = require('url').parse
  ;
function get(url, cb) {
  var req = http.request(parseURL(url), function(response) {
    var bodyParts = [];
    response.on('data', function(b) { bodyParts.push(b); });
    response.on('end', function() {
      var body = bodyParts.join('');
      if (response.statusCode === 200) {
        cb(null, body, response);
      }
      else {
        var err = new Error('http error');
        err.statusCode = response.statusCode;
        err.body = body;
        cb(err);
      }
    });
  });
  req.on('error', cb);
  req.end();
}

// Usage

get('http://google.ca', function(err, body, response) {
  if (err) {
    // error
  }
  else {
    // success
  }
});
