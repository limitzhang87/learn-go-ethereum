package source

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// MerkleNode 默克尔树节点结构
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree 默克尔树结构
type MerkleTree struct {
	RootNode *MerkleNode
}

// NewMerkleNode 创建节点
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	// 如果left或right为空，代表其对应data就是最原始的数据节点
	if left == nil || right == nil {
		// 计算hash
		hash := sha256.Sum256(data)
		// 将[32]byte转为[]byte
		mNode.Data = hash[:]
	} else {
		// 将最有子树的数据集合在一起
		prevHash := append(left.Data, right.Data...)
		// 计算hash
		hash := sha256.Sum256(prevHash)
		mNode.Data = hash[:]
	}

	// 左右子树复制
	mNode.Left = left
	mNode.Right = right
	return &mNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode
	// 确保必须为2的整数倍节点
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}
	// 两层循环完成节点树形构造
	for i := 0; i < len(data)/2; i++ {
		// i相当于默尔克数的层数
		newLevel := make([]*MerkleNode, 0, len(nodes)/2)
		for j := 0; j < len(nodes); j += 2 {
			newLevel = append(newLevel, NewMerkleNode(nodes[j], nodes[j+1], nil))
		}
		nodes = newLevel
	}
	// 最后nodes只剩一个，就是根节点
	return &MerkleTree{nodes[0]}
}

func ShowMerkleTree(root *MerkleNode) {
	if root == nil {
		return
	}
	PrintNode(root)
	ShowMerkleTree(root.Left)
	ShowMerkleTree(root.Right)
}

func PrintNode(node *MerkleNode) {
	fmt.Printf("%p\n", node)
	if node != nil {
		fmt.Printf("left[%p], right[%p], data(%x)\n", node.Left, node.Right, node.Data)
		fmt.Printf("check:%t \n", CheckNode(node))
	}
}

func CheckNode(node *MerkleNode) bool {
	if node.Left == nil {
		return false
	}

	prevHash := append(node.Left.Data, node.Right.Data...)
	hash32 := sha256.Sum256(prevHash)
	hash := hash32[:]
	return bytes.Compare(node.Data, hash) == 0
}
