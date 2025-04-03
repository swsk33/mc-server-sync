# Docker镜像构建指南

该目录中有下列脚本用于生成Dockerfile文件并完成构建工作：

- `generate-dockerfile-fabric.fish` 用于生成Dockerfile文件，可指定不同Minecraft版本、Fabric版本以生成不同的构建文件
- `build-image.fish` 一键构建

首先，使用`generate-dockerfile-fabric.fish`脚本生成Dockerfile文件：

```bash
./generate-dockerfile-fabric.fish "Minecraft版本号" "Fabric Loader版本号" "Fabric Installer版本号"
```

这会在当前目录下生成一个Dockerfile文件，然后构建即可：

```bash
./build-image.fish "Minecraft版本号" "模组加载器名称" "模组加载器版本号" "模组同步服务器版本号"
```