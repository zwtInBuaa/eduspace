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


# 定义一个函数来实现图的广度优先遍历
def kruskal():
    global highlightSteps
    highlightSteps = [3]
    add_to_res()

    highlightSteps = [4, 5]
    add_to_res()
    parent = [v for v in V]
    rank = [0] * len(V)

    def find(x):
        if parent[x] != x:
            parent[x] = find(parent[x])
        return parent[x]

    def union(x, y):
        px, py = find(x), find(y)
        if px == py:
            return False
        if rank[px] < rank[py]:
            parent[px] = py
        elif rank[px] > rank[py]:
            parent[py] = px
        else:
            parent[px] = py
            rank[py] += 1
        return True

    highlightSteps = [6, 7, 8]
    add_to_res()

    edges = []
    for i in range(len(links)):
        edges.append(
            {'source': trans[links[i]['source']], 'target': trans[links[i]['target']], 'weight': links[i]['weight']}
        )
    edges.sort(key=lambda x: x['weight'])

    mst = []
    weights = 0
    highlightSteps = [9, 10]
    add_to_res()
    if len(edges) == 0:
        highlightSteps = [11]
        add_to_res()
    for e in edges:
        # 按照边权从小到大的顺序遍历待选边集合
        l = find_link_label(e['source'], e['target'])
        tmpColor = linkColors[l]
        linkColors[l] = 'lightblue'
        highlightSteps = [11]
        add_to_res()

        unionResult = union(e['source'], e['target'])
        # 如果所选的边加入后形成环路，则放弃选取
        if not unionResult:
            highlightSteps = [12]
            add_to_res()
            linkColors[l] = tmpColor
            add_to_res()
        # 否则加入生成的最小生成树中
        if unionResult:
            highlightSteps = [12]
            add_to_res()
            mst.append(e)
            weights += e['weight']
            nodeColors[e['source']] = 'steelblue'
            nodeColors[e['target']] = 'steelblue'
            l = find_link_label(e['source'], e['target'])
            linkColors[l] = 'steelblue'
            highlightSteps = [13, 14]
            add_to_res()
    highlightSteps = [15]
    add_to_res()
    return weights, mst


highlightSteps = [17]
add_to_res()
kruskal()
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))