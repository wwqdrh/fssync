package tests

// type DownloadSuite struct {
// 	suite.Suite
// }

// func init() {
// 	server.ServerFlag.Port = 1080
// 	server.ServerFlag.Store = "./testdata/store"
// 	server.ServerFlag.Urlpath = "/files/"
// 	server.ServerFlag.ExtraPath = "./testdata/downloadextra"

// 	client.ClientUploadFlag.Host = "http://127.0.0.1:1080/files/"
// 	client.ClientUploadFlag.SpecPath = "./testdata/uploadspec"

// 	client.ClientUploadFlag.Host = "http://127.0.0.1:1080/files/"
// 	client.ClientDownloadFlag.DownloadPath = "./testdata/downloadstore"
// 	client.ClientDownloadFlag.TempPath = "./testdata/temppath"
// 	client.ClientDownloadFlag.SpecPath = "./testdata/downloadspec"
// }

// func TestDownloadSuite(t *testing.T) {
// 	suite.Run(t, &DownloadSuite{})
// }

// func (s *DownloadSuite) TestFileUrlList() {
// 	req, err := http.NewRequest("GET", "http://127.0.0.1:1080/download/list", nil)
// 	require.Nil(s.T(), err)
// 	resp, err := http.DefaultClient.Do(req)
// 	require.Nil(s.T(), err)
// 	defer resp.Body.Close()
// 	data, err := io.ReadAll(resp.Body)
// 	require.Nil(s.T(), err)
// 	fmt.Println(string(data))
// }

// func (s *DownloadSuite) TestCreateUploadFile() {
// 	client.ClientDownloadFlag.DownloadUrl = "http://localhost:1080"
// 	client.ClientDownloadFlag.FileName = "testdownload.txt"
// 	if err := client.DownloadStart(); err != nil {
// 		s.T().Error(err)
// 	}
// }

// func (s *DownloadSuite) TestResumeDownloadFile() {
// 	_, err := os.Stat("./testdata/video.mp4")
// 	if errors.Is(err, os.ErrNotExist) {
// 		s.T().Skip("大文件未加入版本控制中，要测试请手动加入")
// 	}
// 	client.ClientDownloadFlag.DownloadUrl = "http://localhost:1080"
// 	client.ClientDownloadFlag.FileName = "video.mp4"

// 	if err := client.DownloadStart(); err != nil {
// 		s.T().Error(err)
// 	}
// }

// func (s *DownloadSuite) TestDownloadaAll() {
// 	client.ClientDownloadFlag.DownloadUrl = "http://localhost:1080"
// 	client.ClientDownloadFlag.DownAll = true

// 	if err := client.DownloadStart(); err != nil {
// 		s.T().Error(err)
// 	}
// }
