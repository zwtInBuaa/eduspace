import copy
import queue
import sys
import json

code = '''# color表示顶点状态，lightblue表示正在被处理，steelblue表示已覆盖，gray表示未覆盖
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
    {"id": 8}
]
links = [
    {"source": 0, "target": 1, "weight": 4},
    {"source": 0, "target": 6, "weight": 8},
    {"source": 1, "target": 2, "weight": 8},
    {"source": 1, "target": 6, "weight": 1},
    {"source": 2, "target": 3, "weight": 7},
    {"source": 2, "target": 4, "weight": 4},
    {"source": 2, "target": 7, "weight": 2},
    {"source": 3, "target": 4, "weight": 14},
    {"source": 3, "target": 8, "weight": 9},
    {"source": 4, "target": 5, "weight": 2},
    {"source": 4, "target": 8, "weight": 10},
    {"source": 5, "target": 6, "weight": 1},
    {"source": 5, "target": 7, "weight": 4},
    {"source": 6, "target": 7, "weight": 7}
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


def prim(start):
    global highlightSteps
    highlightSteps = [3]
    add_to_res()
    dist[start] = 0
    highlightSteps = [4]
    add_to_res()
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

        # 最小的安全边被选上
        l = find_link_label(pred[rec], rec)
        if l >= 0:
            linkColors[l] = 'steelblue'
        nodeColors[rec] = 'lightblue'
        # rec节点即将被处理
        color[rec] = 'lightblue'
        highlightSteps = [12]
        add_to_res()

        linked_u = E[rec]
        for u in linked_u:
            # 当前所遍历的邻接边
            l = find_link_label(rec, u[0])
            assert l >= 0
            tmpColor = linkColors[l]
            linkColors[l] = 'lightblue'
            highlightSteps = [13]
            add_to_res()
            # 如果这条边不符合更新条件，则还原
            if not (u[1] < dist[u[0]] and color[u[0]] != 'steelblue'):
                highlightSteps = [14]
                add_to_res()
                linkColors[l] = tmpColor
                add_to_res()
            # 否则，将其更新到安全边集合中
            if u[1] < dist[u[0]] and color[u[0]] != 'steelblue':
                highlightSteps = [14]
                add_to_res()
                if dist[u[0]] != sys.maxsize:
                    l = find_link_label(pred[u[0]], u[0])
                    linkColors[l] = '#d3d3d3'
                l = find_link_label(rec, u[0])
                linkColors[l] = 'green'
                dist[u[0]] = u[1]
                pred[u[0]] = rec
                highlightSteps = [15, 16]
                add_to_res()
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
prim(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))