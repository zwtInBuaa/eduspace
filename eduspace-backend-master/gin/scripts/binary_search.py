import copy
import json

code = '''def binarySearch(array, target):
    left = 0
    right = len(array) - 1

    while ___:
        mid = (left + right) // 2
        if array[mid] == target:
            return mid
        elif array[mid] < target:
            left = ___
        else:
            right = ___
    return -1

array = [___]
target = ___
result = binarySearch(array, target)
print(result)'''


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

array = [$$CodeBlank[3]$$]
if len(array) == 0:
    exit("Array can't be empty!")
if len(array) > 15:
    exit("Array length is to big!")
if not all(isinstance(x, int) for x in array):
    exit("Array elements should all be integers!")
for item in array:
    if item >= 500 or item <= -500:
        exit("Array elements should all be > -500 and < 500!")
if array != sorted(array) and array != reversed(sorted(array)):
    exit("Array should be sorted!")
length = len(array)

res = {
    "code": code,
    "barNum": len(array),
    "steps": [],
    "barValue": [],
    "barHeight": [],
    "highlightSteps": [],
    "varSteps": []
}

for item in array:
    res["barValue"].append(item)
res["barHeight"] = normalize(res["barValue"])


def add_to_res():
    global curStep, curhighlightSteps, curVarSteps
    if len(res["steps"]) > 0 and res["steps"][-1] == curStep \
            and res["highlightSteps"][-1] == curhighlightSteps \
            and res["varSteps"][-1] == curVarSteps:
        return
    res["steps"].append(curStep)
    res["highlightSteps"].append(curhighlightSteps)
    res["varSteps"].append(copy.deepcopy(curVarSteps))
    curhighlightSteps = []


def equals(arr1, arr2):
    if len(arr1) != len(arr2):
        return False
    length = len(arr1)
    for i in range(length):
        if arr1[i] != arr2[i]:
            return False
    return True


def addValueToCurVarSteps(left=-1, mid=-1, right=-1, curMidValue=-1, target=-1, circle=0):
    global curVarSteps
    curVarSteps["values"][0] = left
    curVarSteps["values"][1] = mid
    curVarSteps["values"][2] = right
    curVarSteps["values"][3] = curMidValue
    curVarSteps["values"][4] = target
    curVarSteps["values"][5] = circle
    return


def binary_search(array, target):
    circle = 0
    global curStep, curhighlightSteps, curVarSteps
    curStep = ['lightblue' for _ in range(len(array))]
    curhighlightSteps = [0]
    add_to_res()
    left = 0
    right = len(array) - 1
    mid = (left + right) // 2
    curStep = ['lightblue' for _ in range(len(array))]
    curhighlightSteps = [1, 2]
    addValueToCurVarSteps(left, -1, right, -1, target, circle)
    add_to_res()
    lastL, lastM, lastR = -1, -1, -1
    while $$CodeBlank[0]$$:
        circle += 1
        curhighlightSteps = [4]
        add_to_res()
        lastL, lastM, lastR = left, mid, right
        mid = (left + right) // 2
        cur_step = ['lightblue' for _ in range(len(array))]
        cur_step[left] = 'yellow'
        cur_step[right] = 'yellow'
        cur_step[mid] = 'red'  # 优先显示red
        curStep = cur_step
        curhighlightSteps = [5]
        addValueToCurVarSteps(left, mid, right, array[mid], target, circle)
        add_to_res()
        if not (array[mid] == target):
            curStep = cur_step
            curhighlightSteps = [6]
            add_to_res()
        if array[mid] == target:
            curStep = cur_step
            curhighlightSteps = [6]
            add_to_res()
            targetFind = ['lightblue' for _ in range(len(array))]
            targetFind[mid] = 'green'
            curStep = targetFind
            curhighlightSteps = [7]
            add_to_res()
            return res
        if array[mid] >= target:
            curStep = cur_step
            curhighlightSteps = [8]
            add_to_res()
        if array[mid] < target:
            curStep = cur_step
            curhighlightSteps = [8]
            add_to_res()
            left = $$CodeBlank[1]$$
            curStep = cur_step
            curhighlightSteps = [9]
            addValueToCurVarSteps(left, mid, right, array[mid], target, circle)
            add_to_res()
        else:
            curStep = cur_step
            curhighlightSteps = [10]
            add_to_res()
            right = $$CodeBlank[2]$$
            curStep = cur_step
            curhighlightSteps = [11]
            addValueToCurVarSteps(left, mid, right, array[mid], target, circle)
            add_to_res()
            # res["steps"].append(cur_step)
            # fail to find
        if lastL == left and right == lastR and mid == lastM:
            break
    curStep = ['black' for _ in range(len(array))]
    curhighlightSteps = [12]
    add_to_res()
    return res


curStep = ['lightblue' for _ in range(len(array))]
curhighlightSteps = [14]
curVarSteps = {
    'varNames': ["left", "mid", "right", "array[mid]", "target", "circle times"],
    'values': [-1, -1, -1, -1, -1, 0]
}
add_to_res()
target = $$CodeBlank[4]$$
curStep = ['lightblue' for _ in range(len(array))]
curhighlightSteps = [15]
addValueToCurVarSteps(-1, -1, -1, -1, target, 0)
add_to_res()
curStep = ['lightblue' for _ in range(len(array))]
curhighlightSteps = [16]
add_to_res()
tmpResSteps = [['lightblue' for _ in range(len(array))]] * 2 + res["steps"]
# print(tmpResSteps)
res["steps"] = tmpResSteps
print(json.dumps(binary_search(array, target)))
