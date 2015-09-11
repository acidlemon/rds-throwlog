# RDS ThrowLog

Retrive slowlog, errorlog or general log etc... from Amazon RDS and send to stdout(JSON) or fluentd.

## usage

Prepare your ~/.aws/credentials.

And run rds-throwlog.

```
$ rds-throwlog --database="my-database" --path="slowquery/mysql-slowquery.log"
```

### options

rds-throwlog outputs to stdout as JSON format.

- output raw file

Use `--raw` option.

```
$ rds-throwlog --database="my-database" --path="slowquery/mysql-slowquery.log" \
  --raw
```

- send to fluentd

Use `--fluent-host`, `--fluent-port`, `--fluent-tag` options.

```
$ rds-throwlog --database="my-database" --path="slowquery/mysql-slowquery.log" \
  --fluent-host=localhost --fluent-port=24224 --fluent-tag=mysql.slowquery
```

- `--fluent-host`: Fluent hostname (required to send to fluentd)
- `--fluent-port`: Fluent port (Default: 24224)
- `--fluent-tag`: Fluent tag (Default: "mysql.slowquery")

Enjoy! (?)

## License

The MIT License (MIT)

(c) 2014 acidlemon. (c) 2014 KAYAC Inc.


