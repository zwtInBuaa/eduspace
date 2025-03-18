'''
训练数据特征介绍
sorting : 排序
graph : 图论
tree : 树
partition : 分治
basic_knowledge : 基础知识
loop_and_branch : 循环和分支
simulate : 模拟
weakness : 弱点报告
'''

# 引入库包
import os
import sys

# import joblib
import numpy as np
import pandas as pd
# 为了正确评估模型性能，将数据划分为训练集和测试集，并在训练集上训练模型，在测试集上验证模型性能。
from sklearn.model_selection import train_test_split
# 导入XGBoost模型
from xgboost.sklearn import XGBClassifier
# from sklearn import metrics
import xgboost as xgb

xgb.set_config(verbosity=0)

project_path = os.path.abspath('.') + '/'

dic = {
    'sorting': 0b0000001,
    'graph': 0b0000010,
    'tree': 0b0000100,
    'partition': 0b0001000,
    'basic_knowledge': 0b0010000,
    'loop_and_branch': 0b0100000,
    'simulate': 0b0100000
}

dic_to_chinese = {
    'sorting': '排序',
    'graph': '图论',
    'tree': '树',
    'partition': '分治',
    'basic_knowledge': '基础知识',
    'loop_and_branch': '循环和分支',
    'simulate': '模拟',
}


def strs2int64(df, target):
    list = []
    for v in df[target]:
        num = 0
        if not str(v) == 'nan':
            S = str(v).split(";")
            for s in S:
                num |= dic[s]
        list.append(num)
    df[target] = list


def int64tostr(df, target):
    all = []
    for v in df[target]:
        list = []
        for key in dic:
            if (v & dic[key]) != 0:
                list.append(key)
        a = ';'.join(list)
        all.append(a)
    df[target] = all


def xgboost(train_set_path, save_set_path):
    df = pd.read_csv(project_path + train_set_path)
    # 区分数字特征与非数字特征
    numerical_features = [x for x in df.columns if df[x].dtype == np.int64]
    category_features = [x for x in df.columns if df[x].dtype != np.int64 and x != 'weakness']

    for i in category_features:
        df[i] = df[i].apply(lambda x: int(x * 10))
    #  缺点用状态压缩进行编码
    strs2int64(df, 'weakness')

    data_target_part = df['weakness']

    data_features_part = df[[x for x in df.columns if x != 'weakness']]
    # 测试集大小为20%， 80%/20%分
    x_train, x_test, y_train, y_test = train_test_split(data_features_part, data_target_part, test_size=0.2,
                                                        random_state=2020)
    # best fit
    clf = XGBClassifier(gamma=0.27, colsample_bytree=0.685, learning_rate=0.095, max_depth=4, subsample=0.605)

    # 定义 XGBoost模型，注意这里没有参数后期调优
    # clf = XGBClassifier(learning_rate=0.50, max_depth=9, subsample=0.9)

    # 在训练集上训练XGBoost模型
    clf.fit(x_train, y_train)

    # 在训练集和测试集上分布利用训练好的模型进行预测
    test_predict = clf.predict(x_test)

    # 查看混淆矩阵 (预测值和真实值的各类情况统计矩阵)
    # confusion_matrix_result = metrics.confusion_matrix(test_predict, y_test)

    x_all = df.drop('weakness', axis=1)
    y_all = df['weakness']
    x_all = pd.get_dummies(x_all)

    clf.fit(x_all, y_all)
    clf.save_model('xgboost_classifier_model.model')
    y_predict = clf.predict(x_all)

    df = pd.read_csv(project_path + train_set_path)
    df['Predict'] = y_predict
    int64tostr(df, 'Predict')
    df.to_excel(project_path + save_set_path)
    return


def predict(sorting, graph, tree, partition, basic_knowledge, loop_and_branch, simulate):
    acc_all = [sorting, graph, tree, partition, basic_knowledge, loop_and_branch, simulate]
    df = pd.DataFrame(columns=list(dic.keys()))
    df.loc[len(df)] = acc_all

    numerical_features = [x for x in df.columns if df[x].dtype == np.int64]
    category_features = [x for x in df.columns if df[x].dtype != np.int64 and x != 'weakness']

    for i in category_features:
        df[i] = df[i].apply(lambda x: int(x * 10))

    return model_save_load('xgboost_classifier_model.model', df)


def model_save_load(model, x_transform):
    # 模型加载
    clf = xgb.XGBClassifier()
    booster = xgb.Booster()
    booster.load_model(model)
    clf._Booster = booster

    # 数据预测
    y_pred = clf.predict(x_transform)
    # print(y_pred)
    all_weakness = []
    for key in dic:
        if (y_pred & dic[key]) != 0:
            all_weakness.append(dic_to_chinese[key])
    sys.stdout.reconfigure(encoding='utf-8')
    print(all_weakness)


# 训练模型，模型已经训练并保存，无需调用这个函数
# xgboost('student_questions_acc.csv', 'train.xlsx')
# 示例调用
# predict(0.73, 0.60, 0.82, 0.41, 0.92, 0.37, 0.70)
predict(float(sys.argv[1]), float(sys.argv[2]), float(sys.argv[3]), float(sys.argv[4]), float(sys.argv[5]),
        float(sys.argv[6]), float(sys.argv[7]))
