package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestGetHunks(t *testing.T) {

    // ARRANGEÏ
    tables := []struct {
        condition string;        diff string
    } {
        { "trailing newline",    "aaa\n@@ -1,1 +1,1 @@ bbb\n ccc\n@@ -2,2 +2,2 @@ ddd\n eee\n" }, 
        { "no trailing newline", "aaa\n@@ -1,1 +1,1 @@ bbb\n ccc\n@@ -2,2 +2,2 @@ ddd\n eee" },
    }

    expectedHunks := []string {
        "@@ -1,1 +1,1 @@ bbb\n ccc",
        "@@ -2,2 +2,2 @@ ddd\n eee",
    }

    for _, table := range tables {
        // ACT
        hunks := GetHunks(table.diff)

        // ASSERT
        assert.Equal(t, len(expectedHunks), len(hunks), table.condition + ": wrong number of elements")

        for i, hunk := range hunks {
            assert.Equal(t, expectedHunks[i], hunk, table.condition + ": wrong element")
        }
    }
}