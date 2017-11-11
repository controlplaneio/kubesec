var QUEUE_URL = 'https://sqs.us-east-1.amazonaws.com/{AWS_ACCUOUNT_}/matsuoy-lambda';
var AWS = require('aws-sdk');
var sqs = new AWS.SQS({region: 'us-east-1'});

var enqueue = function (filename, context) {
  var params = {
    MessageBody: JSON.stringify({
      name: filename,
      content: fs.readFileSync(filename)
    }),
    QueueUrl: QUEUE_URL
  };
  sqs.sendMessage(params, function (err, data) {
    if (err) {
      console.log('error:', "Fail Send Message" + err);
    } else {
      console.log('data:', data.MessageId);
    }
  });
}

filename = "/tmp/test"
enqueue(filename)
