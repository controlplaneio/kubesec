const AWS = require('aws-sdk');

if (!process.env.LAMBDA_TASK_ROOT) {
  const credentials = new AWS.SharedIniFileCredentials({profile: 'binslug-s3'});
  AWS.config.credentials = credentials;
}

AWS.config.logger = process.stdout;

const firehose = new AWS.Firehose({region: 'us-east-1'});
const streamName = "kubesec-staging"
const fs = require('fs')

const filename = process.argv[2]
const record = {
  filename: filename,
  content: fs.readFileSync(filename, 'utf8'),
  created_at: (new Date()).toISOString().substr(0, 19).replace('T', ' ')
}

// http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/Firehose.html#putRecord-property
function putRecord(dStreamName, data, callback) {
  var recordParams = {
    Record: {
      Data: JSON.stringify(data) + '\n'
    },
    DeliveryStreamName: dStreamName
  };

  firehose.putRecord(recordParams, callback);
}

putRecord(streamName, record, function (err, res) {
  if (err) throw new Error(err);
  console.log('done', res)
})
