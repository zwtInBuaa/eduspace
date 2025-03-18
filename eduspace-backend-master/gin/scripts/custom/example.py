def hanoi(n, source, target, auxiliary):
    """
    递归函数，将 n 个盘子从 source 移动到 target 上，借助 auxiliary
    """
    if n == 1:
        print(f"Move disk 1 from {source} to {target}")
        return

    hanoi(n - 1, source, auxiliary, target)
    print(f"Move disk {n} from {source} to {target}")
    hanoi(n - 1, auxiliary, target, source)


# 调用 hanoi 函数进行测试
hanoi(4, 'A', 'C', 'B')
