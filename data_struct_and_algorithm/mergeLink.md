# 合并K个升序链表

示例 1：

输入：lists = [[1,4,5],[1,3,4],[2,6]]
输出：[1,1,2,3,4,4,5,6]
解释：链表数组如下：
[
  1->4->5,
  1->3->4,
  2->6
]
将它们合并到一个有序链表中得到。
1->1->2->3->4->4->5->6



```
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeKLists(lists []*ListNode) *ListNode {
    newLists := make([]*ListNode, 0 )
    if len(lists) == 0 {
        return nil 
    } 
    if len(lists) == 1 {
        return lists[0]
    }
    lastList := &ListNode{}
    for {
        if len(lists) == 1 {
            return lists[0]
        }
        //判断奇数的最后一个
        if len(lists)%2 == 1 {
            lastList = lists[len(lists)-1]
        } else {
            lastList = nil 
        }
        for i := 0; i< len(lists)-1;{
                sortLists := mergeTwoLists(lists[i], lists[i+1])
                newLists = append(newLists, sortLists)
                i = i+2
        }   
        if lastList != nil {
            lists = append(newLists, lastList)
        } else {
            lists = newLists
        }
        newLists = make([]*ListNode, 0 )
    }
    
}

func mergeTwoLists(list1, list2 *ListNode) *ListNode{

    if list1 == nil && list2 == nil {
        return nil 
    }

    if list1 != nil && list2 == nil {
        return list1
    }

    if list2 != nil && list1 == nil {
        return list2 
    }

    if list1.Val > list2.Val {
        list2.Next = mergeTwoLists(list1, list2.Next)
        return list2 
    } else {
        list1.Next = mergeTwoLists(list1.Next, list2)
        return list1
    }
}
```
