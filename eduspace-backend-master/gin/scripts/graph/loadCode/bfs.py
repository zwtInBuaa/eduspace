import copy
import queue
import sys
import json

code = '''# color表示顶点状态，gray表示未发现，lightblue表示已处理，steelblue表示已加入队列
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

inputDefault = ["0"]

nodes = [
    {"id": 0},
    {"id": 1},
    {"id": 2},
    {"id": 3},
    {"id": 4},
    {"id": 5},
    {"id": 6},
    {"id": 7},
    {"id": 8},
    {"id": 9}
]
links = [
    {"source": 0, "target": 1, "weight": ''},
    {"source": 0, "target": 5, "weight": ''},
    {"source": 1, "target": 2, "weight": ''},
    {"source": 1, "target": 7, "weight": ''},
    {"source": 2, "target": 3, "weight": ''},
    {"source": 2, "target": 9, "weight": ''},
    {"source": 3, "target": 4, "weight": ''},
    {"source": 3, "target": 8, "weight": ''},
    {"source": 4, "target": 5, "weight": ''},
    {"source": 4, "target": 6, "weight": ''},
    {"source": 6, "target": 7, "weight": ''},
    {"source": 6, "target": 8, "weight": ''},
    {"source": 7, "target": 9, "weight": ''},
    {"source": 8, "target": 9, "weight": ''}
]

v_nums = len(nodes)
e_nums = len(links)
V = [v['id'] for v in nodes]
E = {}
for i in range(v_nums):
    E[V[i]] = []
for i in range(e_nums):
    E[links[i]["source"]].append([links[i]["target"], links[i]["weight"]])
    E[links[i]["target"]].append([links[i]["source"], links[i]["weight"]])

color = ['gray' for _ in range(v_nums)]
pred = [-1 for _ in range(v_nums)]
int_max = sys.maxsize
dist = [int_max for _ in range(v_nums)]

res = {
    "code": code,
    "inputDefault": inputDefault,
    "nodeColors": [],
    "linkColors": [],
    "nodeLabelColors": [],
    "linkTextColors": [],
    "highlightSteps": [],
    "links": links,
    "nodes": nodes
}

nodeColors = ['gray' for _ in range(len(nodes))]
linkColors = ['#d3d3d3' for _ in range(len(links))]
linkTextColors = []
nodeLabelColors = ['' for _ in range(len(links))]
highlightSteps = []


def is_equal_to_pre_step() -> bool:
    if len(res["nodeColors"]) > 0 and nodeColors == res["nodeColors"][-1] and linkColors == res["linkColors"][-1] \
            and linkTextColors == res["linkTextColors"][-1] and nodeLabelColors == res["nodeLabelColors"][-1] \
            and highlightSteps == res["highlightSteps"][-1]:
        return True
    return False


def add_to_res():
    if is_equal_to_pre_step():
        return
    res["nodeColors"].append(copy.deepcopy(nodeColors))
    res["linkColors"].append(copy.deepcopy(linkColors))
    res["nodeLabelColors"].append(copy.deepcopy(nodeLabelColors))
    res["linkTextColors"].append(copy.deepcopy(linkTextColors))
    res["highlightSteps"].append(copy.deepcopy(highlightSteps))
    return


def find_link_label(a, b):
    length = len(links)
    for i in range(length):
        if links[i]['source'] == a and links[i]['target'] == b \
                or links[i]['source'] == b and links[i]['target'] == a:
            return i
    return -1


# 定义一个函数来实现图的广度优先遍历
def bfs(start):
    global highlightSteps
    highlightSteps = [3]
    add_to_res()
    Q = queue.Queue()
    highlightSteps = [4, 5, 6, 7, 8]
    add_to_res()
    color[start] = 'lightblue'
    dist[start] = 0
    Q.put(start)
    nodeColors[start] = 'lightblue'
    highlightSteps = [9, 10, 11]
    add_to_res()
    # 遍历队列
    if Q.empty():
        highlightSteps = [13]
        add_to_res()
    while not Q.empty():
        highlightSteps = [13]
        add_to_res()
        # 取出队列首个元素
        u = Q.get()
        linked_v = [v[0] for v in E[u]]
        highlightSteps = [14]
        add_to_res()
        if len(linked_v) == 0:
            highlightSteps = [15]
            add_to_res()

        for v in linked_v:
            # 遍历当前选中节点的邻接节点
            e = find_link_label(u, v)
            tmpColor = linkColors[e]
            linkColors[e] = 'lightblue'
            highlightSteps = [15]
            add_to_res()
            # 如果目标节点已经加入队列，则恢复节点颜色
            if color[v] != 'gray':
                highlightSteps = [16]
                add_to_res()
                linkColors[e] = tmpColor
                add_to_res()
            # 否则，将目标节点加入队列
            if color[v] == 'gray':
                highlightSteps = [16]
                add_to_res()
                color[v] = 'lightblue'
                nodeColors[v] = 'lightblue'
                dist[v] = dist[u] + 1
                pred[v] = u
                e = find_link_label(u, v)
                linkColors[e] = 'steelblue'
                Q.put(v)
                highlightSteps = [17, 18, 19, 20]
                add_to_res()
        # 该节点已被访问过了
        color[u] = 'steelblue'
        nodeColors[u] = 'steelblue'
        highlightSteps = [21]
        add_to_res()
    highlightSteps = [22]
    add_to_res()
    return


start = 0
highlightSteps = [24]
add_to_res()
bfs(start)
highlightSteps = []
add_to_res()

print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))