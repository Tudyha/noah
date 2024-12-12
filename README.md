# 远程主机管理

## 使用方法：

### 直接启动

修改配置文件 configs/config.yaml

```yaml
server:
  host: 127.0.0.1
  port: 8080
```

```shell
chmod +x noah && ./noah
```

### docker 启动

```shell
docker run -d -p 8080:8080 -e HOST=127.0.0.1 -v ./data:/app/data --name noah knodio/noah
```

国内拉镜像慢，可以使用以下方式

```shell
docker run -d -p 8080:8080 -e HOST=127.0.0.1 -v ./data:/app/data --name noah registry.cn-guangzhou.aliyuncs.com/knodio/noah
```

说明：

1. 8080 为服务端端口
2. /app/data 数据文件夹, 最好挂载到宿主机，避免数据丢失

### 访问地址：http://localhost:8080 (默认用户名密码：admin/123456)

## TODO

1. 流量加密（紧急）
2. 用户管理
3. 文件管理用户字段只有 gid、uid，没有名称，需要修改
4. 文件管理图片、json 等文件预览优化
5. ip 归属地

## 截图

![](https://github.com/Tudyha/noah/blob/main/doc/client-list.png?raw=true)
![](https://github.com/Tudyha/noah/blob/main/doc/console.jpeg?raw=true)
