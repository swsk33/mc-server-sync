FROM bellsoft/liberica-runtime-container:jre-21-slim-glibc
WORKDIR /minecraft/data
# 加入服务器核心
ADD https://meta.fabricmc.net/v2/versions/loader/1.21.4/0.16.10/1.0.1/server/jar /minecraft/fabric-server.jar
# 加入同步服务端及其配置文件
ADD mc-sync-server /minecraft/
ADD server-config.yaml /minecraft/data/
# 加入启动脚本和其它数据文件
ADD start.sh /
ADD timezone.tar /usr/share/zoneinfo/
# 初始化
RUN chmod +x /start.sh \
	&& chmod +x /minecraft/mc-sync-server \
	&& java -jar /minecraft/fabric-server.jar --initSettings \
	&& rm -r /minecraft/data/logs/
# 端口
EXPOSE 25565
EXPOSE 25566
# 环境变量
ENV LANG=C.UTF-8
ENV JVM_MIN=1G
ENV JVM_MAX=2G
ENV TZ="Asia/Shanghai"
ENV EULA=false
# 数据卷
VOLUME ["/minecraft/data"]
CMD ["/start.sh"]