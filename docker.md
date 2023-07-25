# docker 基础记录

## docker 构建node_base 镜像 推送远程仓库

nodejs 应用打包docker创建精简1G多镜像

刚开始使用docker打包nodejs应用时，每次生成的镜像文件大小都超过1个G，不仅推送慢，而且占用系统空间。

主要原因
拉取最新node 镜像，发现node镜像有943MB。

```shell
docker pull node:latest
docker images
#node 943MB
```

解决办法
创建一个基础镜像node_base

```xshell
docker pull alpine:latest
# alpine精简Linux 5.6MB
```

所以主要目标就是精简这个基础镜像

Dockerfile如下

```dockerfile
FROM alpine:latest
# 安装nodejs npm
RUN apk add --no-cache --update nodejs npm
```

创建基础镜像node_base
构建镜像

```shell
# 构建镜像 包含用户名
sudo docker build -t meetdocker2020/node_base:latest .
# node_base 54M 54MB 相比于943MB。效果显而易见。

docker login
docker tag meetdocker2020/node_base:latest  meetdocker2020/node_base:v1.0 
docker push meetdocker2020/node_base:v1.0 
```

使用基础镜像打包应用
nodejs应用的Dockerfile

```dockerfile
# 这里会形成第一个层级（layers）
FROM node_base:latest

WORKDIR /app
COPY package.json /app/

# 这里会形成第二个层级（layers），如果package.json依赖文件不变，就会使用缓存的node_modules
RUN npm install

COPY . /app/

# 这里会形成第三个层级（layers）
RUN npm run build
CMD ["npm","run","start"]
```

设置需要忽略的文件夹（不复制打包），.dockerignore文件如下

```file
node_modules
.git
.vscode
.nuxt
.next
```

打包

```shell
docker build -t your_project_name:0.1 .
```

[参考来源](https://zhuanlan.zhihu.com/p/130738206 "nodejs 应用打包docker创建精简1G多镜像")

## docker 搭建 nginx,php,mysql 环境

```shell
mac安装mysql扩展需要支持 v8版本


# 启动mysql
docker run --name mysql8.0 --env=MYSQL_ROOT_PASSWORD=root  --volume=/Users/qin/dev/mysql/8.0/conf:/etc/mysql/conf.d --volume=/Users/qin/dev/mysql/8.0/logs:/logs --volume=/Users/qin/dev/mysql/8.0/data:/var/lib/mysql --volume=/var/lib/mysql  -p 3307:3306 -d mysql:8.0.33


# 启动php 5.6
docker run --name php5.6 -v /Users/qin/dev/docker/php5.6/porject:/var/www/html -v /Users/qin/dev/docker/php5.6/conf:/usr/local/etc/php -v /Users/qin/dev/docker/php5.6/logs:/phplogs -d --link mysql8.0 php:5.6-fpm


# 启动php8.1
docker run --name php8.1 -v /Users/qin/dev/docker/php5.6/porject:/var/www/html -v /Users/qin/dev/docker/php8.1/conf:/usr/local/etc/php -v /Users/qin/dev/docker/php8.1/logs:/phplogs -d --link mysql8.0 php:8.1.18-fpm

说明：
配置文件创建在 conf.d 下
/Users/qin/dev/docker/php8.1/conf/conf.d

# 启动Nginx
docker run --name nginx -p 81:80 -d -v /Users/qin/dev/docker/php5.6/porject:/var/www/html:ro -v /Users/qin/dev/docker/nginx/conf.d:/etc/nginx/conf.d:ro --link php5.6   --link php8.1  --link mysql8.0 nginx


# 启动Nginx 配置虚拟域名 本地host配置 test.pay.com 127.0.0.1  启动容器时添加 --add-host test.pay.com:172.20.10.4（本地ip）
docker run --name nginx_all --add-host test.pay.com:172.20.10.4 -p 8000:80 -p 8001:8001 -p 8002:8002  -p 8003:8003  -d -v /Users/qin/dev/docker/php5.6/porject:/var/www/html:ro -v /Users/qin/dev/docker/nginx/conf.d:/etc/nginx/conf.d:ro --link php5.6   --link php8.1  --link mysql8.0  --link mysql3306 nginx


# 启动Nginx 配置虚拟域名 本地host配置 test.pay.com:127.0.0.1  启动容器时添加
docker run --name nginx_all2 --add-host test.pay.com:127.0.0.1 -p 8000:80 -p 8001:8001 -p 8002:8002  -p 8003:8003  -d -v /Users/qin/dev/docker/php5.6/porject:/var/www/html:ro -v /Users/qin/dev/docker/nginx/conf.d:/etc/nginx/conf.d:ro --link php5.6   --link php8.1  --link mysql8.0  --link mysql3306 nginx


# 启动Nginx 映射本地host
docker run --name nginx_all_host --network=host -p 8000:80 -p 8001:8001 -p 8002:8002  -p 8003:8003  -d -v /Users/qin/dev/docker/php5.6/porject:/var/www/html:ro -v /Users/qin/dev/docker/nginx/conf.d:/etc/nginx/conf.d:ro --link php5.6   --link php8.1  --link mysql8.0  --link mysql3306 nginx

# 报错不允许 需要自己构建网络 docker: links are only supported for user-defined networks.

# 重启nginx 进入nginx容器重启
service nginx reload

# 连接名称使用 link名称 或者 host.docker.internal
# 例如：mysql -h mysql8.0 -u root -p
# 例如：mysql -h host.docker.internal -u root -p
```

相关配置demo

```shell
# nginx config

server {
    listen  80;
    server_name localhost;
    set $root /var/www/html/localhost;

    #access_log  /tmp/nginx/logs/localhost.net.access.log main;
    #error_log  /tmp/nginx/logs/localhost.net.error.log notice;

    location ~ .*.(gif|jpg|jpeg|bmp|png|ico|txt|js|css)$ {
        root $root;
    }

    location / {
        root $root;
        index  index.php index.html index.htm;
        if ( -f $request_filename) {
            break;
        }
        if (!-e $request_filename) {
            rewrite ^(.*)$ /index.php/$1 last;
            break;
        }
    }

    location ~ .php(.*)$ {
        root $root;
        set $script $uri;
        set $path_info "";
        if ($uri ~ "^(.+.php)(/.+)") {
            set $script $1;
            set $path_info $2;
        }
        fastcgi_pass     php5.6:9000;
        #fastcgi_index index.php;
        fastcgi_index    index.php?IF_REWRITE=1;
        fastcgi_param    PATH_INFO    $path_info;
        fastcgi_param    SCRIPT_FILENAME    $document_root$fastcgi_script_name;
        fastcgi_param    SCRIPT_NAME    $script;
        include          fastcgi_params;

    }

    location ~ /.ht {
        deny  all;
    }
    location ~ /.svn {
        deny  all;
    }
    location ~ /.git/ {
        deny  all;
    }
    location ~ /Logs/ {
         deny  all;
    }
    location ~ /Logs/.* {
    }
    location ~ /Logs/.* {
        deny  all;
    }
    location ~ .*.(sql|tar.gz|zip|gz|tar|rariso|rpm|apk|bak)$ {
        deny  all;
    }

}
```

安装php扩展

```shell
cd /usr/local/bin
docker-php-ext-install mysqli mysql pdo_mysql

# php.ini 配置放开 位置 在本机配置文件中创建 conf.d文件夹以及任意名称.ini文件
extension = pdo_mysql.so
cgi.force_redirect = 0
extension = mysqli.so
extension = mysql.so
```

连接host说明

```shell
For all platforms
Docker v 20.10及以上版本（自2020年12月14日）。

使用你的内部IP地址或连接到特殊的DNS名称host.docker.internal，这将解析到主机使用的内部IP地址。

On Linux, using the Docker命令，在你的Docker命令中添加--add-host=host.docker.internal:host-gateway以启用该功能。

要在下列情况下启用该功能Docker Compose在Linux上，在容器定义中添加以下几行。

extra_hosts:
    - "host.docker.internal:host-gateway"
For older macOS and Windows versions of Docker
Docker v 18.03及以上版本（自2018年3月21日）。

使用你的内部IP地址或连接到特殊的DNS名称host.docker.internal，这将解析到主机使用的内部IP地址。

Linux support pending https://github.com/docker/for-linux/issues/264

For older macOS versions of Docker
Docker for Mac v 17.12 to v 18.02

同上，但用docker.for.mac.host.internal代替。

Docker for Mac v 17.06 to v 17.11

同上，但用docker.for.mac.localhost代替。

Docker for Mac 17.05及以下版本

要从docker容器中访问主机，你必须在网络接口上附加一个IP别名。你可以绑定任何你想要的IP，只要确保你不把它用在其他地方。

sudo ifconfig lo0 alias 123.123.123.123/24

然后确保你的服务器正在监听上面提到的IP或者0.0.0.0。如果它监听的是localhost 127.0.0.1，它将不接受连接。

然后，只要将你的docker容器指向这个IP，你就可以访问主机了。

为了测试，你可以在容器内运行类似curl -X GET 123.123.123.123:3000的东西。

该别名将在每次重启时重置，所以如果有必要，可以创建一个启动脚本。

解决方案和更多文件在这里。https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds
```
