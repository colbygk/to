package lib

import (
  "bytes"
  "testing"

  "cgk.sh/to/lib"
)

func TestSymbols(t *testing.T) {
  if lib.UnknownState([]byte(`-`)) {
    t.Errorf(`unknownState did not recognize -`)
  }
  if lib.UnknownState([]byte(`/`)) {
    t.Errorf(`unknownState did not recognize /`)
  }
  if lib.UnknownState([]byte(`x`)) {
    t.Errorf(`unknownState did not recognize x`)
  }
  if lib.UnknownState([]byte(`r`)) {
    t.Errorf(`unknownState did not recognize r`)
  }
  if lib.UnknownState([]byte(`n`)) {
    t.Errorf(`unknownState did not recognize n`)
  }
  if lib.UnknownState([]byte(`!`)) {
    t.Errorf(`unknownState did not recognize !`)
  }
  if lib.UnknownState([]byte(`\\`)) {
    t.Errorf(`unknownState did not recognize \\`)
  }
  if !lib.UnknownState([]byte(`A`)) {
    t.Errorf(`Did not recognize unknownState`)
  }
}

func TestReadItem(t *testing.T) {
  var d int
  var l []byte
  var c []byte

  depth2 := []byte("  - depth 2")
  d, l, c = lib.ReadItem(&depth2)
  if d != 2 || l[0] != '-' || bytes.Compare(c, []byte("depth 2")) != 0 {
    t.Errorf("depth 2 did not pass: %d %q %q", d, l, c)
  }

  depth0 := []byte("- depth 0")
  d, l, c = lib.ReadItem(&depth0)
  if d != 0 || l[0] != '-' || bytes.Compare(c, []byte("depth 0")) != 0 {
    t.Errorf("depth0 did not pass: %d %q %q", d, l, c)
  }

  notes5 := []byte("     n notes depth 5")
  d, l, c = lib.ReadItem(&notes5)
  if d != 5 || l[0] != 'n' || bytes.Compare(c, []byte("notes depth 5")) != 0 {
    t.Errorf("notes5 did not pass: %d %q %q", d, l, c)
  }

  delete0 := []byte("x delete 0")
  d, l, c = lib.ReadItem(&delete0)
  if d != 0 || l[0] != 'x' || bytes.Compare(c, []byte("delete 0")) != 0 {
    t.Errorf("delete0 did not pass: %d %q %q", d, l, c)
  }

  important0 := []byte("! this is very important!")
  d, l, c = lib.ReadItem(&important0)
  if d != 0 || l[0] != '!' || bytes.Compare(c,
    []byte("this is very important!")) != 0 {
    t.Errorf("important0 did not pass: %d %q %q", d, l, c)
  }

  continue0 := []byte("this is a zero continuation")
  d, l, c = lib.ReadItem(&continue0)
  if d != 0 || l != nil || bytes.Compare(c,
    []byte("this is a zero continuation")) != 0 {
    t.Errorf("continue0 did not pass: %d %q %q", d, l, c)
  }

  continue6 := []byte("      this is a continuation starting on char 6")
  d, l, c = lib.ReadItem(&continue6)
  if d != 6 || l != nil || bytes.Compare(c,
    []byte("this is a continuation starting on char 6")) != 0 {
    t.Errorf("continue6 did not pass: %d %q %q", d, l, c)
  }

  literalcontinue2 := []byte(`  \ this is a literal 2 continuation`)
  d, l, c = lib.ReadItem(&literalcontinue2)
  // This comparison requires the decimal representation of the backslash
  if d != 2 || l[0] != 92 || bytes.Compare(c,
    []byte("this is a literal 2 continuation")) != 0 {
    t.Errorf("literalcontinue2 did not pass: %d %d %q", d, l, c)
  }

  rename10 := []byte(`          r rename`)
  d, l, c = lib.ReadItem(&rename10)
  if d != 10 || l[0] != 'r' || bytes.Compare(c, []byte("rename")) != 0 {
    t.Errorf("rename10 did not pass: %d %d %q", d, l, c)
  }

  completed20 := []byte(`                    / completed`)
  d, l, c = lib.ReadItem(&completed20)
  if d != 20 || l[0] != '/' || bytes.Compare(c, []byte("completed")) != 0 {
    t.Errorf("completed20 did not pass: %d %d %q", d, l, c)
  }

}
