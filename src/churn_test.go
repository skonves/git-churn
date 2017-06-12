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

func TestParseHunk(t *testing.T) {
    // ARRANGE
    hunkStr := "@@ -7,4 +1,5 @@ location\n nc1\n-del1\n-del2\n+add1\n nc2\n+add2\n nc3"

    expectedHunk := Hunk{ Range{7,4}, Range{1,5}, []int{8,9}, []int{2,4}, " nc1\n-del1\n-del2\n+add1\n nc2\n+add2\n nc3"}

    // ACT
    result := ParseHunk(hunkStr)

    // ASSERT
    assert.Equal(t, expectedHunk, result, "fail!")
}

func TestParseAdds(t *testing.T) {
    // ARRANGE
    start := 1
    body  := " nc1\n-del1\n-del2\n+add1\n nc2\n+add2\n nc3"

    expectedAdds := []int{2, 4}

    // ACT
    result := ParseAdds(start, body)

    // ASSERT
    assert.Equal(t, expectedAdds, result, "fail!")
}

func TestParseDels(t *testing.T) {
    // ARRANGE
    start := 7
    body  := " nc1\n-del1\n-del2\n+add1\n nc2\n+add2\n nc3"

    expectedDels := []int{8, 9}

    // ACT
    result := ParseDels(start, body)

    // ASSERT
    assert.Equal(t, expectedDels, result, "fail!")
}

func TestApplyDels(t *testing.T) {
    // ARRANGE
    matrix := [][]bool{
        []bool{ true,  true  },
        []bool{ false, false },
        []bool{ true,  true  },
        []bool{ true,  true  },
        []bool{ true,  true  },
    }

    dels := []int{1, 3}

    expectedResult := [][]bool{
        []bool{ true,  false }, // line 1 - Removed
        []bool{ false, false },
        []bool{ true,  true  }, // line 2
        []bool{ true,  false }, // line 3 - Removed
        []bool{ true,  true  }, // line 4
    }

    // ACT
    result := ApplyDels(matrix, dels)

    // ASSERT
    assert.Equal(t, expectedResult, result, "fail")
}

func TestApplyAdds(t *testing.T) {
    // ARRANGE
    matrix := [][]bool{
        []bool{ true,  false },
        []bool{ false, false },
        []bool{ true,  true  },
        []bool{ true,  false },
        []bool{ true,  true  },
    }

    adds := []int{1, 3}

    expectedResult := [][]bool{
        []bool{ true,  false },
        []bool{ false, false },
        []bool{ false, true  }, // Added
        []bool{ true,  true  },
        []bool{ true,  false },
        []bool{ false, true  }, // Added
        []bool{ true,  true  },
    }

    // ACT
    result := ApplyAdds(matrix, adds)

    // ASSERT
    assert.Equal(t, expectedResult, result, "fail")
}
