package directory

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/dispatch/pkg/code"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestBed(t *testing.T, infos []code.Info) (string, []string) {
	rootDir := t.TempDir()
	codeOwnerDirs := make([]string, 0)
	for i, info := range infos {
		testBedDir := fmt.Sprintf("%s/%d", rootDir, i)
		codeOwnerDirs = append(codeOwnerDirs, testBedDir)

		if err := os.MkdirAll(testBedDir, 0700); err != nil {
			t.Fatal(err)
		}

		emptyTestBedDir := fmt.Sprintf("%s/%d/empty", rootDir, i)
		codeOwnerDirs = append(codeOwnerDirs, emptyTestBedDir)

		if err := os.MkdirAll(emptyTestBedDir, 0700); err != nil {
			t.Fatal(err)
		}

		content, err := json.Marshal(&info)
		if err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(testBedDir+"/"+".code_owners", content, 0700); err != nil {
			t.Fatal(err)
		}
		codeOwnerDirs = append(codeOwnerDirs, testBedDir)
	}

	return rootDir, codeOwnerDirs
}

func TestNewAnalyzer(t *testing.T) {
	analyzer, err := NewAnalyzer(RootPwdDir, []string{"1/2/3", "1/2/3/4"})

	assert.Nil(t, err)
	assert.Len(t, analyzer.childs, 2)
	assert.Equal(t, defaultSearchedFileName, analyzer.searchedFileName)
}

func TestNewAnalyzer_should_forward_error(t *testing.T) {
	analyzer, err := NewAnalyzer(RootPwdDir, []string{"1/2/3", "1/2/3/4"}, WithSearchedFileName(".other"))

	assert.Nil(t, err)
	assert.Len(t, analyzer.childs, 2)
	assert.Equal(t, ".other", analyzer.searchedFileName)
}

func TestAnalyzer_Analyze(t *testing.T) {
	// given
	codeOwners := []code.Info{
		{
			OwnerId:      "team-1",
			LocationName: "team-1-service",
		},
	}

	rootDir, dirs := createTestBed(t, codeOwners)
	analyzer, err := NewAnalyzer(rootDir, dirs)
	if err != nil {
		t.Fatal(err)
	}

	// when
	foundOwnerRefs, err := analyzer.Analyze()

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, foundOwnerRefs)
	assert.GreaterOrEqual(t, len(foundOwnerRefs), len(codeOwners))

	for _, ref := range codeOwners {
		assert.Contains(t, foundOwnerRefs, ref)
	}
}

func TestAnalyzer_Analyze_should_fail(t *testing.T) {
	// given
	codeOwners := []code.Info{
		{
			OwnerId:      "team-1",
			LocationName: "team-1-service",
		},
	}

	rootDir, dirs := createTestBed(t, codeOwners)
	analyzer, err := NewAnalyzer(rootDir, dirs)
	if err != nil {
		t.Fatal(err)
	}

	// when
	foundOwnerRefs, err := analyzer.Analyze()

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, foundOwnerRefs)
	assert.Equal(t, len(foundOwnerRefs), len(codeOwners))

	for _, ref := range codeOwners {
		assert.Contains(t, foundOwnerRefs, ref)
	}
}
