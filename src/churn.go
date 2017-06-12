package main
import "fmt"
import "os/exec"
import "strconv"
import "strings"

func main() {
    filename := "src/churn.go"
    shas := GetShas(filename)
    diff := GetDiff(filename, shas[0], shas[1])

    ct := GetFileLineCount(filename, shas[0])
    matrix := MakeMatrix(ct)
    matrix = DupLastRow(matrix)

    dels := []int{}
    adds := []int{}
    for _, hunk := range GetHunks(diff) {
        hunk := ParseHunk(hunk)
        dels = append(dels, hunk.Dels...)
        adds = append(adds, hunk.Adds...)
    }

    matrix = ApplyDels(matrix, dels)

    fmt.Println(matrix)
    // matrix = AddDiff(matrix, diff)
    // fmt.Println(matrix)

    // hunks := GetHunks(diff)

    // for _, h := range hunks {
    //     p := ParseHunk(h)
    //     fmt.Println(p)
    // }

    //fmt.Println(hunks)
}

func MakeMatrix(initialLineCount int) [][]bool {
    matrix := make([][]bool, initialLineCount)

    for i := range matrix {
        matrix[i] = []bool{true}
    }

    return matrix
}

func DupLastRow(matrix [][]bool) [][]bool {
    width := len(matrix[0])
    for i := range matrix {
        matrix[i] = append(matrix[i], matrix[i][width-1])
    }
    return matrix
}

func ApplyDels(matrix [][]bool, dels []int) [][]bool {
    lhscol := len(matrix[0]) - 2
    rhscol := lhscol + 1
    lhsrow := 1

    for i, row := range matrix {
        if row[lhscol] == true {
            for _, del := range dels {
                if lhsrow == del {
                    matrix[i][rhscol] = false
                }
            }
            lhsrow += 1
        }
    }

    return matrix
}

func ApplyAdds(matrix [][]bool, adds []int) [][]bool {
    lhscol := len(matrix[0]) - 2
    rhscol := lhscol + 1
    rhsrow := 1

    width := len(matrix[0])

    for i := 0; i < len(matrix); i++ {
        //fmt.Println(i, rhsrow)
        row := matrix[i]
        if row[rhscol] == true {
            for _, add := range adds {
                if rhsrow == add {
                    // Insert row

                    addedRow :=  make([]bool, width)
                    addedRow[rhscol] = true

                    prior := matrix[:i]
                    after := matrix[i:]

                    fmt.Println(prior, " <= ", len(prior), " : ", addedRow, " : ", len(after), " => ", after)

                    matrix = append(prior, append([][]bool{addedRow}, after...)...)

                    fmt.Println(matrix)
                    fmt.Println("")
                    //rhsrow += 1

                    //matrix[i][rhscol] = false
                }
            }
            rhsrow += 1
        }
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
    Adds []int // Line numbers on rhs that are new
    Dels []int // Line numbers on lhs that are removed
    Body string
}

type Range struct {
    Offset int
    Length int
}

func Insert(matrix [][]bool, record []bool, i int) [][]bool {
    return append(append(matrix[:i], record), matrix[i+1:]...)
}

func GetShas(filename string) []string { // TODO: add start and end date params
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
