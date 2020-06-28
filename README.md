# AWS Aurora Serverless
**cdk:** 
Creates an aurora serverless cluster (needs nodejs and CDK). The database needs to be manually initialised, the script is located in the lambda directory. You can use the AWS console to run aurora queries.

**client:**
A react app with a single table that calls the server.  (yarn start to start)

**server:**
A Golang server that connects to the cluster created via CDK. The servers require two environmental variables: AURORA_ARN and SECRET_ARN. 