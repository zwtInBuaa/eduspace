# 所有图相关算法的伪代码保存在此


bfs_code = '''# color表示顶点状态，gray表示未发现，lightblue表示已处理，steelblue表示已加入队列
# pred表示顶点u由pred[u]发现
# dist表示顶点u距离源点start的距离
def bfs(start):
    Q = queue.Queue()  # 新建空队列
    for u in V:  # 初始化辅助数组
        color[u] = 'gray'
        pred[u] = NULL
        dist[u] = Infinity
    color[start] = 'lightblue'
    dist[start] = 0
    Q.put(start)
    # 广度优先搜索
    while not Q.empty():
        u = Q.get()
        for v in G.Adj[u]:  # u所临接的边
            if color[v] == 'gray':
                color[v] = 'lightblue'  # 该点已被发现
                dist[v] = dist[u] + 1
                pred[v] = u
                Q.put(v)
        color[u] = 'steelblue'  # 该点已被搜索
    return
# Graph[V,E]
bfs(___)'''

dfs_code = '''# color表示顶点状态，gray表示未发现，lightblue表示正在处理，steelblue表示已被处理
# pred表示顶点u由pred[u]发现
# d表示各顶点的发现时刻
# f表示各顶点的处理完成时刻
def dfs(start):
    time = 0
    dfsVisit(start)  # 先对起始节点进行DFS
    for u in V:
        if color[u] == 'gray':  # 保证遍历完全
            dfsVisit(u)
    return
def dfsVisit(u):
    color[u] = 'lightblue'
    time += 1
    d[u] = time  # 发现时刻
    for v in G.Adj[u]:  # u所临接的边
        if color[v] == 'gray':
            pred[v] = u
            disVisit(v)
    color[u] = 'steelblue'
    time += 1
    f[u] = time  # 完成时刻
    return
# Graph[V,E]
dfs(___)'''

dijkstra_code = '''# color表示顶点状态，lightblue表示正在被处理，steelblue表示最短路已被确定，gray表示未确定
# dist记录距离上届
# pred表示前驱节点
def dijkstra(start):
    dist[start] = 0  # 源点到自身距离为0
    for i in len(V):  # 依次计算源点到各顶点的最短路
        minDist = Infinity
        rec = 0
        for j in len(V):  # 选择距源点最近的gray色顶点
            if color[j] != 'steelblue' and dist[j] < minDist:
                minDist = dist[j]
                rec = j
        color[rec] = 'lightblue'  # 该节点正在被处理
        for u in G.Adj[rec]:  # 对rec出发的边进行松弛
            if dist[rec] + weight(rec, u) < dist[u]:
                dist[u] = dist[rec] + weight(rec, u)  # 松弛操作
                pred[u] = rec  # 记录前驱顶点
        color[rec] = 'steelblue'  # 该节点已被处理完成
    return
# Graph[V,E]
dijkstra(___)'''

kruskal_code = '''# 采用并查集判断和维护所选边的顶点是否在一棵子树
# find(x)寻找顶点x的树根结点
# union(x,y)将2顶点所在树结构合并，合并成功返回True
def kruskal():
    parent = [v for v in V] # 所有顶点的树根节点
    rank = [0] * len(V) # 并查集辅助数组
    # e的结构{'source':'','target':'','weight':''}
    edges = [e for e in E]
    edges.sort(key=lambda x: x['weight']) # 边排序
    mst = [] # 存放着选取的所有边
    weights = 0 # 最小生成树总权值
    for e in edges: # 每次加入不形成环路的最小权值边
        if union(e['source'], e['target']):
            mst.append(e)
            weights += e['weight']
    return weights, mst
# Graph[V,E]
kruskal()'''

prim_code = '''# color表示顶点状态，lightblue表示正在被处理，steelblue表示已覆盖，gray表示未覆盖
# dist记录横跨(Va,V - Va)边的权重
# pred表示前驱节点
def prim(start):
    dist[start] = 0  # 从start节点开始
    for i in len(V):  # 循环次数
        minDist = Infinity
        rec = 0
        for j in len(V):  # 记录新增的安全边
            if color[j] != 'steelblue' and dist[j] < minDist:
                minDist = dist[j]
                rec = j
        color[rec] = 'lightblue'  # 该节点正在被处理
        for u in G.Adj[rec]:  # 更新dist数组，即安全边集合
            if weight(rec, u) < dist[u]:
                dist[u] = weight(rec, u)
                pred[u] = rec
        color[rec] = 'steelblue'  # 该节点处理完成
    return
# Graph[V,E]
prim(___)'''