import copy
import queue
import sys
import json

graphInfo = json.loads('''$$GraphInfo$$''')
nodes = graphInfo["nodes"]
links = graphInfo["links"]

v_nums = len(nodes)
e_nums = len(links)
# V = [v['id'] for v in nodes]
trans = [-1 for _ in range(30)]  # 最多20个顶点 trans[i] = j 表示id为i的节点对应数组下标为 j
for i in range(v_nums):
    trans[nodes[i]['id']] = i
V = [_ for _ in range(v_nums)]
E = {}
for i in range(v_nums):
    E[V[i]] = []
for i in range(e_nums):
    E[trans[links[i]["source"]]].append([trans[links[i]["target"]], links[i]["weight"]])
    E[trans[links[i]["target"]]].append([trans[links[i]["source"]], links[i]["weight"]])

color = ['gray' for _ in range(v_nums)]
int_max = sys.maxsize
dist = [int_max for _ in range(v_nums)]
pred = [-1 for _ in range(v_nums)]

res = {
    "nodeColors": [],
    "linkColors": [],
    "linkTextColors": [],
    "nodeLabelColors": [],
    "highlightSteps": []
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
        if trans[links[i]['source']] == a and trans[links[i]['target']] == b \
                or trans[links[i]['source']] == b and trans[links[i]['target']] == a:
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


start = $$CodeBlank[0]$$
nodeIDs = [t['id'] for t in nodes]
if start not in nodeIDs:
    exit("error")
start = trans[start]
highlightSteps = [20]
add_to_res()
dijkstra(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))