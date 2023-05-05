package directory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer, err := NewAnalyzer([]string{"1/2/3", "1/2/3/4"})

	assert.Nil(t, err)
	assert.Len(t, analyzer.childs, 2)
	assert.Equal(t, defaultSearchedFileName, analyzer.searchedFileName)
}

func TestNewAnalyzerWithOtherSearchedFileName(t *testing.T) {
	analyzer, err := NewAnalyzer([]string{"1/2/3", "1/2/3/4"}, WithSearchedFileName(".other"))

	assert.Nil(t, err)
	assert.Len(t, analyzer.childs, 2)
	assert.Equal(t, ".other", analyzer.searchedFileName)
}

func TestNewAnalyzer_should_forwad_error(t *testing.T) {
	_, err := NewAnalyzer([]string{"/1/2/3"})
	assert.NotNil(t, err)
}
