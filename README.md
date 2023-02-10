# 家庭相册整理
扫描指定文件路径下的所有图片和视频进行去重,并将剩余有效信息记录至`ipicture.db`文件

## 功能与规则
### 去重
对于各种原因产生的重复图片或视频只保留唯一
### 归类
保留原有的分类文件夹,不对原文件或文件夹做修改,如果设置`delDuplicate`为`true`,则会从中删除重复的部分

### TODO
- [ ] 分类：将有效文件移动到指定目录下，并按一定规则分类

## 使用
需要先安装 [exiftool](https://www.sno.phy.queensu.ca/~phil/exiftool/)

程序启动后会扫描指定目录下的所有文件,并在当前目录下生成`ipicture.db`文件,其中记录了去重后的可用的图片和视频文件的有效物理路径及相关信息.
程序最多只删除重复的文件,不对现有其他数据做任何修改。

### 编译
`./build.sh xxx` windows / arm64 / linux
### 运行
`./bin/arm64/ipicture -path=path to photos root directory`
