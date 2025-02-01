本项目除了能够直接部署之外，还提供了Docker容器的方式部署，在Docker镜像中已经同时包含了下列内容：

- Minecraft Fabric模组服务器
- 模组同步服务端

只需创建并运行Docker容器，即可直接同时搭建起Minecraft模组服务器和模组同步服务器，更加便捷，推荐大家使用Docker容器化部署。

Docker镜像的`Tag`格式为：`Minecraft版本-模组加载器-模组加载器版本-sync-模组同步服务器版本`

不带`Tag`的就是最新版本镜像。

## 1，拉取镜像

执行下列命令：

```bash
docker pull swsk33/minecraft-server-with-sync
```

## 2，创建数据卷

拉取容器后，创建容器之前先创建个具名数据卷用于持久化世界存档文件和配置文件等等，也便于我们修改：

```bash
docker volume create minecraft-fabric-data
```

这样，就将全部服务端文件持久化到了宿主机。

## 3，创建容器

通过下列命令：

```bash
docker run -itd --name minecraft-fabric-server-sync \
	-p 25565:25565 \
	-p 25566:25566 \
	-v minecraft-fabric-data:/minecraft/data \
	-e EULA=true \
	swsk33/minecraft-server-with-sync
```

第一次需要等待世界创建，过几分钟服务端即启动，可以通过游戏连接。上述启动命令已同时映射了Minecraft服务器和模组同步服务端的端口到宿主机以供访问，请确保你的云服务器防火墙配置对这些端口放行。

需要指定`EULA`环境变量为`true`，表示你同意[Minecraft的最终用户许可协议⁠](https://aka.ms/MinecraftEULA)，否则将无法启动服务端。

此外，可以指定环境变量`JVM_MIN`和`JVM_MAX`来限制服务端所使用的最小启动堆内存和最大堆内存，也就是`java`命令的`-Xms`和`-Xmx`参数值：

```bash
# 限制堆内存范围为2G ~ 4G
docker run -itd --name minecraft-fabric-server-sync \
	-p 25565:25565 \
	-p 25566:25566 \
	-e JVM_MIN=2G \
	-e JVM_MAX=4G \
	-e EULA=true \
	-v minecraft-fabric-data:/minecraft/data \
	swsk33/minecraft-server-with-sync
```

默认情况下，最小启动堆内存限制为`1G`，最大为`2G`。

## 4，修改配置文件

通过上述配置数据卷之后，所有配置文件如下：

- Minecraft服务器配置：`/var/lib/docker/volumes/minecraft-fabric-data/_data/server.properties`
- `mc-sync-server`同步服务器配置：`/var/lib/docker/volumes/minecraft-fabric-data/_data/server-config.yaml`

使用文本编辑器编辑即可，编辑完成需要重启服务端，正确地停止并重启服务端请看下面说明。

## 5，加入模组

通过上述配置数据卷之后，模组文件夹位于：`/var/lib/docker/volumes/minecraft-fabric-data/_data/mods`，将模组放进去后，重启容器即可，正确地停止并重启服务端请看下面说明。

此外，同步服务端的仅客户端模组文件夹默认位于：`/var/lib/docker/volumes/minecraft-fabric-data/_data/client-mods`，可将想要同步给玩家的仅客户端模组放入。

## 6，正确地停止服务端

由于Minecraft服务端是交互式的，因此我们不能直接通过`docker stop`或者`docker restart`命令来停止和重启容器，否则可能导致服务端数据丢失。

先通过以下命令连接容器内服务端的交互式控制台：

```bash
docker attach minecraft-fabric-server-sync
```

此时就进入了Minecraft服务端的交互式控制台，可能你会发现连接上后并没有显示服务端先前在控制台输出的内容，看起来像卡死了一样的，但是事实上你已经成功连接上了服务端控制台，可以直接输入Minecraft服务端的指令，这时你试着输入`/help`试试，能够如你所愿。

输入`/save-all`命令可以保存世界，若要停止服务端输入`/stop`，这时服务端和容器都会停止，并回到宿主机终端。再次启动时使用`docker start`命令即可。修改配置文件之后就需要通过这种方式退出停止服务端并重新启动。

进入容器内Minecraft服务端控制台后，如果想退出容器但是仍然保持容器运行，依次按下`Ctrl + P`和`Ctrl + Q`组合键即可。

## 7，关于世界

根据上述配置之后，世界存档数据卷位于`/var/lib/docker/volumes/minecraft-fabric-data/_data/world`目录，如果想重新生成世界，可以先按照上述方式停止服务端，然后删除该目录（`world`文件夹），再重启容器即可。

如果想使用自己的世界，同样地停止服务端后，先删除`/var/lib/docker/volumes/minecraft-fabric-data/_data`中的`world`目录，然后把你的世界存档目录改为`world`上传至该目录下，再启动容器即可。

个别模组包含一些世界生成矿物、生物群系等，因此可以在停止容器后添加模组文件，并删除世界文件夹，再次启动容器，就会重新生成世界，并包含模组内容。

## 8，查看模组同步服务端日志

在容器中，模组同步服务端以守护进程模式运行，其日志文件位于容器内`/minecraft/data/mc-sync-server.log`，可使用下列命令查看：

```bash
docker exec -it minecraft-fabric-server-sync tail -f /minecraft/data/mc-sync-server.log
```

此时，就可以实时查看同步服务端日志，按下`Ctrl + C`组合键退出查看。