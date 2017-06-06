package main
import "fmt"
import "os/exec"
import "strconv"
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
    fmt.Println(hunks)
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

    for i := range lines {
        isEndOfHunk := i == len(lines) - 1 || strings.HasPrefix(lines[i+1], "@@")

        if isEndOfHunk {
            if hunkStart > 0 {
                hunkEnd := i + 1

                if len(lines[i]) == 0 {
                    hunkEnd -= 1
                }

                hunk := strings.Join(lines[hunkStart:hunkEnd], linebreak)

                // Append hunk
                hunks = append(hunks, hunk)
            }
            hunkStart = i + 1
        }
    }

    return hunks
}

func ParseHunk(hunk string) Hunk {
    x := strings.Split(hunk, " ")

    a := strings.Split(x[1][1:], ",")
    b := strings.Split(x[2][1:], ",")

    body := hunk[strings.Index(hunk, "\n") + 1:]

    aStart, aStartErr := strconv.Atoi(a[0])
    aLen,   aLenErr   := strconv.Atoi(a[1])
    bStart, bStartErr := strconv.Atoi(b[0])
    bLen,   bLenErr   := strconv.Atoi(b[1])

    if aStartErr != nil || aLenErr != nil || bStartErr != nil || bLenErr != nil {

    }

    dels := ParseDels(aStart, body)
    adds := ParseAdds(bStart, body)

    return Hunk{
        Range{aStart, aLen},
        Range{bStart, bLen},
        dels,
        adds,
        body,
    }
}

func ParseAdds(start int, body string) []int {
    adds := make([]int, 0)

    n := start
    for _, line := range strings.Split(body, "\n") {
        if strings.HasPrefix(line, " ") {
            n++
        } else if strings.HasPrefix(line, "+") {
            adds = append(adds, n)
            n++
        }
    }

    return adds
}

func ParseDels(start int, body string) []int {
    dels := make([]int, 0)

    n := start
    for _, line := range strings.Split(body, "\n") {
        if strings.HasPrefix(line, " ") {
            n++
        } else if strings.HasPrefix(line, "-") {
            dels = append(dels, n)
            n++
        }
    }

    return dels
}

type Hunk struct {
    Lhs Range
    Rhs Range
    Adds []int // Line numbers on lhs that are removed
    Dels []int // Line numbers on rhs that are new
    Body string
}

type Range struct {
    Offset int
    Length int
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
