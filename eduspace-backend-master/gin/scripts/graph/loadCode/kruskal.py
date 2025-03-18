import copy
import queue
import sys
import json

code = '''# 采用并查集判断和维护所选边的顶点是否在一棵子树
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

inputDefault = []

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


# 定义一个函数来实现最小生成树kruskal算法
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
    edges = [e for e in links]
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