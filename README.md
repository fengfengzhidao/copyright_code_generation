
# 软著代码生成

## 主要功能
1. 读取代码目录，自动复制代码写入到word文档
2. 自动读取前后30页的代码，不用再手动删除

## 如何使用
```text
  -f string
        需要识别的文件，默认是.go文件，如果需要多个 以;号分隔
  -n string
        页眉的名称
  -r string
        程序目录
```
示例
```text
go run main.go -r D:\IT\fengfeng\test_go_pro -f .go -n "我的go系统 v1.0"
```