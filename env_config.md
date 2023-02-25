[TOC]

## 安装docker

- 安装docker

  - 安装

    ```shell
    # 1. cd进入到/etc/yum.repos.d/中
    cd /etc/yum.repos.d/
    
    # 2. 下载docker的仓库文件
    wget https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
    
    # 3. 安装docker-ce
    yum install -y docker-ce
    ```

  - 配置加速器

    ```shell
    # 1. 进入/etc/docker/中
    cd /etc/docker/
    
    # 2. 创建daemon.json文件
    touch daemon.json
    
    # 3. 在daemon.json中配置一下内容
    {
        "registry-mirrors": [
            "https://1nj0zren.mirror.aliyuncs.com",
            "https://docker.mirrors.ustc.edu.cn",
            "http://f1361db2.m.daocloud.io",
            "https://registry.docker-cn.com"
        ]
    }
    
    # 4. 加载daemon配置文件
    systemctl daemon-reload 
    
    # 5. 重启docker
    systemctl restart docker
    ```

  - 设置开机自启

    ```shell
    systemctl enable docker
    ```

- 安装docker-compose

  - 安装

    ```shell
    curl -L https://get.daocloud.io/docker/compose/releases/download/1.25.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
    ```

  - 设置权限

    ```shell
    sudo chmod +x /usr/local/bin/docker-compose
    ```

  - 检验

    ```shell
    docker-compose -v
    ```

## 安装mysql

- 下载镜像

  ```she
  $ docker pull mysql:5.7
  ```

- 通过镜像启动

  ```she
  $ docker run -p 3306:3306 --name mysql -v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql -e MTSQL_ROOT_PASSWORD=123456 -d mysql:5.7
  
  -p 3306:3306将容器的3306端口映射到主机的3306端口
  
  -v $PWD/conf:/etc/mysql/conf.d 将主机当前目录的 conf/my.cnf 挂载到容器的 etc/mysql/my.cnf
  
  -v $PWD/logs:/logs 将主机当前目录下的 logs 目录挂载到容器的 /var/lib/mysql
  
  -e MTSQL_ROOT_PASSWORD=123456 -d mysql:5.7 初始化 root 密码
  
  运行：docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
  ```

- 进入容器配置

  ```she
  进入容器
  $ docker exec -it XXXXXXXXX(mysql CONTAINER ID) /bin/bash
  
  进入mysql
  $ mysql -uroot -p123456
  
  进入用户并授权
  $ GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'root' WITH GRANT OPTION;
  $ GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' IDENTIFIED BY 'root' WITH GRANT OPTION;
  $ GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' IDENTIFIED BY 'root' WITH GRANT OPTION;
  $ FLUSH PRIVILEGES; 
  
  $ exit
  $ exit
  ```


## 配置go环境

- go version go1.15.3 windows/amd64

- windows

  ```she
  设置代理
  $ go env -w GO111MODULE=on
  $ go env -w GOPROXY=https://goproxy.cn,direct
  ```

- centos

  ```shell
  1.下载
  $ wget https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz
  
  2.解压
  $ tar -xvf go1.15.3.linux-amd64.tar.gz
  
  3.配置环境变量
  $ vi ~/.bashrc
  	export GOROOT=/root/go
  	export GOPATH=/root/projects/go
  	export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
  	
  $ source ~/.bashrc
  
  4.设置代理
  $ go env -w GOPROXY=https://goproxy.cn,direct
  $ go env -w GO111MODULE=on
  ```

## grpc

下载工具

- protoc-3.13.0-win64.zip

下载go依赖包

```she
$ go get github.com/golang/protobuf/protoc-gen-go

$ go get github.com/envoyproxy/protoc-gen-validate
```

proto文件生成

```she
$ protoc --go_out=plugins=grpc:. ./hello.proto

$ protoc --go_out=plugins=grpc:. --validate_out="lang=go:." ./hello.proto
```

## 安装yapi

- 初始化db,开启自定义配置

```shell
git clone https://github.com/Ryan-Miao/docker-yapi.git
cd docker-yapi
docker-compose up

用户名：admin@admin.com
密码：ymfe.org
```

- 打开 localhost:9090
- 默认部署路径为`/my-yapi`(需要修改docker-compose.yml才可以更改)
- 修改管理员邮箱 `ryan.miao@demo.com` (随意, 修改为自己的邮箱)
- 修改数据库地址为 `mongo` 或者修改为自己的mongo实例 (docker-compose配置的mongo服务名称叫mongo)
- 打开数据库认证
- 输入数据库用户名: `yapi`(mongo配置的用户名, 见mongo-conf/init-mongo.js)
- 输入密码: `yapi123456`(mongo配置的密码, 见mongo-conf/init-mongo.js)



- 安装完毕启动

  ```shell
  cd docker-yapi
  docker-compose up
  访问:3000
  ```

  

## 日志库zap

- 安装

	```shell
go get -u go.uber.org/zap

## 安装redis

- docker安装

  ```shell
  docker pull redis:latest
  ```

- 启动

  ```shell
  docker run -p 6379:6379 -d redis:latest redis-server
  ```

- 设置开机自启

  ```shell
  docker container update --restart=always CONTAINER ID
  ```

- go驱动

  ```shell
  https://github.com/go-redis/redis
  ```

## 安装Consul

- 安装

  ```shell
  docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
  
  docker container update --restart=always CONTAINER ID
  ```

- 访问

  ```she
  :8500
  ```

- 访问dns

  ```shell
  dig @192.168.10.105 -p 8600 consul.service.consul SRV
  ```


## 安装nacos

- 下载

  ```shell
  docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -p 9848:9848 -p 9849:9849 -d nacos/nacos-server:latest
  ```

- 访问

  http://192.168.10.105:8848/nacos/index.html

  用户名密码: nacos/nacos

## 安装elasticsearch

- 禁用防火墙

  ```shell
  systemctl stop firewalld.service
  systemctl disable firewalld.service
  systemctl status firewalld.service
  ```

- 安装elasticsearch

  ```shell
  #新建es的config配置文件
  mkdir -p /data/elasticsearch/config
  #新建es的data目录
  mkdir -p /data/elasticsearch/data
  #新建es的logs目录
  mkdir -p /data/elasticsearch/logs
  #给目录设置权限
  chmod 777 -R /data/elasticsearch
  
  #写入配置到elasticsearch.yml中
  echo "http.host: 0.0.0.0" >> /data/elasticsearch/config/elasticsearch.yml
  
  #安装es
  docker run --name elasticsearch -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e ES_JAVA_OPTS="-Xms128m -Xmx256m" \ 
  -v /data/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
  -v /data/elasticsearch/data:/usr/share/elasticsearch/data \
  -v /data/elasticsearch/logs:/usr/share/elasticsearch/logs \
  -v /data/elasticsearch/plugins:/usr/share/elasticsearch/plugins \
  -d elasticsearch:7.10.1
  ```

## 安装kibana

- 安装

  ```shell
  docker run -d --name kibana -e ELASTICSEARCH_HOSTS="http://192.168.10.105:9200" -p 5601:5601 kibana:7.10.1
  ```


## 安装jaeger

- 安装

  ```shell
  docker run \
    --rm \
    --name jaeger \
    -p6831:6831/udp \
    -p16686:16686 \
    jaegertracing/all-in-one:latest
  ```


## 安装kong

- 安装postgresql和migrations

  ```shell
  docker run -d --name kong-database \
  -p 5432:5432 \
  -e "POSTGRES_USER=kong" \
  -e "POSTGRES_DB=kong" \
  -e "POSTGRES_PASSWORD=kong" \
  -e "POSTGRES_DB=kong" postgres:12
  
  docker run --rm \
  -e "KONG_DATABASE=postgres" \
  -e "KONG_PG_HOST=192.168.10.105" \
  -e "KONG_PG_PASSWORD=kong" \
  -e "POSTGRES_USER=kong" \
  -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
  kong kong migrations bootstrap
  ```

- 安装kong

  ```shell
  sudo yum -y install https://download.konghq.com/gateway-2.x-centos-7/Packages/k/kong-2.1.0.el7.amd64.rpm
  ```

- 编辑kong配置

  ```shell
  systemctl stop filewalld.service
  systemctl restart docker
  
  cp /etc/kong/kong.conf.default /etc/kong/kong.conf
  vim /etc/kong/kong.conf
  #修改如下内容
  database = postgres
  pg_host = 192.168.10.105
  pg_port = 5432
  pg_timeout = 5000
  
  pg_user = kong
  pg_password = kong
  pg_database = kong
  
  dns_resolver = 127.0.0.1:8600 #配置的consul端口
  admin_listen = 0.0.0.0:8001 reuseport backlog=16384, 127.0.0.1:8444 http2 ssl reuseport backlog=16384
  proxy_listen = 0.0.0.0:8000 reuseport backlog=16384, 0.0.0.0:8443 http2 ssl reuseport backlog=16384
  ```

- kong启动

  ```shell
  kong start -c /etc/kong/kong.conf
  ```

- 安装konga

  ```shell
  docker run -d -p 1337:1337 --name konga pantsel/konga
  ```


### 安装jenkins

- 安装java

  ```shell
  yum install java-1.8.0-openjdk* -y
  ```

- 下载jenkins

  ```shell
  wget https://mxshop-files.oss-cn-hangzhou.aliyuncs.com/jenkins-2.284-1.1.noarch.rpm
  ```

- 上传安装包并安装

  ```shell
  rpm -ivh jenkins-2.284-1.1.noarch.rpm
  ```

- 修改jenkins配置

  ```shell
  vim /etc/sysconfig/jenkins
  
  JENKINS_USER="root"
  JENKINS_PORT="8088"
  ```

- 启动

  ```shell
  systemctl start jenkins
  ```

- 修改插件下载地址

  ```shell
  sed -i 's/https:\/\/updates.jenkins.io\/download/http:\/\/mirrors.tuna.tsinghua.edu.cn\/jenkins/g' /var/lib/jenkins/updates/default.json && sed -i 's/http:\/\/www.google.com/https:\/\/www.baidu.com/g' /var/lib/jenkins/updates/default.json
  ```

- 进入Manage Jenkins->Manage Plugins->Advanced->update Site修改为

  https://mirrors.tuna.tsinghua.edu.cn/jenkins/update-center.json

- restart

  http://192.168.10.105:8088/restart

  s
# mall
