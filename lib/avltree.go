/*
 * @Date: 2021-11-18 10:00:59
 * @LastEditTime: 2021-11-18 15:53:52
 * @FilePath: \gnet_server\lib\avltree.go
 * @Description: 实现自平衡二叉树
 */
package lib

import (
	"strings"
)

// 定义树节点
type AVLNode struct {
	Key         string
	Value       interface{}
	Left, Right *AVLNode
	Height      int
}

// 定义树
type AVLTree struct {
	Root *AVLNode
}

// 初始化树
func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

// 判断树中是否存在元素
func (t *AVLTree) Contains(key string) bool {
	if t.Root == nil {
		return false
	}

	return t.Root.contains(key)
}

// 子节点操作定义
func (n *AVLNode) contains(key string) bool {
	if strings.Compare(key, n.Key) == 0 {
		return true
	}

	if strings.Compare(key, n.Key) == -1 && n.Left != nil {
		return n.Left.contains(key)
	}

	if strings.Compare(key, n.Key) == 1 && n.Right != nil {
		return n.Right.contains(key)
	}

	return false
}

// 获取元素
func (t *AVLTree) Get(key string) interface{} {
	if t.Root == nil {
		return nil
	}

	return t.Root.get(key)
}

// 子节点操作定义
func (n *AVLNode) get(key string) interface{} {
	if strings.Compare(key, n.Key) == 0 {
		return n.Value
	}

	if strings.Compare(key, n.Key) == -1 && n.Left != nil {
		return n.Left.get(key)
	}

	if strings.Compare(key, n.Key) == 1 && n.Right != nil {
		return n.Right.get(key)
	}

	return nil
}

// 新增元素
func (t *AVLTree) Set(key string, value interface{}) {
	t.Root = t.Root.set(key, value)
}

// 子节点操作定义
func (n *AVLNode) set(key string, value interface{}) *AVLNode {
	if n == nil {
		return &AVLNode{Key: key, Value: value, Height: 1}
	}

	if strings.Compare(key, n.Key) == 0 {
		n.Value = value
		return n
	}

	var _tmpNode *AVLNode

	if strings.Compare(key, n.Key) == -1 {
		n.Left = n.Left.set(key, value)
		// 判断是否平衡
		var bf = n.balanceFactor()
		if bf == 2 {
			if strings.Compare(key, n.Left.Key) == 1 {
				// 右旋
				_tmpNode = RightRotate(n)
			} else {
				// 左旋
				_tmpNode = LeftRotate(n)
			}
		}
	} else {
		n.Right = n.Right.set(key, value)
		// 判断是否平衡
		var bf = n.balanceFactor()
		if bf == -2 {
			if strings.Compare(key, n.Right.Key) == 1 {
				// 左旋
				_tmpNode = LeftRotate(n)
			} else {
				// 右旋
				_tmpNode = RightRotate(n)
			}
		}
	}

	if _tmpNode == nil {
		// 更新子树高度
		n.updateHeight()
		return n
	} else {
		// 更新子树高度
		_tmpNode.updateHeight()
		return _tmpNode
	}
}

func (n *AVLNode) updateHeight() {
	if n == nil {
		return
	}

	leftHeight, rightHeight := 0, 0
	if n.Left != nil {
		leftHeight = n.Left.Height
	}
	if n.Right != nil {
		rightHeight = n.Right.Height
	}

	var maxHeight = leftHeight
	if rightHeight > maxHeight {
		maxHeight = rightHeight
	}

	n.Height = maxHeight + 1
}

func (n *AVLNode) balanceFactor() int {
	leftHeight, rightHeight := 0, 0
	if n.Left != nil {
		leftHeight = n.Left.Height
	}
	if n.Right != nil {
		rightHeight = n.Right.Height
	}

	return leftHeight - rightHeight
}

// RightRotate 右旋操作
func RightRotate(node *AVLNode) *AVLNode {
	pivot := node.Left    // pivot 表示新插入的节点
	pivotR := pivot.Right // 暂存 pivot 右子树入口节点
	pivot.Right = node    // 右旋后最小不平衡子树根节点 node 变成 pivot 的右子节点
	node.Left = pivotR    // 而 pivot 原本的右子节点需要挂载到 node 节点的左子树上

	// 只有 node 和 pivot 的高度改变了
	node.updateHeight()
	pivot.updateHeight()

	// 返回右旋后的子树根节点指针，即 pivot
	return pivot
}

// LeftRotate 左旋操作
func LeftRotate(node *AVLNode) *AVLNode {
	pivot := node.Right  // pivot 表示新插入的节点
	pivotL := pivot.Left // 暂存 pivot 左子树入口节点
	pivot.Left = node    // 左旋后最小不平衡子树根节点 node 变成 pivot 的左子节点
	node.Right = pivotL  // 而 pivot 原本的左子节点需要挂载到 node 节点的右子树上

	// 只有 node 和 pivot 的高度改变了
	node.updateHeight()
	pivot.updateHeight()

	// 返回旋后的子树根节点指针，即 pivot
	return pivot
}

// LeftRightRotation 双旋操作（先左后右）
func LeftRightRotation(node *AVLNode) *AVLNode {
	node.Left = LeftRotate(node.Left)
	return RightRotate(node)
}

// RightLeftRotation 先右旋后左旋
func RightLeftRotation(node *AVLNode) *AVLNode {
	node.Right = RightRotate(node.Right)
	return LeftRotate(node)
}

// Delete 删除指定节点
func (tree *AVLTree) Remove(key string) {
	// 空树直接返回
	if tree.Root == nil {
		return
	}

	// 删除指定节点，和插入节点一样，根节点也会随着 AVL 树的旋转动态变化
	tree.Root = tree.Root.delete(key)
}

func (n *AVLNode) delete(key string) *AVLNode {
	// 空节点直接返回 nil
	if n == nil {
		return nil
	}
	if strings.Compare(key, n.Key) == -1 {
		// 如果删除节点值小于当前节点值，则进入当前节点的左子树删除元素
		n.Left = n.Left.delete(key)
		// 删除后要更新左子树高度
		n.Left.updateHeight()
	} else if strings.Compare(key, n.Key) == 1 {
		// 如果删除节点值大于当前节点值，则进入当前节点的右子树删除元素
		n.Right = n.Right.delete(key)
		// 删除后要更新右子树高度
		n.Right.updateHeight()
	} else {
		// 找到待删除节点后
		// 第一种情况，删除的节点没有左右子树，直接删除即可
		if n.Left == nil && n.Right == nil {
			// 返回 nil 表示直接删除该节点
			return nil
		}

		if n.Left != nil && n.Right != nil {
			// 左子树更高，拿左子树中值最大的节点替换
			if n.Left.Height > n.Right.Height {
				maxNode := n.Left
				for maxNode.Right != nil {
					maxNode = maxNode.Right
				}

				// 将值最大的节点值赋值给待删除节点
				n.Key = maxNode.Key
				n.Value = maxNode.Value

				// 然后把值最大的节点删除
				n.Left = n.Left.delete(maxNode.Key)
				// 删除后要更新左子树高度
				n.Left.updateHeight()
			} else {
				// 右子树更高，拿右子树中值最小的节点替换
				minNode := n.Right
				for minNode.Left != nil {
					minNode = minNode.Left
				}

				// 将值最小的节点值赋值给待删除节点
				n.Key = minNode.Key
				n.Value = minNode.Value

				// 然后把值最小的节点删除
				n.Right = n.Right.delete(minNode.Key)
				// 删除后要更新右子树高度
				n.Right.updateHeight()
			}
		} else {
			// 只有左子树或只有右子树
			// 只有一棵子树，该子树也只是一个节点，则将该节点值赋值给待删除的节点，然后置子树为空
			if n.Left != nil {
				// 第三种情况，删除的节点只有左子树
				// 根据 AVL 树的特征，可以知道左子树其实就只有一个节点，否则高度差就等于 2 了
				n.Key = n.Left.Key
				n.Value = n.Left.Value
				n.Height = 1
				n.Left = nil
			} else if n.Right != nil {
				// 第四种情况，删除的节点只有右子树
				// 根据 AVL 树的特征，可以知道右子树其实就只有一个节点，否则高度差就等于 2 了
				n.Key = n.Right.Key
				n.Value = n.Right.Value
				n.Height = 1
				n.Right = nil
			}
		}

		// 找到指定值对应的待删除节点并进行替换删除后，直接返回该节点
		return n
	}

	// 左右子树递归删除节点后需要平衡
	var newNode *AVLNode
	// 相当删除了右子树的节点，左边比右边高了，不平衡
	if n.balanceFactor() == 2 {
		if n.Left.balanceFactor() >= 0 {
			newNode = RightRotate(n)
		} else {
			newNode = LeftRightRotation(n)
		}
		//  相当删除了左子树的节点，右边比左边高了，不平衡
	} else if n.balanceFactor() == -2 {
		if n.Right.balanceFactor() <= 0 {
			newNode = LeftRotate(n)
		} else {
			newNode = RightLeftRotation(n)
		}
	}

	if newNode == nil {
		n.updateHeight()
		return n
	} else {
		newNode.updateHeight()
		return newNode
	}
}

func (t *AVLTree) IteratorAsc(f func(k interface{}, v interface{}) bool) {
	t.Root.traverse(f)
}

func (n *AVLNode) traverse(f func(k interface{}, v interface{}) bool) {
	if n == nil {
		return
	}
	// 否则先从左子树最左侧节点开始遍历
	n.Left.traverse(f)
	// 打印位于中间的根节点
	f(n.Key, n.Value)
	// 最后按照和左子树一样的逻辑遍历右子树
	n.Right.traverse(f)
}
