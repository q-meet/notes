# mysql 日记

## mysql分组取最大(最小、最新、前N条)条记录

在数据库开发过程中，我们要为每种类型的数据取出前几条记录，或者是取最新、最小、最大等等，这个该如何实现呢，本文章向大家介绍如何实现mysql分组取最大(最小、最新、前N条)条记录。需要的可以参考一下。

先看一下本示例中需要使用到的数据

创建表并插入数据：

```sql
CREATE TABLE `tb` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) CHARACTER SET latin1 DEFAULT NULL,
  `val` int(11) DEFAULT NULL,
  `memo` varchar(20) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

insert into tb values('a',    2,   'a2');
insert into tb values('a',    1,   'a1');
insert into tb values('a',    3,   'a3');
insert into tb values('b',    1,   'b1');
insert into tb values('b',    3,   'b3');
insert into tb values('b',    2,   'b2');
insert into tb values('b',    4,   'b4');
insert into tb values('b',    5,   'b5');

```

### 按name分组取val最大的值所在行的数据

- 方法一：

```sql
select a.* from tb a where val = (select max(val) from tb where name = a.name) order by a.name
```

- 方法二：

```sql
select a.* from tb a where not exists(select 1 from tb where name = a.name and val > a.val)
```

- 方法三：

```sql
select a.* from tb a,(select name,max(val) val from tb group by name) b 
where a.name = b.name and a.val = b.val order by a.name
```

- 方法四：

```sql
select a.* from tb a inner join (select name , max(val) val from tb group by name) b on a.name = b.name and a.val = b.val order by
```

- 方法五：

```sql
select a.* from tb a where 1 > (select count(*) from tb where name = a.name and val > a.val ) order by a.name
```

- 方法六:

```sql
select * from (select * from tb ORDER BY val desc) temp GROUP BY name ORDER BY val desc;
```

以上六种方法运行的结果均为如下所示：

| name | val | memo |
| ---- | --- | ---- |
| a    | 3   | a3   |
| b    | 5   | b5   |

小编推荐使用第一、第三、第四钟方法,结果显示第1,3,4种方法效率相同，第2，5种方法效率差些。

### 按name分组取val最小的值所在行的数据

- 方法一：

```sql
select a.* from tb a where val = (select min(val) from tb where name = a.name) order by a.name
```

- 方法二：

```sql
select a.* from tb a where not exists(select 1 from tb where name = a.name and val < a.val)
```

- 方法三：

```sql
select a.* from tb a,(select name,min(val) val from tb group by name) b where a.name = b.name and a.val = b.val order by a.name
```

- 方法四：

```sql
select a.* from tb a inner join (select name , min(val) val from tb group by name) b on a.name = b.name and a.val = b.val order by a.name
```

- 方法五：

```sql
select a.* from tb a where 1 > (select count(*) from tb where name = a.name and val < a.val) order by a.name
```

以上五种方法运行的结果均为如下所示：

| name | val | memo |
| ---- | --- | ---- |
| a    | 1   | a1   |
| b    | 1   | b1   |

### 按name分组取第一次出现的行所在的数据

- 方法一：

```sql
select a.* from tb a where val = (select top 1 val from tb where name = a.name) order by a.name
//这个是sql server的
//mysql应该是
select a.* from tb a where val = (select val from tb where name = a.name limit 1) order by a.name
```

结果如下：

| name | val | memo |
| ---- | --- | ---- |
| a    | 2   | a2   |
| b    | 1   | b1   |
