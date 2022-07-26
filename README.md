# 介绍

dt-runner是一个集中式的CI触发平台。在传统的gitlab ci中，研发人员使用.gitlab-ci.yaml文件控制项目的ci流程，但是，这种流程存在权限的管控漏洞，研发人员可以在CI脚本中写入任意内容，从而导致 CI流程的不可控。dt-runner是一个集中式的CI触发工具，其禁用了项目的CI文件，改有集中控制，加强项目的可控性。

* 监听gitlab所有事件
* 创建CRD
* 运行k8s pod
* 采集运行结果
* 发送邮件

## 使用手册

* 编译

* 运行

