package version

import "github.com/blang/semver/v4"

// Satisfies returns true if version v matches semver constraint expression c.
// Supported operators: = (default), >=, <=, >, <
// Example constraints: "<=1.2.0", ">1.0.0", "1.2.3" (exact)
func Satisfies(constraint, v string) bool {
    if constraint == "" {
        return true
    }
    var op rune
    if len(constraint) > 1 {
        switch constraint[:2] {
        case ">=":
            op = 'g' // greater or equal
            constraint = constraint[2:]
        case "<=":
            op = 'l' // less or equal
            constraint = constraint[2:]
        }
    }
    // single-char ops
    if op == 0 {
        switch constraint[0] {
        case '>':
            op = 'G'
            constraint = constraint[1:]
        case '<':
            op = 'L'
            constraint = constraint[1:]
        case '=':
            op = 'E'
            constraint = constraint[1:]
        }
    }

    verC, err1 := semver.ParseTolerant(constraint)
    verV, err2 := semver.ParseTolerant(v)
    if err1 != nil || err2 != nil {
        return false
    }

    switch op {
    case 'g':
        return verV.GTE(verC)
    case 'l':
        return verV.LTE(verC)
    case 'G':
        return verV.GT(verC)
    case 'L':
        return verV.LT(verC)
    default: // exact match
        return verV.Equals(verC)
    }
}
