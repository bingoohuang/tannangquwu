# 探囊取物

《三国演义·第四二回》：「我向曾闻云长言，翼德于百万军中，取上将之首，如探囊取物。」

本机docker MySQL一百万张卡，取可用卡，无池化时(每次请求访问数据库)TPS可以达到5百，池化时TPS则可以达到4万。

1. 卡表(card)：数据表（百万/千万级别），使用自增ID作为主键/唯一索引
2. 序号表(seq)：控制表（有限行），维护当前消费到的最大序号和可用最大序号
3. 消费步骤：序号表叫号，以消费序号去卡表中占用，消费

1. 优势1：避免低效查询(state=0)，以及查询出来的结果被竞争消费
2. 优势2：实现简单高效(无复杂的查询语句，只是最简单的带主键的改查)，池化处理后效率非常高

![image](https://user-images.githubusercontent.com/1940588/92542108-7ae54880-f27a-11ea-8fd2-21aadf9984c9.png)

验证步骤：

1. 招募百万兵甲
    ```bash
    🕙[2020-09-08 10:30:43.202] ❯ tannangquwu gen -r 1000000
    2020/09/08 10:30:57 Generating records by config &{NumRecs:1000000 BatchSize:1000 LogSeconds:10}
    2020/09/08 10:30:57 Starting inserts
    2020/09/08 10:31:07  526999/1000000 ( 52.70%) written in 10.000039601s, avg: 18.975µs/record, 52699.69 records/s
    2020/09/08 10:31:17  919999/1000000 ( 92.00%) written in 20.000474198s, avg: 21.739µs/record, 45998.86 records/s
    2020/09/08 10:31:20 1000000/1000000 (100.00%) written in 23.176722212s, avg: 23.176µs/record, 43146.74 records/s
    ```
1. 开启探囊取物服务
    ```bash
    🕙[2020-09-08 11:50:56.164] ❯ tannangquwu http
    2020/09/08 12:00:14 探囊取物服务开启 &{addr::8000 ctx:<nil> db:<nil>}
    ```
1. 取上将首级，如同探囊取物
    ```bash
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
    ```
1. "取上将首级"所使用SQL
    ```sql
    update seq set num = num + 1 where name = '兵甲';
    set @num = (select num from seq where name = '兵甲');
    update card set state = 1 where id = @num and state = 0;
    select @num;
    ```
1. 池化处理之后:
    ```bash
   🕙[2020-09-08 13:43:26.938] ❯ gobench -u http://127.0.0.1:8000 -rr 10000
   Dispatching 100 goroutines
   Waiting for results...

   Total Requests:			10000 hits
   Successful requests:		10000 hits
   Network failed:			0 hits
   Bad requests(!2xx):		0 hits
   Successful requests rate:	44323 hits/sec
   Read throughput:		5.9 MiB/sec
   Write throughput:		3.6 MiB/sec
   Test time:			225.616011ms
   ```