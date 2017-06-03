package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestAddDiff(t *testing.T) {
    // ARRANGE
    diff :=
`asdfasdf
@@ -1,1 +1,1 @@ package main
 asdf

@@ -1,1 +1,1 @@ package main
 asdf
`

    expectedHunks := []string {
`@@ -1,1 +1,1 @@ package main
 asdf`,
 `@@ -1,1 +1,1 @@ package main
 asdf` }

    // ACT
    hunks := GetHunks(diff)

    // ASSERT
    assert.Equal(t, len(hunks), len(expectedHunks), "wrong number of elements")

    // for i, hunk := range hunks {
    //     t.Log(hunk)
    //     if hunk != expectedHunks[i] {
    //         t.Error("wrong element")
    //     }
    // }
}