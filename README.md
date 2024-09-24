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
docker run -d -p 80:9527 -e ADMIN_PASSWORD=123456 --name noah knodio/noah
```
国内拉镜像慢，可以使用以下方式
```shell
docker run -d -p 80:9527 -e ADMIN_PASSWORD=123456 --name noah registry.cn-guangzhou.aliyuncs.com/knodio/noah
```
说明：
1. 9527为web管理访问端口
2. ADMIN_PASSWORD 为管理密码，默认为123456

## TODO
1. 流量加密（紧急）
2. 用户管理
3. 文件管理用户字段只有gid、uid，没有名称，需要修改
4. 文件管理图片、json等文件预览优化

