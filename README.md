# 远程主机管理
## 使用方法：

### docker 启动
```shell
docker run -d -p 80:9527 -p 1234:8080 -e ADMIN_PASSWORD=123456 --name noah knodio/noah
```
说明：
1. 9527为web管理访问端口
2. 8080为服务端端口，用于客户端与服务器通信
3. ADMIN_PASSWORD 为管理密码，默认为123456

### 源码部署
#### 1. 启动服务端
修改配置文件configs/config.yaml
```yaml
server:
  port: 8080
admin:
  password: 123456
```

启动服务端
```shell
go run cmd/noah/main.go
```

#### 2.部署前端
以下操作在web目录下执行

安装依赖
```shell
npm install
```
修改配置文件.env.staging
```text
# 服务端地址
VUE_APP_BASE_API = 'http://127.0.0.1:8080'

VUE_APP_WS_ADDR = 'ws://127.0.0.1:8080'
```
打包
```shell
npm run build:stage
```
