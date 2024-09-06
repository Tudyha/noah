# 前端文件打包
cd ../web || exit
npm run build:prod

# 构建镜像
cd ..
docker build -t knodio/noah .