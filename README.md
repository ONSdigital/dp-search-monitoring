dp-search-monitoring
==================

### Configuration

| Environment variable  | Default          | Type    | Description
| --------------------  | ---------------- | ------- | ------------------------------------------------
| AWS_ACCESS_KEY_ID     | N/A              | String  | AWS access key.
| AWS_SECRET_ACCESS_KEY | N/A              | String  | AWS secret access key.
| AWS_REGION            | N/A              | String  | AWS region to use.
| ANALYTICS_SQS_URL     | N/A              | String  | URL of SQS queue to use.
| SQS_WAIT_TIMEOUT      | 20               | int64   | Timeout in seconds.
| MONGODB_URL           | localhost:27017  | String  | MongoDB URL.
| MONGO_DB              | local            | String  | Database to use in MongoDB.
| MONGO_COLLECTION      | searchstats      | String  | Collection to use in MongoDB.
| RUN_ON_STARTUP        | true             | Boolean | Run import on startup.
| TIME_UNIT             | DAYS             | String  | Schedule one import per day (DAYS) or hour (HOURS).
| AT_TIME               | 00:00            | String  | Time to run job (when using TIME_UNIT='DAYS' only).

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
