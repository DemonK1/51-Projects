### 负载均衡器

本项目是负责跟踪哪个服务器拥有最多的流量以及哪个服务器现在可用,然后把高负载的服务器的流量转移到其他服务器上以实现负载均衡和提高整体性能,从而减轻任何一个服务器的负担,确保应用程序的高效运行.

**关于负载均衡器:**
1. 动态分配请求：负载均衡器可以停止将新的请求发送到高负载的服务器，并将这些请求转发到其他负载较低的服务器上。这样可以避免高负载服务器过载，确保所有服务器的负载均衡。

2. 会话迁移：如果应用程序涉及用户会话管理，负载均衡器可以将用户会话从高负载的服务器迁移到其他负载较低的服务器上。这样用户可以无缝地继续他们的会话，而不会受到高负载服务器的影响。

**关于反向代理:**
1. 反向代理是一种服务器加速和负载均衡的技术，它可以隐藏真实的服务器地址，同时将客户端请求转发到后端的多个服务器上