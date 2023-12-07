package source

import (
	"testing"
)

func TestNewMerkleTree(t *testing.T) {
	type args struct {
		data [][]byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "123",
			args: args{
				data: [][]byte{
					[]byte("node1"),
					[]byte("node2"),
					[]byte("node3"),
					[]byte("node4"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMerkleTree(tt.args.data); got != nil {
				ShowMerkleTree(got.RootNode)
			}
		})
	}
}
