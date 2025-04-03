#!/bin/fish

# Docker镜像构建脚本

# 处理参数
set minecraft_version $argv[1]
set mod_loader_name $argv[2]
set mod_loader_version $argv[3]
set sync_server_version $argv[4]

function printHelp
    echo
    echo 用法：
    echo "./build-image.fish <Minecraft版本号> <模组加载器名称> <模组加载器版本号> <模组同步服务器版本号>"
end

if test -z $minecraft_version
    echo 请指定Minecraft版本号！
    printHelp
    exit
end

if test -z $mod_loader_name
    echo 请指定模组加载器名称！
    printHelp
    exit
end

if test -z $mod_loader_version
    echo 请指定模组加载器版本号！
    printHelp
    exit
end

if test -z $sync_server_version
    echo 请指定模组同步服务器版本号！
    printHelp
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
set image_name swsk33/minecraft-server-with-sync:$minecraft_version-$mod_loader_name-$mod_loader_version-sync-$sync_server_version
cd ./dockerfile
echo 已切换工作目录到：(pwd)
echo 正在构建Docker镜像...
docker build -f ./Dockerfile -t $image_name --network host --build-arg ALL_PROXY="http://127.0.0.1:7500" .
echo 构建完成！

# 清理文件
echo 正在清理...
rm mc-sync-server
rm server-config.yaml
