## 基于go的日志采集工具

a log collector which in support of sending data to kafka or other mq

how to use ?

```go

  cd $path/logshot/ 目录下
  执行 `make build` , 如果生成$path下生成/builds/logshot二进制文件,则表示构建成功


  进入 $path/builds

  执行命令 ./logshot --config=$configPath/conf.ini

```

命令相关参数说明:

--readall 监控整个文件(默认为false)
--sender=xxx 指定自定义的采集器(默认为default, eg: --sender=kafka,则采用kafka采集器进行日志采集)
--hb 是否开启心跳检测
--profile 生成性能快照 (cpu 和 内存)
--check 检测配置文件合法性