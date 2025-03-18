import json
import copy
# import gen_selection_sort_code as gssc


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
    # "code": gssc.gen_selection_sort_code()["code"],
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


def selection_sort():
    global curStep, curDetail, curhighlightSteps
    curhighlightSteps = [0]
    add_to_res()
    curhighlightSteps = [1]
    add_to_res()
    for i in range(length):
        curhighlightSteps = [2]
        add_to_res()
        swapIndex = i
        curhighlightSteps = [3]
        add_to_res()
        for j in range(i + 1, length):
            curhighlightSteps = [4]
            add_to_res()
            # if array[j] < array[swapIndex]:
            if not ($$CodeBlank[0]$$):
                curhighlightSteps = [5]
                add_to_res()
            if $$CodeBlank[0]$$:
                curhighlightSteps = [5]
                add_to_res()
                swapIndex = j
                curhighlightSteps = [6]
                add_to_res()
        array[i], array[swapIndex] = array[swapIndex], array[i]
        if swapIndex != i:
            curStep = [i, swapIndex]
            curhighlightSteps = [7]
            add_to_res()
    curDetail = ['green' for _ in range(length)]
    curhighlightSteps = [8]
    add_to_res()
    return res


curhighlightSteps = [9]
add_to_res()
curhighlightSteps = [10]
add_to_res()
curhighlightSteps = [11]
add_to_res()
selection_sort()
print(json.dumps(res, ensure_ascii=False))