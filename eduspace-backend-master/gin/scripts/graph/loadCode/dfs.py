import copy
import queue
import sys
import json

code = '''# color表示顶点状态，gray表示未发现，lightblue表示正在处理，steelblue表示已被处理
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
d = [-1 for _ in range(v_nums)]
f = [-1 for _ in range(v_nums)]

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
    res["linkTextColors"].append(copy.deepcopy(linkTextColors))
    res["nodeLabelColors"].append(copy.deepcopy(nodeLabelColors))
    res["highlightSteps"].append(copy.deepcopy(highlightSteps))
    return


def find_link_label(a, b):
    length = len(links)
    for i in range(length):
        if links[i]['source'] == a and links[i]['target'] == b \
                or links[i]['source'] == b and links[i]['target'] == a:
            return i
    return -1


# 定义一个函数来实现图的深度优先搜索
def dfs(start):
    global highlightSteps
    highlightSteps = [4]
    add_to_res()
    highlightSteps = [5]
    add_to_res()
    highlightSteps = [6]
    add_to_res()
    dfs_visit(start)

    if len(V) == 0:
        highlightSteps = [7]
        add_to_res()
    for u in V:
        highlightSteps = [7]
        add_to_res()

        if color[u] != 'gray':
            highlightSteps = [8]
            add_to_res()
        if color[u] == 'gray':
            highlightSteps = [8]
            add_to_res()
            highlightSteps = [9]
            add_to_res()
            dfs_visit(u)
    highlightSteps = [10]
    add_to_res()
    return


time = 0


def dfs_visit(u: int):
    global time, highlightSteps
    highlightSteps = [11]
    add_to_res()
    color[u] = 'lightblue'
    nodeColors[u] = 'lightblue'
    highlightSteps = [12, 13, 14]
    add_to_res()
    time += 1
    d[u] = time
    linked_v = [v[0] for v in E[u]]
    if len(linked_v) == 0:
        highlightSteps = [15]
        add_to_res()
    for v in linked_v:
        # 遍历当前节点的所有邻接节点
        l = find_link_label(u, v)
        tmpColor = linkColors[l]
        linkColors[l] = 'lightblue'
        highlightSteps = [15]
        add_to_res()
        # 如果邻接节点已被处理，则跳过
        if color[v] != 'gray':
            highlightSteps = [16]
            add_to_res()
            linkColors[l] = tmpColor
            add_to_res()
        # 如果邻接节点也被处理，则处理该节点
        if color[v] == 'gray':
            highlightSteps = [16]
            add_to_res()
            pred[v] = u
            l = find_link_label(u, v)
            linkColors[l] = 'steelblue'
            highlightSteps = [17, 18]
            add_to_res()
            dfs_visit(v)
    # 该节点已被处理完成
    color[u] = 'steelblue'
    nodeColors[u] = 'steelblue'
    highlightSteps = [19, 20, 21]
    add_to_res()
    time += 1
    f[u] = time
    highlightSteps = [22]
    add_to_res()
    return


start = 0
highlightSteps = [24]
add_to_res()
dfs(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))