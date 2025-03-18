import random

dic = ['sorting',
       'graph',
       'tree',
       'partition',
       'basic_knowledge',
       'loop_and_branch',
       'simulate']

file = open('student_questions_acc.csv', 'w')
file.write('sorting,graph,tree,partition,basic_knowledge,loop_and_branch,simulate,weakness\n')
for i in range(10000):
    all = []
    weakness = []
    for j in range(len(dic)):
        random_number = round(random.uniform(0.43, 0.99), 2)
        all.append(str(random_number))
        if random_number <= 0.59:
            weakness.append(dic[j])
    front = ','.join(all)
    behind = ';'.join(weakness)
    line = front + ',' + behind + '\n'
    file.write(line)
file.close()
