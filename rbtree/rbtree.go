package rbtree

type nodetype bool

type Comparator[T any] func(src T, dst T) int

const (
	red   nodetype = true
	black nodetype = false
)

type node[T any] struct {
	typ    nodetype
	data   T
	left   *node[T]
	right  *node[T]
	parent *node[T]
}

// isLeftChild 判断当前节点是否是其父节点的左孩子
func (n *node[T]) isLeftChild() bool {
	return n.parent != nil && n.parent.left == n
}

// isRightChild 判断当前节点是否是其父节点的右孩子
func (n *node[T]) isRightChild() bool {
	return n.parent != nil && n.parent.right == n
}

// sibling 查找当前节点的兄弟节点
func (n *node[T]) sibling() *node[T] {
	if n.isLeftChild() {
		return n.parent.right
	}
	if n.isRightChild() {
		return n.parent.left
	}
	return nil
}

// findSuccessor 寻找后继
func (n *node[T]) findSuccessor() *node[T] {
	if n == nil {
		return nil
	}
	if n.right != nil {
		root := n.right
		for root.left != nil {
			root = root.left
		}
		return root
	}

	root := n
	for root.isRightChild() {
		root = root.parent
	}
	return root.parent
}

type RedBlackTree[T any] struct {
	compare Comparator[T]
	root    *node[T]
}

func NewNode[T any](val T) *node[T] {
	return &node[T]{
		typ:  red,
		data: val,
	}
}

func NewRedBlackTree[T any](compare Comparator[T]) *RedBlackTree[T] {
	return &RedBlackTree[T]{
		compare: compare,
	}
}

// Insert 向红黑树里面插入节点
func (r *RedBlackTree[T]) Insert(val T) {
	node := r.insert(r.root, val)
	r.adjustAfterInsert(node)
}

func (r *RedBlackTree[T]) insert(root *node[T], val T) *node[T] {
	new_node := NewNode[T](val)
	if root == nil {
		r.root = new_node
		return new_node
	}
	if r.compare(root.data, val) == -1 {
		if root.right == nil {
			root.right = new_node
			root.right.parent = root
		} else {
			return r.insert(root.right, val)
		}
	} else {
		if root.left == nil {
			root.left = new_node
			root.left.parent = root
		} else {
			return r.insert(root.left, val)
		}
	}
	return new_node
}

// adjustAfterInsert 调整新加入节点的树，使其满足红黑树的性质
func (r *RedBlackTree[T]) adjustAfterInsert(n *node[T]) {
	parent := n.parent
	// 节点没有父亲说明为root直接染色成黑色
	if parent == nil {
		r.black(n)
		return
	}
	// 节点的父亲为黑色直接添加
	if r.isBlack(parent) {
		return
	}
	uncle := parent.sibling()
	grand := parent.parent
	// 叔父节点为红色
	if r.isRed(uncle) {
		r.black(parent)
		r.black(uncle)
		r.red(grand)
		r.adjustAfterInsert(grand)
		return
	}
	// 叔父节点为黑色
	if parent.isLeftChild() {
		if n.isLeftChild() {
			// ll
			r.black(parent)
			r.red(grand)
			r.rightRotate(grand)
		} else {
			// lr
			r.black(n)
			r.red(grand)
			r.leftRotate(parent)
			r.rightRotate(grand)
		}
	} else {
		if n.isLeftChild() {
			// rl
			r.black(n)
			r.red(grand)
			r.rightRotate(parent)
			r.leftRotate(grand)
		} else {
			// rr
			r.black(parent)
			r.red(grand)
			r.leftRotate(grand)
		}
	}
}

// isBlack 判断当前节点是不是黑色节点
func (r *RedBlackTree[T]) isBlack(n *node[T]) bool {
	if n == nil {
		return true
	}
	if n.typ == black {
		return true
	}
	return false
}

// isRed 判断当前节点是不是红色节点
func (r *RedBlackTree[T]) isRed(n *node[T]) bool {
	if n == nil {
		return false
	}
	if n.typ == red {
		return true
	}
	return false
}

// black 将节点染色成黑色
func (r *RedBlackTree[T]) black(n *node[T]) {
	n.typ = black
}

// red 将节点染色成红色
func (r *RedBlackTree[T]) red(n *node[T]) {
	n.typ = red
}

// leftRotate 左旋
func (r *RedBlackTree[T]) leftRotate(grand *node[T]) {
	parent := grand.right
	child := parent.left
	grand.right = child
	parent.left = grand

	parent.parent = grand.parent
	if grand.isLeftChild() {
		grand.parent.left = parent
	} else if grand.isRightChild() {
		grand.parent.right = parent
	} else {
		r.root = parent
	}
	if child != nil {
		child.parent = grand
	}
	grand.parent = parent
}

// rightRotate 右旋
func (r *RedBlackTree[T]) rightRotate(grand *node[T]) {
	parent := grand.left
	child := parent.right
	grand.left = parent.right
	parent.right = grand
	parent.parent = grand.parent
	if grand.isLeftChild() {
		grand.parent.left = parent
	} else if grand.isRightChild() {
		grand.parent.right = parent
	} else {
		r.root = parent
	}
	if child != nil {
		child.parent = grand
	}
	grand.parent = parent
}

// color 染色
func (r *RedBlackTree[T]) color(n *node[T], typ nodetype) {
	n.typ = typ
}

// get 获取红黑树中的节点,第一返回值为节点，第二个返回值是是否存在值为val的节点
func (r *RedBlackTree[T]) get(n *node[T], val T) (*node[T], bool) {
	if n == nil {
		return nil, false
	}
	if r.compare(n.data, val) == 0 {
		return n, true
	} else if r.compare(n.data, val) == -1 {
		return r.get(n.right, val)
	}
	return r.get(n.left, val)

}

// Delete 返回红黑树是否存在val这个节点
func (r *RedBlackTree[T]) Delete(val T) bool {
	n, ok := r.get(r.root, val)
	if !ok {
		return false
	}
	if n.left != nil && n.right != nil {
		s := n.findSuccessor()
		n.data = s.data
		n = s
	}
	var replacement *node[T]
	if n.left != nil && n.right == nil {
		replacement = n.left
	} else if n.right != nil && n.left == nil {
		replacement = n.right
	}

	if replacement != nil {
		replacement.parent = n.parent
		if n.parent == nil {
			r.root = replacement
		}
		if n.isRightChild() {
			n.parent.right = replacement
		} else if n.isLeftChild() {
			n.parent.left = replacement
		}

		r.adjustAfterDelete(n, replacement)
	} else if n.parent == nil {
		r.root = nil
		r.adjustAfterDelete(n, nil)
	} else {
		if n.isRightChild() {
			n.parent.right = nil
		} else {
			n.parent.left = nil
		}
		r.adjustAfterDelete(n, nil)
	}
	return true
}

// adjustAfterDelete 删除元素之后调整树，使其满足红黑树的性质
func (r *RedBlackTree[T]) adjustAfterDelete(n *node[T], replace *node[T]) {
	// 节点为红色直接删除
	if r.isRed(n) {
		return
	}
	// 节点为黑色，替代他的节点为红色
	if r.isRed(replace) {
		r.black(replace)
		return
	}
	parent := n.parent
	if parent == nil {
		return
	}
	var sibling *node[T]
	right := parent.right == nil || n.isRightChild()
	if right {
		sibling = parent.left
	} else {
		sibling = parent.right
	}
	// 当前节点是右孩子
	if right {
		// 兄弟节点是红色需要处理一下是其兄弟节点变成黑色
		if r.isRed(sibling) {
			r.black(sibling)
			r.red(parent)
			r.rightRotate(parent)
			sibling = parent.left
		}
		// 兄弟节点为黑色且兄弟节点没有子节点
		if r.isBlack(sibling.left) && r.isBlack(sibling.right) {
			parentisBlack := r.isBlack(parent)
			r.black(parent)
			r.red(sibling)
			if parentisBlack {
				r.adjustAfterDelete(parent, nil)
			}
		} else {
			// 兄弟节点存在红色的子节点
			if r.isBlack(sibling.left) {
				r.leftRotate(sibling)
				sibling = parent.left
			}
			r.color(sibling, parent.typ)
			r.black(sibling.left)
			r.black(parent)
			r.rightRotate(parent)
		}
	} else {
		if r.isRed(sibling) {
			r.black(sibling)
			r.red(parent)
			r.leftRotate(parent)
			sibling = parent.right
		}
		if r.isBlack(sibling.left) && r.isBlack(sibling.right) {
			parentisBlack := r.isBlack(parent)
			r.black(parent)
			r.red(sibling)
			if parentisBlack {
				r.adjustAfterDelete(parent, nil)
			}
		} else {
			if r.isBlack(sibling.right) {
				r.rightRotate(sibling)
				sibling = parent.right
			}
			r.color(sibling, parent.typ)
			r.black(sibling.right)
			r.black(parent)
			r.leftRotate(parent)
		}
	}
}
