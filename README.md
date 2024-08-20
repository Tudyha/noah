# 远程主机管理
## 使用方法：
### 1. 启动服务端
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

### 2.部署前端
以下操作在web目录下执行

安装依赖
```shell
npm install
```
修改配置文件.env.staging
```text
# 服务端地址
VUE_APP_BASE_API = 'http://127.0.0.1:8080'
```
打包
```shell
npm run build:stage
```
