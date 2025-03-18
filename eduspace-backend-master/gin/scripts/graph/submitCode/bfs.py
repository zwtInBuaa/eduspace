import copy
import queue
import sys
import json

graphInfo = json.loads('''$$GraphInfo$$''')
nodes = graphInfo["nodes"]
links = graphInfo["links"]

# 初始化图结构
v_nums = len(nodes)
e_nums = len(links)

trans = [-1 for _ in range(30)]  # 最多30个顶点 trans[i] = j 表示id为i的节点对应数组下标为 j
for i in range(v_nums):
    trans[nodes[i]['id']] = i
# V = [v['id'] for v in nodes]
V = [_ for _ in range(v_nums)]
E = {}

# 存放所有顶点的邻接边信息
for i in range(v_nums):
    E[V[i]] = []
for i in range(e_nums):
    E[trans[links[i]["source"]]].append([trans[links[i]["target"]], links[i]["weight"]])
    E[trans[links[i]["target"]]].append([trans[links[i]["source"]], links[i]["weight"]])

# 辅助数组
color = ['gray' for _ in range(v_nums)]
pred = [-1 for _ in range(v_nums)]
int_max = sys.maxsize
dist = [int_max for _ in range(v_nums)]

res = {
    "nodeColors": [],
    "linkColors": [],
    "linkTextColors": [],
    "nodeLabelColors": [],
    "highlightSteps": [],
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


start = $$CodeBlank[0]$$
nodeIDs = [t['id'] for t in nodes]
if start not in nodeIDs:
    exit("error")
# default start = 0
start = trans[start]
highlightSteps = [24]
add_to_res()
bfs(start)
highlightSteps = []
add_to_res()
print(json.dumps(res, ensure_ascii=False))
# file = open("../../res.json", 'w')
# file.write(json.dumps(res))