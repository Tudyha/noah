# 远程主机管理
## 使用方法：

### 直接启动
修改配置文件configs/config.yaml
```yaml
server:
  port: 8080
admin:
  password: 123456
```
```shell
chmod +x noah && ./noah
```

### docker 启动
```shell
docker run -d -p 8080:8080 -e ADMIN_PASSWORD=123456 -v ./data:/app/data --name noah knodio/noah
```
国内拉镜像慢，可以使用以下方式
```shell
docker run -d -p 8080:8080 -e ADMIN_PASSWORD=123456 -v ./data:/app/data --name noah registry.cn-guangzhou.aliyuncs.com/knodio/noah
```
说明：
1. 8080为服务端端口
2. ADMIN_PASSWORD 为管理密码，默认为123456

### 访问地址：http://localhost:8080

## TODO
1. 流量加密（紧急）
2. 用户管理
3. 文件管理用户字段只有gid、uid，没有名称，需要修改
4. 文件管理图片、json等文件预览优化

## 截图
![](https://github.com/Tudyha/noah/blob/main/doc/client-list.png?raw=true)
![](https://github.com/Tudyha/noah/blob/main/doc/console.jpeg?raw=true)