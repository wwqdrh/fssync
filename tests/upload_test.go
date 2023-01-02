package tests

// var f = int64(0) // atomic

// func init() {
// 	server.ServerFlag.Port = 1080
// 	server.ServerFlag.Store = "./testdata/store"
// 	server.ServerFlag.Urlpath = "/files/"
// 	server.ServerFlag.Store = "./testdata/store"

// 	client.ClientUploadFlag.Host = "http://127.0.0.1:1080/files/"
// 	client.ClientUploadFlag.SpecPath = "./testdata/uploadspec"
// }

// type UploadSuite struct {
// 	suite.Suite
// }

// func TestUploadSuite(t *testing.T) {
// 	suite.Run(t, &UploadSuite{})
// }

// func (s *UploadSuite) TestCreateUploadFile() {
// 	client.ClientUploadFlag.Uploadfile = "./testdata/testupload.txt"
// 	if err := client.UploadStart(); err != nil {
// 		s.T().Error(err)
// 	}
// }

// func (s *UploadSuite) TestResumeUploadFile() {
// 	_, err := os.Stat("./testdata/video.mp4")
// 	if errors.Is(err, os.ErrNotExist) {
// 		s.T().Skip("大文件未加入版本控制中，要测试请手动加入")
// 	}
// 	client.ClientUploadFlag.Uploadfile = "./testdata/video.mp4"

// 	if err := client.UploadStart(); err != nil {
// 		s.T().Error(err)
// 	}
// }
