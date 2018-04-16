dp-search-monitoring
==================

### Configuration

| Environment variable  | Description
| --------------------  | --------------------------------------
| AWS_ACCESS_KEY_ID     | AWS access key.
| AWS_SECRET_ACCESS_KEY | AWS secret access key.
| AWS_REGION            | AWS region to use.
| ANALYTICS_SQS_URL     | URL of SQS queue to use.
| SQS_WAIT_TIMEOUT      | Timeout in seconds.
| MONGODB_URL           | MongoDB URL.
| MONGO_DB              | Database to use in MongoDB.
| MONGO_COLLECTION      | Collection to use in MongoDB.
| RUN_ON_STARTUP        | Boolean - whether to run import on startup.
| TIME_UNIT             | Schedule one import per day (DAYS) or hour (HOURS).
| AT_TIME               | Time to run job (when using TIME_UNIT='DAYS' only).

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
