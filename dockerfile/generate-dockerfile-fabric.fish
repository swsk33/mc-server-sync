#!/bin/fish

# 生成构建镜像的Dockerfile文件
# 脚本参数：
# 1. Minecraft版本号
# 2. Fabric Loader版本号
# 3. Fabric Installer版本号
#
# Minecraft版本可参考：https://mcversions.net/
# Fabric版本可参考：https://fabricmc.net/use/server/

# 处理参数
set minecraft_version $argv[1]
set fabric_loader_version $argv[2]
set fabric_installer_version $argv[3]

function printHelp
    echo
    echo 用法：
    echo "./generate-dockerfile.fish <Minecraft版本号> <Fabric Loader版本号> <Fabric Installer版本号>"
end

if test -z $minecraft_version
    echo 请指定Minecraft版本号！
    printHelp
    exit
end

if test -z $fabric_loader_version
    echo 请指定Fabric Loader版本号！
    printHelp
    exit
end

if test -z $fabric_installer_version
    echo 请指定Fabric Installer版本号！
    printHelp
    exit
end

# 准备生成文件
echo 正在生成Dockerfile文件...

set file_path "./Dockerfile"
set file_content "FROM bellsoft/liberica-runtime-container:jre-21-slim-glibc
WORKDIR /minecraft/data
# 加入服务器核心
ADD https://meta.fabricmc.net/v2/versions/loader/$minecraft_version/$fabric_loader_version/$fabric_installer_version/server/jar /minecraft/fabric-server.jar
# 加入同步服务端及其配置文件
ADD mc-sync-server /minecraft/
ADD server-config.yaml /minecraft/data/
# 加入启动脚本和其它数据文件
ADD start.sh /
ADD timezone.tar /usr/share/zoneinfo/
# 初始化
RUN chmod +x /start.sh \\
	&& chmod +x /minecraft/mc-sync-server \\
	&& java -jar /minecraft/fabric-server.jar --initSettings \\
	&& rm -r /minecraft/data/logs/
# 端口
EXPOSE 25565
EXPOSE 25566
# 环境变量
ENV LANG=C.UTF-8
ENV JVM_MIN=1G
ENV JVM_MAX=2G
ENV TZ=\"Asia/Shanghai\"
ENV EULA=false
# 数据卷
VOLUME [\"/minecraft/data\"]
CMD [\"/start.sh\"]"

echo $file_content >$file_path
echo 已生成Dockerfile到$file_path，请执行build-image.fish脚本构建镜像！
