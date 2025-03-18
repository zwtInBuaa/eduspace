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
pred = [-1 for _ in range(v_nums)]
int_max = sys.maxsize
dist = [int_max for _ in range(v_nums)]

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


start = $$CodeBlank[0]$$
nodeIDs = [t['id'] for t in nodes]
if start not in nodeIDs:
    exit("error")
start = trans[start]
highlightSteps = [20]
add_to_res()
prim(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))