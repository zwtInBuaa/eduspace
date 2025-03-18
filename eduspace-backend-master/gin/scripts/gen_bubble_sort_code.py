import json
import copy


# import gen_bubble_sort_code as gbsc

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


array = [8, 1, 6, 2, 9, 4, 3, 7, 5]

res = {
    "code": '''def bubbleSort(array):
    length = len(array)
    for i in range(length - 1):
        for j in range(length - i - 1):
            if ___ :
                array[j], array[j + 1] = array[j + 1], array[j]
    return
array = [___]
bubbleSort(array)
print("Sorted array:", array)''',
    "inputDefault": [
        "array[j] > array[j + 1]",
        "8, 1, 6, 2, 9, 4, 3, 7, 5"
    ],
    # "code": gbsc.gen_bubble_sort_code()["code"],
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


def bubble_sort():
    global curStep, curDetail, curhighlightSteps
    curhighlightSteps = [0]
    add_to_res()
    curhighlightSteps = [1]
    add_to_res()
    # 排序流程
    for i in range(length - 1):
        curhighlightSteps = [2]
        add_to_res()
        count = 0
        for j in range(length - i - 1):
            curhighlightSteps = [3]
            add_to_res()
            # if array[j] > array[j + 1]:
            if not (array[j] > array[j + 1]):
                curhighlightSteps = [4]
                add_to_res()
            if array[j] > array[j + 1]:
                curhighlightSteps = [4]
                add_to_res()
                array[j], array[j + 1] = array[j + 1], array[j]
                curStep = [j, j + 1]
                curhighlightSteps = [5]
                add_to_res()
                count += 1
        if count == 0:
            break
    curDetail = ['green' for _ in range(length)]
    curhighlightSteps = [6]
    add_to_res()
    return res


curhighlightSteps = [7]
add_to_res()
curhighlightSteps = [8]
add_to_res()
curhighlightSteps = [9]
add_to_res()
bubble_sort()
print(json.dumps(res, ensure_ascii=False))