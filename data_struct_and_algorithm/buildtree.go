# 构造二叉树 

从中序 前序遍历数据，还原二叉树


/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func buildTree(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 || len(inorder) == 0  {
        return nil 
    }
    valMap := make(map[int]int)
    for i,val := range inorder {
        valMap[val] = i 
    }
    node := buildTreeMain(preorder, inorder, valMap)
    return node 
}
func buildTreeMain(preorder []int, inorder []int, valMap map[int]int ) *TreeNode {
    if len(preorder) == 0 || len(inorder) == 0 {
        return nil 
    }
    rootVal := preorder[0]
    node := &TreeNode{Val:rootVal}
    if len(preorder) == 1 {
        return node 
    }
    j := 0 
    lenMax := len(inorder)
    for {
        if j == lenMax {
            break
        }
        if preorder[0] == inorder[j] {
            node.Left  = buildTreeMain(preorder[1:j+1], inorder[:j],valMap)
            node.Right = buildTreeMain(preorder[j+1:], inorder[j+1:],valMap)
            break
        }
        j++
    }
   
    return node 
}