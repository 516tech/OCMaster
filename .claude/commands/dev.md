# 原则

原则：有任何疑问，第一时间向我提问
要求：所有文档行数不超过200，表格不超过3个，同级别标题不超过5个
目标：参考prd.md，如果方案有变更，重写prd.md
搜索：使用browser harness 
流程：分step进行，step 1计划，step 2实现，step 3评估，step 4总结
编译：github的action

# step 1

设计：需要拆分成多个子任务，任务通过md的checkbox标注是否开发完成
修复：看下现有bug.md，设计修复内容
总结：生成architecture.md文件，如果方案有变更，重写architecture.md

# step 2

选型：C端使用nodejs+electron实现windows/linux/mac跨平台，S端使用golang+gorm+zerolog，前端使用vue，用DDD实现对象注入
架构：C端是单文件，S端用docker compose实现bin+mysql+nginx
编译：本地使用docker+colima
测试：所有模块增加单元测试，整体增加集成测试，避免使用配置文件启动
配置：都提供默认配置，提供用户可修改配置的页面，保存到本地
总结：增量更新evaluate.md文件，添加本轮迭代生成内容评估项

# step 3

评估：根据evaluate.md，逐项评估目标是否实现
测试：运行单元测试和集成测试，能解决的bug直接解决
总结：无法解决的问题，生成bug.md

# step 4

优化：把所有涉及文档prd.md，architecture.md，evaluate.md，bug.md进行优化
