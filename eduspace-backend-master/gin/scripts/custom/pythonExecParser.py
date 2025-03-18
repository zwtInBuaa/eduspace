import json
from pytutor import generate_trace
import sys

'''
    1. pip install pytutor
    2. python3 pythonExecParser.py arg1 arg2
        arg1: path of the python file or string of python code
        arg2: True means arg1 is a file path 
              False means arg1 is a string of python code
'''


def trans_heap_info(info, is_recursive):
    global heap_dict
    ret_val = 'UNKNOWN: ' + str(info)
    if not isinstance(info, list):
        # if info isn't a list, return string directly
        ret_val = info
    elif info[0] == 'INSTANCE':
        # eg: INSTANCE: module
        ret_val = info[0] + ': ' + info[1]
    elif info[0] == 'FUNCTION':
        # eg: FUNCTION: add(a, b)
        ret_val = info[0] + ':' + info[1]
    elif info[0] == 'LIST':
        temp_list = list()
        if len(info) > 1:
            for item in info[1:]:
                temp_list.append(trans_heap_info(item, True))
        ret_val = temp_list
    elif info[0] == 'SET':
        temp_set = set()
        if len(info) > 1:
            for item in info[1:]:
                temp_set.add(trans_heap_info(item, True))
        ret_val = temp_set
    elif info[0] == 'TUPLE':
        temp_list = list()
        if len(info) > 1:
            for item in info[1:]:
                temp_list.append(trans_heap_info(item, True))
        ret_val = tuple(temp_list)
    elif info[0] == 'DICT':
        if len(info) > 1:
            ret_val = str({sublist[0]: trans_heap_info(sublist[1], True) for sublist in info[1:]})
        else:
            ret_val = dict()
    elif info[0] == 'REF':
        # if info is a REF, recursively parse this item
        ret_val = trans_heap_info(heap_dict[str(info[1])], True)
    if is_recursive:
        return ret_val
    else:
        return str(ret_val)


def parse_src_code(_trace_str):
    global heap_dict
    trace = _trace_str
    trace_dict = json.loads(trace)
    # code = trace_dict['code']['main_code']
    trace_steps = trace_dict['trace']
    return_steps = []
    highlightSteps = []
    if len(trace_steps) > 1000:
        exit('代码运行步骤过长')
    for item in trace_steps:
        highlightSteps.append([item['line'] - 1])
        step = {'event': item['event']}
        if item['event'] != 'uncaught_exception':
            step['func_name'] = item['func_name']
            step['stdout'] = item['stdout']
        else:
            step['func_name'] = 'Error'
            step['stdout'] = item['exception_msg']
            return_steps.append(step)
            continue
        heap_dict = item['heap']
        globals_vars = []
        if item['globals']:
            for k, v in item['globals'].items():
                if isinstance(v, list):
                    globals_vars.append({'key': k, 'value': trans_heap_info(heap_dict[str(v[1])], False)})
                else:
                    globals_vars.append({'key': k, 'value': v})
        local_vars = []
        if item['stack_to_render']:
            cur_stack = item['stack_to_render'][-1]  # "is_highlighted"
            for stack in item['stack_to_render']:
                if stack["is_highlighted"]:
                    cur_stack = stack
            for k, v in cur_stack['encoded_locals'].items():
                if isinstance(v, list):
                    local_vars.append({'key': k, 'value': trans_heap_info(heap_dict[str(v[1])], False)})
                else:
                    local_vars.append({'key': k, 'value': v})
        step['global_vars'] = globals_vars
        step['local_vars'] = local_vars
        return_steps.append(step)
    return_steps.append(return_steps[-1])
    return {'steps': return_steps, 'highlightSteps': highlightSteps}


if __name__ == '__main__':
    args = sys.argv
    heap_dict = dict()
    source_code = './example.py' if len(args) < 2 else args[1]
    is_file = True if len(args) < 3 else args[2]
    src = open(source_code, encoding='utf-8').read() if is_file else source_code
    trace_str = generate_trace.run_logger(src, '', {})
    # print(json.dumps(json.loads(trace_str)))
    res = parse_src_code(trace_str)  # 默认为True
    print(json.dumps(res))
