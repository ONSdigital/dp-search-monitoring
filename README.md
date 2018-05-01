dp-search-monitoring
==================

### General Configuration

| Environment variable  | Type    | Default          | Description
| --------------------  | ------- | ---------------- | ------------------------------------------------------------------
| AWS_ACCESS_KEY_ID     | String  | N/A              | AWS access key.
| AWS_SECRET_ACCESS_KEY | String  | N/A              | AWS secret access key.
| AWS_REGION            | String  | N/A              | AWS region to use.
| RUN_ON_STARTUP        | Boolean | true             | Run import on startup.
| TIME_UNIT             | String  | DAYS             | Schedule imports by day (DAYS), hour (HOURS) or minutes (MINS).
| TIME_WINDOW           | uint64  | 1                | Time window for job to run (i.e every 'x' DAYS).
| AT_TIME               | String  | 00:00            | Time to run job (when using TIME_UNIT='DAYS' only).
| BACKEND               | String  | RDS_POSTGRES     | Which backend to use (RDS_POSTGRES or MONGO).

### SQS Configuration

| Environment variable  | Type    | Default          | Description
| --------------------  | ------- | ---------------- | ------------------------------------------------------------------
| ANALYTICS_SQS_URL     | String  | N/A              | URL of SQS queue to use.
| SQS_WAIT_TIMEOUT      | int64   | 20               | Timeout in seconds (must be between 0 and 20).
| SQS_DELETE_ENABLED    | Boolean | true             | Specify if deletion of SQS messages is permitted.

### MongoDB Configuration

| Environment variable  | Type    | Default          | Description
| --------------------  | ------- | ---------------- | ------------------------------------------------------------------
| MONGODB_URL           | String  | localhost:27017  | MongoDB URL.
| MONGO_DB              | String  | local            | Database to use in MongoDB.
| MONGO_COLLECTION      | String  | searchstats      | Collection to use in MongoDB.

### AWS RDS Postrges Configuration

| Environment variable  | Type    | Default          | Description
| --------------------  | ------- | ---------------- | ------------------------------------------------------------------
| RDS_DB_USER           | String  | N/A              | Postgres DB username.
| RDS_DB_PASSWORD       | String  | N/A              | Postgres DB password.
| RDS_DB_NAME           | String  | N/A              | Postgres DB name.
| RDS_DB_ENDPOINT       | String  | N/A              | AWS Endpoint of Postgres DB.
| RDS_PORT              | int     | 5432             | Port to connect to.

### Licence

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
