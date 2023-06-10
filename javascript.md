# JS Date对象时间戳和时区一些注意点

时区，时间戳，日期

## 什么是时区？

整个地球分为二十四时区，每个时区都有自己的本地时间，本地时间 = UTC + 时区差，时区差东为正，西为负。因此，把东八时区（北京时间 ）记为 UTC+8。

## 时间戳

时间戳是一个绝对值，代表格林威治时间（1970年01月01日00时00分00秒）到现在的总毫秒数，跟时区没有任何关系，你在地区任何一个地方都是这个值

### 解决方案

#### 第一种：约定好都传时间戳

最简单的方式肯定就是约定好都传时间戳，时间戳是一个绝对数，所有的逻辑都在后台的服务器上处理，服务器肯定都在一个时区，就不会出现时区上的问题

#### 第二种：原生JS

  getTimezoneOffset方法，表示UTC与本地时区之间的差值，单位为 分钟。
如果本地时区后于协调世界时，则该差值为正值，如果先于协调世界时则为负值。例如北京时区为 UTC+8，将会返回 -480。对于同一个时区，夏令时将会改变这个值。

```js
var timeZoneDiff = 8; // 目标时区(北京时区)，GMT+8
var offsetTime = new Date().getTimezoneOffset(); // 本地时间和UTC0时的时间差，单位为分钟
var timeStamp = new Date().getTime(); // 本地时间距 1970 年 1 月 1 日0时之间的毫秒数
var targetDate = new Date(timeStamp + offsetTime * 60 * 1000 + timeZoneDiff * 60 * 60 * 1000); // 本地时间戳 + 偏移时间戳 + 目标时间戳 =》 目标时间
console.log("现在北京时间是：" + targetDate);
```

#### 第三种：引入Timezone插件，以dayjs为例

```js
import * as dayjs from 'dayjs';
import * as utc from 'dayjs/plugin/utc';
import * as timezone from 'dayjs/plugin/timezone';
dayjs.extend(utc);
dayjs.extend(timezone);

// 无论在哪个地方，都把本地的时间转成中国时间传给后端
dayjs(new Date()).tz("Asia/Shanghai").format('YYYY-MM-DD HH:mm:ss')
```

不同地区标识符是不一样的，国内的是 Asia/Shanghai，其他的值可以参考 [时区](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

原文链接：<https://juejin.cn/post/7204085076504854583>
