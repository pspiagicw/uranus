package repl

import "testing"

func TestExample(t *testing.T) {
    if true != true {
        t.Errorf("true is not equal to true")
    }
}
