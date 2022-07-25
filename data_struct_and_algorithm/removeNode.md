
# 删除链表的倒数第 N 个结点

给你一个链表，删除链表的倒数第 n 个结点，并且返回链表的头结点。

示例 1：


输入：head = [1,2,3,4,5], n = 2
输出：[1,2,3,5]

```
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func removeNthFromEnd(head *ListNode, n int) *ListNode {
   p , q := head,head 

//快指针。先走N步,这样当快指针到达尾部的时候，慢指针就是需要被删除的
   for i:=0;i<n ; i++ {
       p = p.Next 
   }

   if p == nil {
       return head.Next
   }
//q 移动到需要删除的位置
  for {
      if p.Next == nil {
            break
      }
      p = p.Next
      q = q.Next 
  }
  //删除指定的节点
  q.Next = q.Next.Next

  return head
}
```