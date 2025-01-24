#!/bin/fish
# 构建脚本
set image_version $argv[1]

if test -z $image_version
    echo 请指定版本号！
    exit
end

if test ! -d $image_version
    echo 文件夹：$image_version 不存在！
    exit
end

# 构建同步服务端程序
set exe_path dockerfile/mc-sync-server
set config_path dockerfile/server-config.yaml
echo 正在构建服务端...
cd ../
echo 已切换工作目录到：(pwd)
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o $exe_path gitee.com/swsk33/mc-server-sync/cmd/server
upx -9 $exe_path
echo 复制配置文件...
cp ./config-template/server-config.yaml $config_path

# 构建Docker镜像
set image_name swsk33/minecraft-server-with-sync:$image_version
cd ./dockerfile
echo 已切换工作目录到：(pwd)
echo 正在构建Docker镜像...
docker build -f ./$image_version/Dockerfile -t $image_name --network host --build-arg ALL_PROXY="http://127.0.0.1:7500" .
echo 构建完成！

# 清理文件
echo 正在清理...
rm mc-sync-server
rm server-config.yaml
