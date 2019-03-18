package lib

import (
	"bytes"
	"fmt"
	"hash"
)

type ToTree struct {
	content string
	depth   int
	lineNo  int       // up to 2.4 billion lines
	state   []byte
	parent  *ToTree   // can only have one root
	branches  []*ToTree // can have many leaves
	hash    hash.Hash
}

func min(a, b int) int {
    if a > b {
        return b
    }
    return a
}

func max(a, b int) int {
    if a < b {
        return b
    }
    return a
}



func newLeaf(leaves []*ToTree, depth int, state []byte, content []byte, lineNo int) ([]*ToTree, *ToTree){
	leaf := new(ToTree)
	leaf.depth = depth
	leaf.state = state
	leaf.content = string(content)
	leaf.lineNo = lineNo
	return append(leaves, leaf), leaf
}

func addBranch(parent, branch *ToTree, trackingFile *ToFile, startLine int) (i int) {
	i = startLine
	branch.parent = parent

	noItems := len(trackingFile.lines)
    for i < noItems {

		depth, state, content := ReadItem(&trackingFile.lines[i])
		if unknownState(state) {
			if len(branch.branches) > 0 {
				last := branch.branches[len(branch.branches)-1]
				last.content += string(state) + string(content)
			} else {
				// ignoring this mystery sauce
			}
			i += 1
		} else {
			// Leaf
			if depth > branch.depth {
				var leaf *ToTree
				branch.branches, leaf = newLeaf(branch.branches, depth, state, content, i)
				i = addBranch(branch, leaf, trackingFile, i + 1)
			}
			// Sibling
			if depth == branch.depth {
				branch.branches, _ = newLeaf(branch.branches, depth, state, content, i)
				i += 1
			}
			// End of tree
			if depth < branch.depth {
				return
			}
		}

        //if false {
        //    currentTree.hash.Write(line)
        //    fmt.Printf("%d %x -> %s\n", i, currentTree.hash.Sum(nil), line)
        //}
    }

	return
}

func PrintTree(tree *ToTree) {
	k := len(tree.branches)
	n := 0
	for n < k {
		s := fmt.Sprintf("%%d - %%%ds\n", tree.branches[n].depth+len(tree.branches[n].content))
		fmt.Printf(s, tree.branches[n].lineNo, tree.branches[n].content)
		PrintTree(tree.branches[n])
		n += 1
	}
}

func newRoot() (*ToTree) {
	root := new(ToTree)
	root.depth = -1
	root.parent = root
	root.content = string("root")
	return root
}

func ParseTracking(trackingFile *ToFile) (root *ToTree){
	trackingFile.lines = bytes.Split(trackingFile.rawData, []byte("\n"))
	root = newRoot()
    addBranch(root, root, trackingFile, 0)
	return
}
