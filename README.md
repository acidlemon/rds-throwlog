# RDS ThrowLog

Retrive slowlog, errorlog or general log etc... from Amazon RDS and send to fluentd.

## usage

Prepare your ~/.aws/credentials.

And run rds-throwlog.

```
$ rds-throwlog --database="my-database" --path="slowquery/mysql-slowquery.log"
```

Enjoy! (?)

## License

The MIT License (MIT)

(c) 2014 acidlemon. (c) 2014 KAYAC Inc.


