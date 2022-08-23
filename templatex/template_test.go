package templatex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanner_ScanPaths(t *testing.T) {
	scanner := NewScanner()

	scanner.ScanPaths(context.Background(), "test_data")
	txt, err := scanner.Parse(context.Background(), "1", "txt")
	assert.Nil(t, err)
	assert.Equal(t, "test 1 txt", txt)

	txt, err = scanner.Parse(context.Background(), "no.txt", "txt")
	assert.Equal(t, "", txt)
	assert.NotNil(t, err)

	txt, err = scanner.Parse(context.Background(), "3", "txt")
	assert.Nil(t, err)
	assert.Equal(t, "test 3 txt. import test 2 txt", txt)

	txt, err = scanner.Parse(context.Background(), "4", "txt")
	assert.Nil(t, err, err)
	assert.Equal(t, "test 4\n\n    def 111\n\n\n    def 222\n", txt)

	txt, err = scanner.Parse(context.Background(), "sql-template-root", "txt")
	println("sql-template-root:", txt)
}

func Test_defineRename(t *testing.T) {
	raw := `{{define "def"}}def 111{{end}}`
	txt := defineRename(context.Background(), "test", raw)
	assert.Equal(t, `{{define "test/def"}}def 111{{end}}`, txt)
}
