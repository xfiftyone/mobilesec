#!/bin/bash
# created by x51.
time=$(date +%m-%d--%H:%M:%S)
echo "【👤】x51"
echo "【🕙】当前时间："$time


is_device_connected=$(adb devices -l|grep usb)
if [ ! "$1" ];then
    echo "【❌】未指定项目名称！使用方法：$0 yourProjectName"
    exit
else
    echo "【❕】当前项目名：$1"
    echo "【❕】创建远程项目文件夹..."
    remoteProjectPath=$(adb shell "mkdir /sdcard/$1")
    savePath="/sdcard/$1"
    echo "【❕】项目保存路径：/sdcard/$1"
fi
if [ "$is_device_connected" != "" ];then
    echo "【✅】设备已通过USB连接"
    echo $(adb devices)
    is_root=$(adb shell "su -c 'whoami'")
    if [ "$is_root" = "root" ];then
        echo "【✅】设备已root"
        pkgPath=$(adb shell "su -c 'find /data/data/com.tencent.mm/MicroMsg/*/appbrand/pkg/ | head -n1'")
        echo "【❕】小程序包路径："$pkgPath
        echo "【❕】清理pkgPath..."
        clearPkgs=$(adb shell "su -c 'rm $pkgPath*.wxapkg'")
        echo $clearPkgs
        echo "【✅】pkg文件夹已清空"
        echo "【❕】正在等待重新打开小程序"
        while :
        do
            newPkgs=$(adb shell "su -c 'find /data/data/com.tencent.mm/MicroMsg/*/appbrand/pkg/ -name *.wxapkg'")
            if [ "$newPkgs" != "" ];then
                echo "【❕】发现新增${#newPkgs[*]}个wxapkg文件"
                echo "【❕】耐心等待10s..."
                sleep 10 # 小程序下载
                for var in ${newPkgs[*]}
                do
                    echo $var
                done
                copyToProjectPath=$(adb shell "su -c 'cp $pkgPath* $savePath'")
                pullToLocal=$(adb pull $savePath ./)
                echo "【✅】导出完毕！"
                ls -la "$1"
                deleteRemoteProjectPath=$(adb shell "rm -rf $savePath")
                echo "【✅】远程项目文件夹已清理"
                break
            else
                continue
            fi
        done
    else
        echo "【❌】设备未root"
    fi
else
    echo "【❌】没有发现android设备！"
fi
