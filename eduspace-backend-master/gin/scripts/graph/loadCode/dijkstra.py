import copy
import queue
import sys
import json

code = '''# color表示顶点状态，lightblue表示正在被处理，steelblue表示最短路已被确定，gray表示未确定
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

inputDefault = ["0"]

nodes = [
    {"id": 0},
    {"id": 1},
    {"id": 2},
    {"id": 3},
    {"id": 4},
    {"id": 5},
    {"id": 6}
]
links = [
    {"source": 0, "target": 1, "weight": 12},
    {"source": 0, "target": 5, "weight": 16},
    {"source": 0, "target": 6, "weight": 14},
    {"source": 1, "target": 2, "weight": 10},
    {"source": 1, "target": 5, "weight": 7},
    {"source": 2, "target": 3, "weight": 3},
    {"source": 2, "target": 4, "weight": 5},
    {"source": 2, "target": 5, "weight": 6},
    {"source": 3, "target": 4, "weight": 4},
    {"source": 4, "target": 5, "weight": 2},
    {"source": 4, "target": 6, "weight": 8},
    {"source": 5, "target": 6, "weight": 9}
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
int_max = sys.maxsize
dist = [int_max for _ in range(v_nums)]
pred = [-1 for _ in range(v_nums)]

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


def dijkstra(start):
    global highlightSteps
    highlightSteps = [3]
    add_to_res()

    highlightSteps = [4]
    add_to_res()
    dist[start] = 0
    if v_nums == 0:
        highlightSteps = [5]
        add_to_res()
    for i in range(v_nums):
        highlightSteps = [5]
        add_to_res()
        minDist = sys.maxsize
        rec = 0
        highlightSteps = [6, 7]
        add_to_res()
        for j in range(v_nums):
            if color[j] != 'steelblue' and dist[j] < minDist:
                minDist = dist[j]
                rec = j
        highlightSteps = [8, 9, 10, 11]
        add_to_res()
        # rec节点即将被处理
        nodeColors[rec] = 'lightblue'
        color[rec] = 'lightblue'
        highlightSteps = [12]
        add_to_res()

        linked_u = E[rec]
        for u in linked_u:
            # 遍历rec节点的所有邻接边
            highlightSteps = [13]
            l = find_link_label(rec, u[0])
            tmpColor = linkColors[l]
            linkColors[l] = 'lightblue'
            add_to_res()
            # 如果e(rec,u)不满足松弛条件，则恢复该边
            if not (dist[rec] + u[1] < dist[u[0]]):
                highlightSteps = [14]
                add_to_res()
                linkColors[l] = tmpColor
                highlightSteps = [14]
                add_to_res()
            # 如果e(rec,u)满足松弛条件，则对该边进行松弛
            if dist[rec] + u[1] < dist[u[0]]:
                highlightSteps = [14]
                add_to_res()
                dist[u[0]] = dist[rec] + u[1]
                pred[u[0]] = rec
                highlightSteps = [15, 16]
                l = find_link_label(rec, u[0])
                for s in E[u[0]]:
                    tmp = find_link_label(u[0], s[0])
                    assert tmp >= 0
                    if linkColors[tmp] == 'steelblue':
                        linkColors[tmp] = '#d3d3d3'
                linkColors[l] = 'steelblue'
                add_to_res()
        # rec已经被处理完成
        nodeColors[rec] = 'steelblue'
        color[rec] = 'steelblue'
        highlightSteps = [17]
        add_to_res()
    highlightSteps = [18]
    add_to_res()
    return


start = 0
highlightSteps = [20]
add_to_res()
dijkstra(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))