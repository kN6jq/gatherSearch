# gatherSearch

集成fofa,hunter,shodan,0zone等搜索引擎的信息收集工具,带有自动保存xlsx文件功能

持续更新中...

# 使⽤方法

> 举例如下,日后有空还会新增其他搜索引擎,方便后续加入我的信息收集系统

## fofa


> tips: 如果查询文件内的为ip,为了快速,通过配置文件中的batchSize设置数量为一组拼接进行查询

```bash
gatherSearch fofa -d baidu.com # 搜索域名
gatherSearch fofa -i 1.1.1.1 # 搜索ip
gatherSearch fofa -f 1.txt # 从文件中读取域名或ip
gatherSearch fofa -f 1.txt -c # 从文件中读取自定义的搜索语句
```

## hunter

> tips: 如果查询文件内的为ip,为了快速,通过配置文件中的batchSize设置数量为一组拼接进行查询

> tips: hunter最大支持5个拼接,同时请注意api积分!!!

```bash
gatherSearch hunter -d baidu.com # 搜索域名
gatherSearch hunter -i 1.1.1.1 # 搜索ip
gatherSearch hunter -f 1.txt # 从文件中读取域名或ip
gatherSearch hunter -f 1.txt -c # 从文件中读取自定义的搜索语句
```

## shodan

```bash
gatherSearch shodan -d baidu.com # 搜索域名
```


## shodandb

```bash
gatherSearch shodandb -i 1.1.1.1 # 搜索ip
gatherSearch shodandb -i 1.1.1.1/24 # 搜索ip c段
gatherSearch shodandb -f 1.txt # 从文件中读取ip
gatherSearch shodandb -i 1.1.1.1/24 | httpx --title --status-code # 结合httpx进行扫描

```

## 0zone

```bash
gatherSearch 0zone -n 零零信安 # 搜索单位信息系统
gatherSearch 0zone -d baidu.com # 搜索域名(这里咨询了客服,需要提供域名进行查询)
```

# 注意事项

- 使用前请在config.yaml配置api
- fofa及hunter搜索时可设置是否获取全部或者指定的数量,请参考config.yaml


# 免责申明

本项目仅面向安全研究与学习，禁止任何非法用途

如您在使用本项目的过程中存在任何非法行为，您需自行承担相应后果

除非您已充分阅读、完全理解并接受本协议，否则，请您不要使用本项目