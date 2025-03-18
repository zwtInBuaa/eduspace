import json
import copy

# from gen_quick_sort_code import gen_quick_sort_code as gqsc

res = {
    "values": [],
    "heights": [],
    "steps": [],
    "detail": []
}


def normalize(arr):
    min_val = min(arr)
    max_val = max(arr)

    normalized_arr = []
    if max_val == min_val:
        normalized_arr = [75 for _ in range(len(arr))]
    else:
        for val in arr:
            normalized_val = round(20 + ((val - min_val) / (max_val - min_val)) * 150)
            normalized_arr.append(normalized_val)

    return normalized_arr


# array = [8, 1, 6, 2, 9, 4, 3, 7, 5]
array = [$$CodeBlank[1]$$]
if len(array) == 0:
    exit("Array can't be empty!")
if len(array) > 15:
    exit("Array length is to big!")
if not all(isinstance(x, int) for x in array):
    exit("Array elements should all be integers!")
for item in array:
    if item >= 500 or item <= -500:
        exit("Array elements should all be > -500 and < 500!")

res = {
    # "code": gqsc()["code"],
    "values": copy.deepcopy(array),
    "heights": normalize(array),
    "steps": [],
    "detail": [],
    "highlightSteps": []
}
length = len(array)

curStep = []
curDetail = []
curhighlightSteps = []


def add_to_res():
    global curStep, curDetail, curhighlightSteps
    if len(res["steps"]) > 0 and res["steps"][-1] == curStep \
            and res["detail"][-1] == curDetail \
            and res["highlightSteps"][-1] == curhighlightSteps:
        return
    res["steps"].append(curStep)
    res["detail"].append(curDetail)
    res["highlightSteps"].append(curhighlightSteps)
    curStep = []
    curDetail = []
    curhighlightSteps = []


def quick_sort(array, left, right):
    global curhighlightSteps
    curhighlightSteps = [0]
    add_to_res()
    if left >= right:
        curhighlightSteps = [1]
        add_to_res()
    if left < right:
        curhighlightSteps = [1]
        add_to_res()
        curhighlightSteps = [2]
        add_to_res()
        pivot = partition(array, left, right)
        curhighlightSteps = [3]
        add_to_res()
        quick_sort(array, left, pivot - 1)
        curhighlightSteps = [4]
        add_to_res()
        quick_sort(array, pivot + 1, right)
    curhighlightSteps = [5]
    add_to_res()
    return


def partition(arr, left, right):
    global length, curStep, curDetail, curhighlightSteps
    curhighlightSteps = [7]
    add_to_res()
    #   每次子迭代只会给本次迭代的元素赋色
    color = ['white' for _ in range(length)]
    for i in range(left, right):
        color[i] = 'lightblue'
    #   锚点元素设置成红色
    color[right] = 'red'
    pivot = arr[right]
    curDetail = color
    curhighlightSteps = [8]
    add_to_res()
    i = left - 1
    curhighlightSteps = [9]
    add_to_res()
    for j in range(left, right):
        curhighlightSteps = [10]
        add_to_res()
        # if arr[j] <= pivot:
        if not ($$CodeBlank[0]$$):
            curhighlightSteps = [11]
            add_to_res()
        if $$CodeBlank[0]$$:
            curhighlightSteps = [11]
            add_to_res()
            i += 1
            curhighlightSteps = [12]
            add_to_res()
            arr[i], arr[j] = arr[j], arr[i]
            if i != j:
                '''
                交换步骤：
                1.选中一个比锚点小（大）的元素
                2.选中一个比锚点大（小）的元素
                3.交换他们
                '''
                c1 = copy.deepcopy(color)
                c1[j] = 'black'
                curDetail = c1
                curhighlightSteps = [13]
                add_to_res()

                c2 = copy.deepcopy(c1)
                c2[i] = 'pink'
                curDetail = c2
                curhighlightSteps = [13]
                add_to_res()

                curStep = [i, j]
                curDetail = color
                curhighlightSteps = [13]
                add_to_res()

    arr[i + 1], arr[right] = arr[right], arr[i + 1]
    if i + 1 != right:
        #   最后把锚点移动到合适的位置
        c3 = copy.deepcopy(color)
        c3[i + 1] = 'blue'
        curDetail = c3
        curhighlightSteps = [14]
        add_to_res()

        curStep = [right, i + 1]
        curDetail = color
        curhighlightSteps = [14]
        add_to_res()
    curhighlightSteps = [15]
    add_to_res()
    return i + 1


color = ['white' for _ in range(length)]
for i in range(0, len(array) - 1):
    color[i] = 'lightblue'
#   锚点元素设置成红色
color[len(array) - 1] = 'red'
curDetail = color
curhighlightSteps = [17]
add_to_res()
curhighlightSteps = [18]
add_to_res()
curhighlightSteps = [19]
add_to_res()
quick_sort(array, 0, len(array) - 1)

#   排序结束
curDetail = ['green' for i in range(length)]
add_to_res()

# 返回关键信息
print(json.dumps(res, ensure_ascii=False))