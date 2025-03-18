r ## 摘要

这是机器学习预测学生弱点的程序，采用的是xgboost分类模型，训练采用随机生成数据，其中弱点项为数据中正确率低于0.60对应的题目类型类型，数据量为5000条。

## 各个文件的含义

- **gen_data.py**，用于生成随机数据。

- **student_questions_acc.csv**，生成的数据集，用于训练。

- **student_questions_to_predict.csv**，待预测学生弱点的数据集，相较于**student_questions_acc.csv**，不同在于没有weakness这一项，这也是该程序的**输入**。

- **train.xlsx**这是训练数据预测的结果，可以忽略。

- **predict.xlsx**这是预测的结果。

- **main.py**，机器学习模型主文件，主要包括**xgboost**和**predict**两个函数。

  - **xgboost(train_set_path, save_set_path)**，机器学习训练函数，预测的数据从**predict_set_path**读入，并将训练结果到**save_test_path**中。训练的数据为学生每种题型的正确率和真实的弱点项，示例数据如下。预测目标为弱点项，由于是多项字符串型，设计状态压缩函数**strs2int64(df, target)**，训练的结果存放在**predict.xlsx**中，并会将训练好的模型存放在**xgboost_classifier_model.model**在，这保证了一次训练，多次使用。

  ```csv
  branch,loop,binary_search,sorting,graph_traversal,graph_path,tree,weakness
  0.94,0.66,0.68,0.88,0.52,0.69,0.76,graph_traversal
  ```

  - **predict(predict_set_path, save_test_path)**，预测函数，预测的数据从**predict_set_path**读入，并将预测结果保存到**save_test_path**中，当然控制台也会输出结果。主要是调用**xgboost**生成的**xgboost_classifier_model.model**模型进行预测。

- **xgboost_classifier_model.model**，生成的预测模型。

## 使用方法

如果已经训练好了，只要将训练数据放到**student_questions_to_predict.csv**中，运行**main.py**中的**predict**函数即可；

如果没有训练好，需要先运行**main.py**中的**xgboost**函数获得模型。

## 注意事项

修改数据项的类型后，需要重新训练，并修改**main.py**的状态压缩字典**dic**。