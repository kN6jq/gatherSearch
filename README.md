# gatherSearch

集成fofa,hunter,shodan三个搜索引擎的信息收集工具


## 使⽤方法

- main.exe -p shodan -d baidu.com # 指定使用shodan搜索引擎搜索baidu.com域名
- main.exe -p fofa -i 1.1.1.1 # 指定使用fofa搜索引擎搜索ip
- main.exe -p hunter -f ip.txt # 指定使用hunter搜索引擎搜索ip.txt文件中的数据,数据内容可以为ip或者域名
- main.exe -p fofa -f data.txt -c # data.txt文件中的数组为相对搜索引擎的语法,如fofa的语法,需要加上-c参数


## 注意事项

- 请在config.json中配置自己的api_key
- 目前shodan只支持domain搜索,其他搜索方式不支持
- 配置文件可设置最大搜索数量,默认为3000条,超过3000条的数据将设置数量为配置文件中的数量