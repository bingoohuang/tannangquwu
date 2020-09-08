# 探囊取物

《三国演义·第四二回》：「我向曾闻云长言，翼德于百万军中，取上将之首，如探囊取物。」

1. 招募百万兵士

    🕙[2020-09-08 10:30:43.202] ❯ tannangquwu gen -r 1000000
    2020/09/08 10:30:57 Generating records by config &{NumRecs:1000000 BatchSize:1000 LogSeconds:10}
    2020/09/08 10:30:57 Opening database
    2020/09/08 10:30:57 Starting progress logging
    2020/09/08 10:30:57 Starting inserts
    2020/09/08 10:31:07  526999/1000000 ( 52.70%) written in 10.000039601s, avg: 18.975µs/record, 52699.69 records/s
    2020/09/08 10:31:17  919999/1000000 ( 92.00%) written in 20.000474198s, avg: 21.739µs/record, 45998.86 records/s
    2020/09/08 10:31:20 1000000/1000000 (100.00%) written in 23.176722212s, avg: 23.176µs/record, 43146.74 records/s

1. 开启探囊取物服务

    🕙[2020-09-08 11:50:56.164] ❯ tannangquwu http                              
    2020/09/08 12:00:14 探囊取物 &{addr::8000 ctx:<nil> db:<nil>}
    2020/09/08 12:00:14 Opening database

1. 取上将首级，如同探囊取物

    🕙[2020-09-08 11:59:21.401] ❯ gobench -u http://127.0.0.1:8000 -rr 1000
    Dispatching 100 goroutines
    Waiting for results...
    
    Total Requests:			1000 hits
    Successful requests:		1000 hits
    Network failed:			0 hits
    Bad requests(!2xx):		0 hits
    Successful requests rate:	584 hits/sec
    Read throughput:		79 KiB/sec
    Write throughput:		49 KiB/sec
    Test time:			1.712166417s

1. 所使用SQL

    update seq set num = num + 1 where name = '步兵';
    update card set state = 1 where id = (select num from seq where name = '步兵') and state = 0;
    select num from seq where name = '步兵';
