package lib

import (
  "encoding/json"
  "errors"
  "fmt"
  "hash"
)

type ToTree struct {
  Content   string    `json:"content"`
  Depth     int       `json:"depth,string"`
  LineStart int       `json:"lineStart,string"` // up to 2.4 billion lines
  NumLines  int       `json:"numLines,string"`  // how many were part of this one
  State     []byte    `json:"state"`
  parent    *ToTree   // can only have one root
  Branches  []*ToTree `json:"branches,omitempty"` // can have many branches
  Hash      hash.Hash `json:"hash,string"`
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

func PrintTree(tree *ToTree) {
  k := len(tree.Branches)
  n := 0
  for n < k {
    s := fmt.Sprintf("%%d %s %%%ds\n", tree.Branches[n].State, tree.Branches[n].Depth+len(tree.Branches[n].Content))
    fmt.Printf(s, tree.Branches[n].LineStart, tree.Branches[n].Content)
    PrintTree(tree.Branches[n])
    n += 1
  }
}

func PrintJSONTree(tree *ToTree) {
  jsonTree, err := json.Marshal(*tree)
  if err != nil {
    fmt.Printf("Sooorry: %s\n", err)
  }
  fmt.Println(string(jsonTree))
}

func newRoot() *ToTree {
  root := new(ToTree)
  root.Depth = 0
  root.parent = nil
  root.Content = string("")
  return root
}

func newLeaf(depth int, state []byte, content []byte, lineStart int) (leaf *ToTree) {
  leaf = new(ToTree)
  leaf.Depth = depth
  leaf.State = state[:]
  leaf.Content = string(content)
  leaf.LineStart = lineStart + 1
  leaf.NumLines = 1
  return
}

func addLeaf(parent, leaf *ToTree) {
  parent.Branches = append(parent.Branches, leaf)
}

func addBranch(parent, branch *ToTree, trackingFile *ToFile, startLine int) (i int) {
  addLeaf(parent, branch)
  i = startLine
  for i < trackingFile.numLines {
    depth, state, content := ReadItem(&trackingFile.lines[i])

    // Four possibilties for each new entry:
    // 1) This is a continuation of the previous entry, depth does not matter
    // 2) This is a sibling entry: new depth == old depth
    // 3) This is a new branch: new depth > old depth
    // 4) This is the end of this particular branch, new depth < old depth

    // 1) continuation
    if UnknownState(state) && Symbol(state[0]) != Continue {
      // This is an assumed continuation of the previous line
      // Because it is either an explicit continuation '\' or
      // it starts with a character that is not a recognized explicit state
      branch.Content += " " + string(content)
      branch.NumLines++
      i++
    } else if depth == branch.Depth {
      branch = newLeaf(depth, state, content, i)
      addLeaf(parent, branch)
      i++
    } else if depth > branch.Depth {
      i = addBranch(branch, newLeaf(depth, state, content, i), trackingFile, i+1)
    } else if depth < branch.Depth {
      return
    }
  }
  return
}

func ParseTracking(trackingFile *ToFile) (root *ToTree, err error) {
  root = newRoot()

  if trackingFile.numLines == 0 {
    addBranch(newRoot(), nil, trackingFile, 0)
    return root, errors.New("Empty tracking file")
  }

  // Handle the first line of the file
  depth, state, content := ReadItem(&trackingFile.lines[0])

  if depth > 0 {
    addBranch(newRoot(), nil, trackingFile, 0)
    return root, errors.New("Tracking file did not start at 0 depth")
  }

  addBranch(root, newLeaf(depth, state, content, 0), trackingFile, 1)
  return root, nil
}
