# gatherSearch

集成fofa,hunter,shodan,0zone搜索引擎的信息收集工具,带有自动保存xlsx文件功能


## 使⽤方法

> 举例如下,日后有空还会新增其他搜索引擎,方便后续加入我的信息收集系统

### fofa

```bash
gatherSearch -p fofa -d baidu.com # 搜索域名
gatherSearch -p fofa -i 1.1.1.1 # 搜索ip
gatherSearch -p fofa -f 1.txt # 从文件中读取域名或ip
gatherSearch -p fofa -f 1.txt -c # 从文件中读取自定义的搜索语句
```

### hunter

```bash
gatherSearch -p hunter -d baidu.com # 搜索域名
gatherSearch -p hunter -i 1.1.1.1 # 搜索ip
gatherSearch -p hunter -f 1.txt # 从文件中读取域名或ip
gatherSearch -p hunter -f 1.txt -c # 从文件中读取自定义的搜索语句
```

### shodan

```bash
gatherSearch -p shodan -d baidu.com # 搜索域名
```


### 0zone

```bash
gatherSearch -p 0zone -n 零零信安 # 搜索单位信息系统
gatherSearch -p 0zone -d baidu.com # 搜索域名(这里咨询了客服,需要提供域名进行查询)
```

## 注意事项

- 使用前请在config.yaml配置api
- 当超过搜索引擎设置的最大结果时,使用配置文件最大搜索数量