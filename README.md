dp-search-monitoring
==================

### Configuration

| Environment variable  | Default          | Description
| --------------------  | -----------------  -------------------------
| AWS_ACCESS_KEY_ID     | N/A              | AWS access key.
| AWS_SECRET_ACCESS_KEY | N/A              | AWS secret access key.
| AWS_REGION            | N/A              | AWS region to use.
| ANALYTICS_SQS_URL     | N/A              | URL of SQS queue to use.
| SQS_WAIT_TIMEOUT      | 20               | Timeout in seconds.
| MONGODB_URL           | localhost:27017  | MongoDB URL.
| MONGO_DB              | local            | Database to use in MongoDB.
| MONGO_COLLECTION      | searchstats      | Collection to use in MongoDB.
| RUN_ON_STARTUP        | true             | (bool) Whether to run import on startup.
| TIME_UNIT             | DAYS             | Schedule one import per day (DAYS) or hour (HOURS).
| AT_TIME               | 00:00            | Time to run job (when using TIME_UNIT='DAYS' only).

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
