package audiofile

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"talkliketv.click/tltv/internal/test"
	"talkliketv.click/tltv/internal/util"
	"testing"
)

type audioFileTestCase struct {
	name       string
	buildFile  func(*testing.T) *os.File
	checkLines func([]string, error)
}

func TestGetLines(t *testing.T) {

	testCases := []audioFileTestCase{
		{
			name: "No error",
			buildFile: func(t *testing.T) *os.File {
				return createFile(
					t,
					"noerror",
					"This is the first sentence.\nThis is the second sentence.\n")
			},
			checkLines: func(lines []string, err error) {
				require.NoError(t, err)
				require.Equal(t, len(lines), 2)
			},
		},
		{
			name: "parse srt",
			buildFile: func(t *testing.T) *os.File {
				srtString := `654
				00:34:22,393 > 00:34:25,271
				¿El camión a Tepatitlán?
				Saliendo, segundo andén.

				655
				00:34:25,354 > 00:34:28,441
				Por favor, nada más debo entregar esto.
					Un segundo, por favor.

				656
				00:34:29,192 > 00:34:31,444
				Déjala pasar, mi Johnny.
					Gracias.`
				return createFile(
					t,
					"parsesrt",
					srtString)
			},
			checkLines: func(lines []string, err error) {
				require.NoError(t, err)
				require.Equal(t, len(lines), 3)
			},
		},
		{
			name: "Multi newline",
			buildFile: func(t *testing.T) *os.File {
				return createFile(
					t,
					"noerror",
					"This is the first sentence.\n\n\n\n\n\n\nThis is the second sentence.\n")
			},
			checkLines: func(lines []string, err error) {
				require.NoError(t, err)
				require.Equal(t, len(lines), 2)
			},
		},
		{
			name: "paragraph",
			buildFile: func(t *testing.T) *os.File {
				return createFile(
					t,
					"noerror",
					"This is the first. This is the second. This is the third. this is the fourth\nThis is the fifth")
			},
			checkLines: func(lines []string, err error) {
				require.NoError(t, err)
				require.Equal(t, len(lines), 5)
			},
		},
		{
			name: "too short",
			buildFile: func(t *testing.T) *os.File {
				return createFile(
					t,
					"noerror",
					"This is the. This is. This is the. this is the\nThis is the")
			},
			checkLines: func(lines []string, err error) {
				require.Errorf(t, err, "unable to parse file")
			},
		},
		{
			name: "empty file",
			buildFile: func(t *testing.T) *os.File {
				return createFile(
					t,
					"noerror",
					"")
			},
			checkLines: func(lines []string, err error) {
				require.Errorf(t, err, "unable to parse file")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/fakeurl", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			file := tc.buildFile(t)
			audioFile := AudioFile{}
			stringsSlice, err := audioFile.GetLines(c, file)
			tc.checkLines(stringsSlice, err)
		})
	}
}

func TestBuildAudioInputFiles(t *testing.T) {
	// BuildAudioInputFiles(e echo.Context, ids []int64, t db.Title, pause, from, to, tmpDir string)

	title := test.RandomTitle()
	pause := util.RandomString(4)
	from := util.RandomString(4)
	to := util.RandomString(4)
	tmpDir := test.AudioBasePath + "TestBuildAudioInputFiles/" + title.Title + "/"
	fromPath := fmt.Sprintf("%s%s/", tmpDir, from)
	toPath := fmt.Sprintf("%s%s/", tmpDir, to)
	err := os.MkdirAll(tmpDir, 0777)
	require.NoError(t, err)
	testCases := []audioFileTestCase{
		{
			name: "No error",
			checkLines: func(lines []string, err error) {
				require.NoError(t, err)
				require.Equal(t, len(lines), 2)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/fakeurl", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			audioFile := AudioFile{}
			err := audioFile.BuildAudioInputFiles(
				c,
				[]int64{util.RandomInt64(),
					util.RandomInt64()},
				title,
				pause,
				fromPath,
				toPath,
				tmpDir,
			)
			require.NoError(t, err)
			filePath := tmpDir + title.Title + "-input-1"
			require.FileExists(t, filePath)
		})
	}
}

func createFile(t *testing.T, filename, fileString string) *os.File {
	// Create a new file
	file, err := os.Create(filename)
	require.NoError(t, err)
	defer os.Remove(filename)

	// Write to the file
	_, err = file.WriteString(fileString)
	require.NoError(t, err)
	// Ensure data is written to disk
	err = file.Sync()
	require.NoError(t, err)

	return file
}
