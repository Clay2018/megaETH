# nft_element_aggregation_contract


算法基本原理(O(n^2 * m):
    - 1.找到A token 到 B token 的所有最短通路, 获取路径
    - 2.根据路径求得换的B token数

    - demo目的: 
        获取算法实现所用的gas(因为第2步的gas基本固定, 这里只考虑第1步gas消耗).

第1步实现解析:

<img width="612" alt="image" src="https://user-images.githubusercontent.com/37411084/218982306-616f4755-46cd-466c-b261-6d6f9276973f.png">
    
    - 例如从1到6
    - 路径数共有
        - 1, 2, 4, 8, 5, 6
        - 1, 2, 5, 6
        - 1, 3, 6 
        - 1, 3, 7, 6(抛弃)
    
    - 初始化两个栈,一个主栈(mainStack)存放路径节点，一个附栈(subjectStack)存放邻接节点
    - sptSet集合存放两个栈中已使用的节点, 是一个bool集合

    - 将1压栈进入mainStack；将1的邻接点压入subjectStack
    - 进入for循环, 循环终止条件(subjectStack为空)
        - subjectStack弹出元素，查看该元素是否为空; 如果不为空, 压入mainStack，subjectStack压入一个empty元素 更新sptSet; 如果为空, 
        - 找到主栈最上层元素的所有邻接元素, 经过sptSet去重, 获得新集合adjacentVertex
        - 查看adjacentVertex是否含有dest节点，若有，获得一条路径
        - 若果没有，将相关元素压入附栈

- 参考： https://segmentfault.com/a/1190000020445075

