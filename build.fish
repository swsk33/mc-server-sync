#!/bin/fish

set app_version $argv[1]

if test -z $app_version
    echo 请指定版本号！
    exit
end

# 基本构建路径
set base_output target

# 文件夹存在则清理上次构建
if test -d $base_output
    echo 清理上次构建...
    rm -r $base_output
end

mkdir -p $base_output

# 文件名称
set base_exe_name mc-sync-

# 构建一个发行版
# 参数1：构建目标操作系统，可用值：windows, linux
# 参数2：构建目标架构，可用值：386, amd64
# 参数3：构建客户端还是服务端，可用值：client, server
function build_release
    # 处理参数
    set build_os $argv[1]
    set build_arch $argv[2]
    set build_target $argv[3]
    # 根据构建目的系统，确认输出文件的格式和打包格式
    set output_exe_name {$base_exe_name}{$build_target}
    set archive_format tar.xz
    if test $build_os = windows
        set output_exe_name {$output_exe_name}.exe
        set archive_format 7z
    end
    # 构建临时文件夹
    set temp_dir $base_output/$build_os-$build_arch-$build_target
    # 构建程序
    echo 正在构建代码...
    set output_exe_path $temp_dir/$output_exe_name
    GOOS=$build_os GOARCH=$build_arch go build -ldflags "-w -s" -o $output_exe_path gitee.com/swsk33/mc-server-sync/cmd/$build_target
    # 压缩可执行文件
    echo 正在压缩...
    upx -9 $output_exe_path
    # 复制配置文件
    echo 复制配置文件...
    cp ./config-template/{$build_target}-config.yaml $temp_dir/
    # 执行打包
    echo 正在打包...
    if test $build_arch = 386
        set build_arch i386
    end
    set output_archive_path $base_output/{$base_exe_name}{$build_target}-$build_os-$build_arch-$app_version
    if test $archive_format = 7z
        7z a -t7z -mx9 $output_archive_path.7z ./$temp_dir/*
    else
        tar -cJvf $output_archive_path.tar.xz -C $temp_dir .
    end
    # 清理
    echo 构建完成！
    echo 正在清理...
    rm -r $temp_dir
end

# 构建参数列表
set os_list windows linux
set arch_list 386 amd64
set target_list client server

# 执行构建
for os in $os_list
    for arch in $arch_list
        for target in $target_list
            build_release $os $arch $target
        end
    end
end
