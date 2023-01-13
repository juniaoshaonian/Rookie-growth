package rbtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedBlackTree_LeftRotate(t *testing.T) {
	r := NewRedBlackTree[int](intcompare)
	testcases := []struct {
		name         string
		wantTree     func() *node[int]
		grandandroot func() (*node[int], *node[int])
	}{
		{
			name: "grand is root",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newLeftGrand()
				return grand, grand
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n3.left = n1
				n3.right = n5
				n1.parent = n3
				n5.parent = n3
				n1.left = n2
				n1.right = n4
				n2.parent = n1
				n4.parent = n1
				return n3
			},
		},
		{
			name: "grand is root left",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newLeftGrand()
				root := &node[int]{
					data: 6,
				}
				root.left = grand
				grand.parent = root
				return grand, root

			},
			wantTree: func() *node[int] {
				n0 := &node[int]{data: 6}
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n3.left = n1
				n3.right = n5
				n1.parent = n3
				n5.parent = n3
				n1.left = n2
				n1.right = n4
				n2.parent = n1
				n4.parent = n1
				n0.left = n3
				n3.parent = n0
				return n0
			},
		},
		{
			name: "grand is root right",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newLeftGrand()
				root := &node[int]{
					data: 0,
				}
				root.right = grand
				grand.parent = root
				return grand, root

			},
			wantTree: func() *node[int] {
				n0 := &node[int]{data: 0}
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n3.left = n1
				n3.right = n5
				n1.parent = n3
				n5.parent = n3
				n1.left = n2
				n1.right = n4
				n2.parent = n1
				n4.parent = n1
				n0.right = n3
				n3.parent = n0
				return n0
			},
		},
		{
			name: "child is nil",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newLeftGrand()
				grand.right.left = nil
				return grand, grand
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n5 := &node[int]{data: 5}
				n3.left = n1
				n3.right = n5
				n1.parent = n3
				n5.parent = n3
				n1.left = n2
				n2.parent = n1
				return n3

			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			grand, root := tc.grandandroot()
			r.root = root
			r.leftRotate(grand)
			wanttree := tc.wantTree()
			assert.Equal(t, wanttree, r.root)
		})
	}
}

func newLeftGrand() *node[int] {
	n1 := &node[int]{data: 1}
	n2 := &node[int]{data: 2}
	n3 := &node[int]{data: 3}
	n4 := &node[int]{data: 4}
	n5 := &node[int]{data: 5}
	n1.left = n2
	n1.right = n3
	n2.parent = n1
	n3.parent = n1
	n3.left = n4
	n3.right = n5
	n4.parent = n3
	n5.parent = n3
	return n1
}

func intcompare(src int, dst int) int {
	if src < dst {
		return -1
	}
	if src == dst {
		return 0
	}
	return 1
}

func TestRedBlackTree_RightRotate(t *testing.T) {
	r := NewRedBlackTree[int](intcompare)
	testcases := []struct {
		name         string
		wantTree     func() *node[int]
		grandandroot func() (*node[int], *node[int])
	}{
		{
			/*

					4                  2
				   /  \               /  \
				  2    5  --右旋->   1     4
				 / \                     /  \
				1	3			        3	  5
			*/
			name: "grand is root",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newRightGrand()
				return grand, grand
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n2.left = n1
				n2.right = n4
				n1.parent = n2
				n4.parent = n2
				n4.left = n3
				n4.right = n5
				n3.parent = n4
				n5.parent = n4
				return n2
			},
		},
		{
			/*
							6                 6
				          /                  /
						4                  2
					   /  \               /  \
					  2    5  --右旋->   1     4
					 / \                     /  \
					1	3			        3	  5
			*/
			name: "grand is root left",
			grandandroot: func() (*node[int], *node[int]) {
				root := &node[int]{data: 6}
				grand := newRightGrand()
				root.left = grand
				grand.parent = root
				return grand, root
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n6 := &node[int]{data: 6}
				n2.left = n1
				n2.right = n4
				n1.parent = n2
				n4.parent = n2
				n4.left = n3
				n4.right = n5
				n3.parent = n4
				n5.parent = n4
				n6.left = n2
				n2.parent = n6
				return n6
			},
		},
		{
			/*        0                 0
			            \                 \
						4                  2
					   /  \               /  \
					  2    5  --右旋->   1     4
					 / \                     /  \
					1	3			        3	  5
			*/
			name: "grand is root right",
			grandandroot: func() (*node[int], *node[int]) {
				root := &node[int]{data: 0}
				grand := newRightGrand()
				root.right = grand
				grand.parent = root
				return grand, root
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n3 := &node[int]{data: 3}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n0 := &node[int]{data: 0}
				n2.left = n1
				n2.right = n4
				n1.parent = n2
				n4.parent = n2
				n4.left = n3
				n4.right = n5
				n3.parent = n4
				n5.parent = n4
				n0.right = n2
				n2.parent = n0
				return n0
			},
		},
		{
			/*
				 4                  2
				/  \               /  \
				2    5  --右旋->   1    4
				/                        \
				1				          5
			*/
			name: "no child",
			grandandroot: func() (*node[int], *node[int]) {
				grand := newRightGrand()
				grand.left.right = nil
				return grand, grand
			},
			wantTree: func() *node[int] {
				n1 := &node[int]{data: 1}
				n2 := &node[int]{data: 2}
				n4 := &node[int]{data: 4}
				n5 := &node[int]{data: 5}
				n2.left = n1
				n2.right = n4
				n1.parent = n2
				n4.parent = n2
				n4.right = n5
				n5.parent = n4
				return n2
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			grand, root := tc.grandandroot()
			r.root = root
			r.rightRotate(grand)
			wanttree := tc.wantTree()
			assert.Equal(t, wanttree, r.root)
		})
	}

}

func newRightGrand() *node[int] {
	n1 := &node[int]{data: 1}
	n2 := &node[int]{data: 2}
	n3 := &node[int]{data: 3}
	n4 := &node[int]{data: 4}
	n5 := &node[int]{data: 5}
	n4.left = n2
	n4.right = n5
	n2.parent = n4
	n5.parent = n4
	n2.left = n1
	n2.right = n3
	n1.parent = n2
	n3.parent = n2
	return n4
}

func TestRedBlackTree_get(t *testing.T) {

	testcases := []struct {
		name    string
		key     int
		wantVal int
		isFound bool
	}{
		{
			name:    "key is root",
			key:     55,
			wantVal: 55,
			isFound: true,
		},
		{
			name:    "key is child right",
			key:     38,
			wantVal: 38,
			isFound: true,
		},
		{
			name:    "key is child left",
			key:     80,
			wantVal: 80,
			isFound: true,
		},
		{
			name:    "key is left grandchild",
			key:     25,
			wantVal: 25,
			isFound: true,
		},
		{
			name:    "key is right grandchild",
			key:     46,
			wantVal: 46,
			isFound: true,
		},
		{
			name:    "key not found",
			key:     90,
			isFound: false,
		},
	}
	rbt := NewRedBlackTree[int](intcompare)
	rbt.root = newTestRbTree()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n, ok := rbt.get(rbt.root, tc.key)
			assert.Equal(t, tc.isFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantVal, n.data)

		})
	}
}

func TestRedBlackTreeInsert(t *testing.T) {
	rbt := NewRedBlackTree(intcompare)
	testcases := []struct {
		name     string
		val      int
		Setroot  func() *node[int]
		wantTree func() *node[int]
	}{
		{
			name: "empty tree",
			val:  1,
			Setroot: func() *node[int] {
				return nil
			},
			wantTree: func() *node[int] {
				return &node[int]{
					typ:  black,
					data: 1,
				}
			},
		},
		{
			name: "parent is black",
			val:  40,
			Setroot: func() *node[int] {
				root := newTestRbTree()
				return root
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n46, _ := rbt.get(root, 46)
				n40 := &node[int]{data: 40, typ: red}
				n46.left = n40
				n40.parent = n46
				return root
			},
		},
		{
			/*

					            						    55_b
				                                  /         		   \
									            38_r              	    80_r
							                 /       \       	  	  /      \
						                   25_b       50_b   		76_b     88_b
							              /    \      /   \         /
							            17_r   33_r  46_r   52_r  72_r
			*/
			name: "uncle is black rr",
			val:  52,
			Setroot: func() *node[int] {
				root := newTestRbTree()
				return root
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n46, _ := rbt.get(root, 46)
				n38, _ := rbt.get(root, 38)
				n50, _ := rbt.get(root, 50)
				n52 := &node[int]{data: 52, typ: red}
				n38.right = n50
				n50.parent = n38
				n50.left = n46
				n50.right = n52
				n50.typ = black
				n46.typ = red
				n46.right = nil
				n46.parent = n50
				n52.parent = n50
				return root
			},
		},
		{
			/*
												55_b
						      			/         		   \
							          38_r              	80_r
					                 /   \       	  	  /      \
				                   25_b   46_b   		72_b      88_b
					              /    \     \         /   \
					            17_r   33_r   50_r    60_r    76_r


			*/
			name: "uncle is black ll",
			val:  60,
			Setroot: func() *node[int] {
				root := newTestRbTree()
				return root
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n80, _ := rbt.get(root, 80)
				n76, _ := rbt.get(root, 76)
				n72, _ := rbt.get(root, 72)
				n60 := &node[int]{data: 60, typ: red}
				n80.left = n72
				n72.parent = n80
				n72.left = n60
				n72.right = n76
				n72.typ = black
				n76.left = nil
				n76.parent = n72
				n76.typ = red
				n60.parent = n72
				return root
			},
		},
		{
			/*
					                           55_b
						      			/         		   \
							          38_r              	80_r
					                 /   	\       	  	  /      \
				                   25_b   	48_b   		76_b      88_b
					              /    \   、 /	\         /
					            17_r   33_r  46_r 50_r   72_r
			*/
			name: "uncle is black rl",
			val:  48,
			Setroot: func() *node[int] {
				return newTestRbTree()
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n46, _ := rbt.get(root, 46)
				n50, _ := rbt.get(root, 50)
				n38, _ := rbt.get(root, 38)
				n48 := &node[int]{
					data: 48,
					typ:  red,
				}
				n38.right = n48
				n48.parent = n38
				n48.left = n46
				n48.right = n50
				n48.typ = black
				n46.right = nil
				n46.typ = red
				n46.parent = n48
				n50.parent = n48
				return root
			},
		},
		{
			/*
				                                   55_b
							      			/         		   \
								          38_r              	80_r
						                 /   \       	  	  /      \
					                   25_b   46_b   		74_b      88_b
						              /    \     \         /  \
						            17_r   33_r   50_r   72_r  76_r

			*/
			name: "uncle is black lr",
			val:  74,
			Setroot: func() *node[int] {
				return newTestRbTree()
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n80, _ := rbt.get(root, 80)
				n76, _ := rbt.get(root, 76)
				n72, _ := rbt.get(root, 72)
				n74 := &node[int]{typ: black, data: 74}
				n80.left = n74
				n74.left = n72
				n74.right = n76
				n74.parent = n80
				n76.parent = n74
				n76.typ = red
				n76.left = nil
				n76.right = nil
				n72.parent = n74
				return root
			},
		},
		{
			/*
				                                                   55_b
											      			/         		   \
												          38_b              	80_b
										                 /   \       	  	  /      \
									                   25_r   46_b   		76_b      88_b
										              /    \     \         /
										            17_b   33_b   50_r   72_r
					 								/
								                 10_r
			*/
			name: "uncle is red",
			val:  10,
			Setroot: func() *node[int] {
				return newTestRbTree()
			},
			wantTree: func() *node[int] {
				root := newTestRbTree()
				rbt := NewRedBlackTree[int](intcompare)
				n38, _ := rbt.get(root, 38)
				n80, _ := rbt.get(root, 80)
				n25, _ := rbt.get(root, 25)
				n17, _ := rbt.get(root, 17)
				n33, _ := rbt.get(root, 33)
				n10 := &node[int]{data: 10, typ: red}
				n38.typ = black
				n80.typ = black
				n25.typ = red
				n17.typ = black
				n33.typ = black
				n17.left = n10
				n10.parent = n17
				return root
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rbt.root = tc.Setroot()
			rbt.Insert(tc.val)
			assert.Equal(t, tc.wantTree(), rbt.root)
		})

	}

}

func newTestRbTree() *node[int] {
	/*
					_b: black
					_r: red
			                           55_b
				      			/         		   \
					          38_r              	80_r
			                 /   \       	  	  /      \
		                   25_b   46_b   		76_b      88_b
			              /    \     \         /
			            17_r   33_r   50_r   72_r


	*/

	n55 := &node[int]{data: 55, typ: black}
	n38 := &node[int]{data: 38, typ: red}
	n80 := &node[int]{data: 80, typ: red}
	n25 := &node[int]{data: 25, typ: black}
	n46 := &node[int]{data: 46, typ: black}
	n76 := &node[int]{data: 76, typ: black}
	n88 := &node[int]{data: 88, typ: black}
	n17 := &node[int]{data: 17, typ: red}
	n33 := &node[int]{data: 33, typ: red}
	n50 := &node[int]{data: 50, typ: red}
	n72 := &node[int]{data: 72, typ: red}
	n55.left = n38
	n55.right = n80
	n38.parent = n55
	n80.parent = n55
	n38.left = n25
	n38.right = n46
	n25.parent = n38
	n46.parent = n38
	n80.left = n76
	n80.right = n88
	n76.parent = n80
	n88.parent = n80
	n25.left = n17
	n25.right = n33
	n17.parent = n25
	n33.parent = n25
	n46.right = n50
	n50.parent = n46
	n76.left = n72
	n72.parent = n76

	return n55
}

func TestRedBlackTree_Delete(t *testing.T) {
	testcases := []struct {
		name     string
		val      int
		setRoot  func() *node[int]
		wantTree func() *node[int]
		isFound  bool
	}{
		{
			name: "delete only one element",
			val:  1,
			setRoot: func() *node[int] {
				return &node[int]{
					typ:  black,
					data: 1,
				}
			},
			wantTree: func() *node[int] {
				return nil
			},
			isFound: true,
		},
		{
			name: "delete red node",
			val:  50,
			setRoot: func() *node[int] {
				n50 := &node[int]{data: 50, typ: red}
				n46 := &node[int]{data: 46, typ: black}
				n38 := &node[int]{data: 38, typ: red}
				n46.left = n38
				n46.right = n50
				n38.parent = n46
				n50.parent = n46
				return n46
			},
			wantTree: func() *node[int] {
				n46 := &node[int]{data: 46, typ: black}
				n38 := &node[int]{data: 38, typ: red}
				n38.parent = n46
				n46.left = n38
				return n46
			},
			isFound: true,
		},
		{
			/*
					55_b                      55_b
				   /    \                    /    \
				  46_b   80_b    --->       50_b  80_b
				     \
				      50_r
			*/
			name: "delete node is black,replacement is red",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n50 := &node[int]{data: 50, typ: red}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n46.right = n50
				n50.parent = n46
				return n55
			},
			val: 46,
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n50 := &node[int]{data: 50, typ: black}
				n55.left = n50
				n55.right = n80
				n50.parent = n55
				n80.parent = n55
				return n55
			},
			isFound: true,
		},
		{
			/*
					55_b                      55_b
				   /   \                      /   \
				46_b   80_r      ---->      46_b   78_r
				       /   \                       /  \
				      76_b   88_b                76_b  80_b
				        \
				         78_r
			*/
			name: "delete black leaf and node is right child and sibling is black and sibling has only right red child",
			val:  88,
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n78 := &node[int]{data: 78, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				n76.right = n78
				n78.parent = n76
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n76 := &node[int]{data: 76, typ: black}
				n78 := &node[int]{data: 78, typ: red}
				n55.left = n46
				n55.right = n78
				n46.parent = n55
				n78.parent = n55
				n78.left = n76
				n78.right = n80
				n76.parent = n78
				n80.parent = n78
				return n55
			},
			isFound: true,
		},
		{
			/*
						55_b                      55_b
					   /   \                      /   \
					46_b   80_r      ---->      46_b   76_r
					       /   \                       /  \
					      76_b   88_b                72_b  80_b
				          /
				         72_r
			*/
			name: "black leaf and node is right child and sibling is black and sibling has only left red child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n72 := &node[int]{data: 72, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				n76.left = n72
				n72.parent = n76
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n76 := &node[int]{data: 76, typ: red}
				n72 := &node[int]{data: 72, typ: black}
				n55.left = n46
				n55.right = n76
				n46.parent = n55
				n76.parent = n55
				n76.left = n72
				n76.right = n80
				n72.parent = n76
				n80.parent = n76
				return n55
			},
			isFound: true,
			val:     88,
		},
		{
			/*
						55_b                      55_b
					   /   \                      /   \
					46_b   80_r      ---->      46_b   76_r
					       /   \                       /  \
					      76_b   88_b                72_b  80_b
				          /  \                              /
				         72_r  78_r                        78_r
			*/
			name: "black leaf and node is right child and sibling is black and sibling has two red child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n72 := &node[int]{data: 72, typ: red}
				n78 := &node[int]{data: 78, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				n76.left = n72
				n76.right = n78
				n72.parent = n76
				n78.parent = n76
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n76 := &node[int]{data: 76, typ: red}
				n72 := &node[int]{data: 72, typ: black}
				n78 := &node[int]{data: 78, typ: red}
				n55.left = n46
				n55.right = n76
				n46.parent = n55
				n76.parent = n55
				n76.left = n72
				n76.right = n80
				n72.parent = n76
				n80.parent = n76
				n80.left = n78
				n78.parent = n80
				return n55
			},
			isFound: true,
			val:     88,
		},
		{
			/*
					55_b           			 55_b
				  /		\       		   /     \
				46_b	80_r      --->	 46_b    80_b
						/   \           	    /
				       76_b  88_b             76_r
			*/
			name: "black leaf node and node is right child and sibling is black and sibling has no child and parent is red",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n76 := &node[int]{data: 76, typ: red}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n76.parent = n80
				return n55
			},
			isFound: true,
			val:     88,
		},
		{
			/*
						55_b               55_b
						/   \             /
				      46_b	 80_b  --->  46_r

			*/
			name: "black leaf node and node is right child and sibling is black and sibling has no child and parent is black",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: red}
				n55.left = n46
				n46.parent = n55
				return n55
			},
			isFound: true,
			val:     80,
		},
		{
			/*
					80_b             55_b
				   /   \            /   \
				  55_r  88_b  -->  46_b   80_b
				 /  \                     /
				46_b 76_b                76_r
			*/
			name: "black leaf node and node is right child and sibling is red",
			setRoot: func() *node[int] {
				n80 := &node[int]{data: 80, typ: black}
				n55 := &node[int]{data: 55, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n76 := &node[int]{data: 76, typ: black}
				n80.left = n55
				n80.right = n88
				n55.parent = n80
				n88.parent = n80
				n55.left = n46
				n55.right = n76
				n46.parent = n55
				n76.parent = n55
				return n80
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n76 := &node[int]{data: 76, typ: red}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n76.parent = n80
				return n55
			},
			isFound: true,
			val:     88,
		},
		{
			/*
						55_b                      55_b
					   /   \                      /   \
					46_b   80_r      ---->      46_b   81_r
					       /   \                       /  \
					      76_b   88_b                80_b  88_b
					             /
				                81_r
			*/
			name: "black leaf node and node is left child and sibling is black and sibling has only red left child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n81 := &node[int]{data: 81, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				n88.left = n81
				n81.parent = n88
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n81 := &node[int]{data: 81, typ: red}
				n80 := &node[int]{data: 80, typ: black}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n81
				n46.parent = n55
				n81.parent = n55
				n81.left = n80
				n81.right = n88
				n80.parent = n81
				n88.parent = n81
				return n55
			},
			isFound: true,
			val:     76,
		},
		{
			/*
						55_b                      55_b
					   /   \                      /   \
					46_b   80_r      ---->      46_b   88_r
					       /   \                       /  \
					      76_b   88_b                80_b  90_b
					                \
				                    90_r
			*/
			name: "black leaf node and node is left child and sibling is black and sibling has only one right red child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n90 := &node[int]{data: 90, typ: red}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				n88.right = n90
				n90.parent = n88
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n90 := &node[int]{data: 90, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n88 := &node[int]{data: 88, typ: red}
				n55.left = n46
				n55.right = n88
				n46.parent = n55
				n88.parent = n55
				n88.left = n80
				n88.right = n90
				n80.parent = n88
				n90.parent = n88
				return n55
			},
			isFound: true,
			val:     76,
		},
		{
			/*
					55_b           			 55_b
				  /		\       		   /     \
				46_b	80_r      --->	 46_b    80_b
						/   \           	    	\
				       76_b  88_b                   88_r
			*/
			name: "black leaf node and node is left child and sibling is black and sibling has no child and parent is red",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n88.parent = n80
				n76.parent = n80
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n88 := &node[int]{data: 88, typ: red}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.right = n88
				n88.parent = n80
				return n55
			},
			isFound: true,
			val:     76,
		},
		{
			/*
						55_b               55_b
						/   \            		\
				      46_b	 80_b  --->        80_r

			*/
			name: "black leaf node and node is left child and sibling is black and sibling has no child and parent is black",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n55.right = n80
				n80.parent = n55
				return n55
			},
			isFound: true,
			val:     46,
		},
		{
			/*
					80_b           			 88_b
				  /		\       		   /     \
				55_b	88_r      --->	 88_b    90_b
						/   \              \
				       81_b  90_b           81_r
			*/
			name: "black leaf node and node is left child and sibling is red ",
			setRoot: func() *node[int] {
				n80 := &node[int]{data: 80, typ: black}
				n88 := &node[int]{data: 88, typ: red}
				n55 := &node[int]{data: 55, typ: black}
				n81 := &node[int]{data: 81, typ: black}
				n90 := &node[int]{data: 90, typ: black}
				n80.left = n55
				n80.right = n88
				n55.parent = n80
				n88.parent = n80
				n88.left = n81
				n88.right = n90
				n81.parent = n88
				n90.parent = n88
				return n80
			},
			wantTree: func() *node[int] {
				n88 := &node[int]{data: 88, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n90 := &node[int]{data: 90, typ: black}
				n81 := &node[int]{data: 81, typ: red}
				n88.left = n80
				n88.right = n90
				n80.parent = n88
				n90.parent = n88
				n80.right = n81
				n81.parent = n80
				return n88
			},
			isFound: true,
			val:     55,
		},
		{
			/*
					55_b                      55_b
				   /   \                      /   \
				46_b   80_r      ---->      46_b   88_b
				       /   \                       /
				      76_b   88_b                76_r

			*/
			name: "delete node with two child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: red}
				n76 := &node[int]{data: 76, typ: black}
				n88 := &node[int]{data: 88, typ: black}
				n55.left = n46
				n55.right = n80
				n46.parent = n55
				n80.parent = n55
				n80.left = n76
				n80.right = n88
				n76.parent = n80
				n88.parent = n80
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n88 := &node[int]{data: 88, typ: black}
				n76 := &node[int]{data: 76, typ: red}
				n55.left = n46
				n55.right = n88
				n46.parent = n55
				n88.parent = n55
				n88.left = n76
				n76.parent = n88
				return n55
			},
			isFound: true,
			val:     80,
		},
		{
			name: "val not found",
			val:  11,
			setRoot: func() *node[int] {
				return nil
			},
			wantTree: func() *node[int] {
				return nil
			},
			isFound: false,
		},
		{
			/*
					55_b           		  55_b
					/  \       --->     /    \
				   46_b  80_b      	43_b    80_b
				   /
				  43_r
			*/
			name: "delete black node with left red child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n46 := &node[int]{data: 46, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n43 := &node[int]{data: 43, typ: red}
				n55.left = n46
				n55.right = n80
				n80.parent = n55
				n46.parent = n55
				n46.left = n43
				n43.parent = n46
				return n55
			},
			wantTree: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n43 := &node[int]{data: 43, typ: black}
				n80 := &node[int]{data: 80, typ: black}
				n55.left = n43
				n55.right = n80
				n43.parent = n55
				n80.parent = n55
				return n55
			},
			val:     46,
			isFound: true,
		},
		{
			/*
			   55_b           43_b
			   /     --->
			  43_r
			*/
			name: "delete root and root has right child",
			setRoot: func() *node[int] {
				n55 := &node[int]{data: 55, typ: black}
				n43 := &node[int]{data: 43, typ: red}
				n55.left = n43
				n43.parent = n55
				return n55
			},
			wantTree: func() *node[int] {
				return &node[int]{data: 43, typ: black}
			},
			isFound: true,
			val:     55,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rbt := NewRedBlackTree[int](intcompare)
			rbt.root = tc.setRoot()
			ok := rbt.Delete(tc.val)
			assert.Equal(t, tc.isFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantTree(), rbt.root)
		})
	}
}
