package main
import "fmt"
import "strconv"
import "os/exec"
import "strings"

func main() {
    filename := "src/exec.go"
    shas := GetShas(filename)
    diff := GetDiff(filename, shas[0], shas[1])

    // ct := GetFileLineCount(filename, shas[0])
    // matrix := MakeMatrix(ct)
    // matrix = AddDiff(matrix, diff)
    // fmt.Println(matrix)

    hunks := GetHunks(diff)
    s := GetHunkStart(hunks[0])
    fmt.Println(s)
    // fmt.Println()
    // fmt.Println(hunks)
}

func MakeMatrix(initialLineCount int) [][]bool {
    matrix := make([][]bool, initialLineCount)

    for i := 0; i < len(matrix); i++ {
        matrix[i] = []bool{true}
    }

    return matrix
}

func AddDiff(matrix [][]bool, diff string) [][]bool {
    // Add new column
    width := len(matrix[0])
    for i := 0; i < len(matrix); i++ {
        matrix[i] = append(matrix[i], matrix[i][width-1])
    }

    // Proc hunks from diff
    hunks := GetHunks(diff)
    for _, hunk := range hunks {
        matrix = AddHunk(matrix, hunk)
    }

    return matrix
}

func GetChanges(hunks []string) map[int]string {
    res := make(map[int]string);



    return res
}

func GetHunkStart(hunk string) int {
    x := strings.Split(hunk, ",")
    n, _ := strconv.Atoi(x[0][4:])
    return n
}

func AddHunk(matrix [][]bool, hunk string) [][]bool {
    return matrix
}

func GetDiff(filename, startSha, endSha string) string {
    cmd := exec.Command("git", "diff", startSha + ".." + endSha, "--", filename)
    stdout, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return "ERROR"
    }

    return string(stdout)
}

func GetFileLineCount(filename, sha string) int {
    cmd := exec.Command("git", "show", sha + ":" + filename)
    stdout, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return 0
    }

    file := string(stdout)

    lines := strings.Split(file, "\n")

    return len(lines)
}

func GetHunks(diff string) []string {
    hunks := make([]string, 0)
    linebreak := "\n"

    lines := strings.Split(diff, linebreak)

    hunkStart := 0
    hunkEnd := 0

    for i, line := range lines {
        isStartOfHunk := strings.HasPrefix(line, "@@")
        isEndOfDiff := i == len(lines) - 1

        if isStartOfHunk || isEndOfDiff {
            if hunkStart > 0 {
                // Build hunk
                if isEndOfDiff {
                    hunkEnd = i
                } else {
                    hunkEnd = i - 1
                }

                hunk := strings.Join(lines[hunkStart:hunkEnd], linebreak)

                // Append hunk
                hunks = append(hunks, hunk)
            }
            hunkStart = i
        }
    }

    return hunks
}

func Insert(matrix [][]bool, record []bool, i int) [][]bool {
    return append(append(matrix[:i], record), matrix[i+1:]...)
}

func GetShas(filename string) []string {
    cmd := exec.Command("git", "log", "--pretty=format:%h", "--", filename)
    stdout, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return []string{}
    }

    lines := strings.Split(string(stdout), "\n")

    // reverse it
    res := make([]string, len(lines))
    for i, v := range lines {
        res[len(lines) - i - 1] = strings.TrimSpace(v)
    }

    return res
}
