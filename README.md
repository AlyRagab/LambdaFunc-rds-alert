# LambdaFunc-rds-alert
Golang Lambda Function for sending RDS Metrics to Slack, The Flow as below :

<img src="https://github.com/AlyRagab/LambdaFunc-rds-alert/blob/main/slack.png" />

#### Steps :
- Create the RDS instance :
```
aws rds create-db-instance \
    --db-instance-identifier $INSTANCE_NAME \
    --db-instance-class $CLASS_NAME \
    --engine mysql \
    --master-username $USER_NAME \
    --master-user-password $PASSWORD \
    --allocated-storage 50
```

- Create the SNS topic :
```
aws sns create-topic --name rds
```
- Create the CloudWatch Alarm :
```
aws cloudwatch put-metric-alarm --alarm-name "HighCPUUtilization"  --alarm-description "High CPU Utilization" \
--actions-enabled  \
--alarm-actions "$ROLE_ARN"  \
--metric-name "CPUUtilization"  \
--namespace AWS/RDS --statistic "Average" \
--period 60  \
--evaluation-periods 60  \
--threshold 70  \
--comparison-operator "GreaterThanOrEqualToThreshold"
```

- Build it :
```
go build -o main main.go
zip deployment.zip main
```

- Create the Lambda Function :
```
aws lambda create-function --function-name rds-slack \
--zip-file fileb://./deployment.zip \
--runtime go1.x --handler main \
--role $ROLE_ARN \
--environment Variables={SLACK_WEBHOOK=""}
```