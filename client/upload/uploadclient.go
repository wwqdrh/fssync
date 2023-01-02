package upload

import (
	"fmt"
	"io"
	netUrl "net/url"
	"strconv"

	"net/http"

	"github.com/wwqdrh/fssync/pkg/store"
)

const (
	ProtocolVersion = "1.0.0"
)

// Config provides a way to configure the Client depending on your needs.
type UploadConfig struct {
	// ChunkSize divide the file into chunks.
	ChunkSize int64
	// Resume enables resumable upload.
	Resume bool
	// OverridePatchMethod allow to by pass proxies sendind a POST request instead of PATCH.
	OverridePatchMethod bool
	// Store map an upload's fingerprint with the corresponding upload URL.
	// If Resume is true the Store is required.
	Store store.UploadStore
	// Set custom header values used in all requests.
	Header http.Header
	// HTTP Client
	HttpClient *http.Client
}

// DefaultConfig return the default Client configuration.
func DefaultUploadConfig() *UploadConfig {
	return &UploadConfig{
		ChunkSize:           2 * 1024 * 1024,
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              make(http.Header),
		HttpClient:          nil,
	}
}

// Validate validates the custom configuration.
func (c *UploadConfig) Validate() error {
	if c.ChunkSize < 1 {
		return ErrChuckSize
	}

	if c.Resume && c.Store == nil {
		return ErrNilStore
	}

	return nil
}

// Client represents the tus client.
// You can use it in goroutines to create parallels uploads.
type UploadClient struct {
	Config  *UploadConfig
	Url     string
	Version string
	Header  http.Header

	client *http.Client
}

// NewClient creates a new tus client.
func NewUploadClient(url string, config *UploadConfig) (*UploadClient, error) {
	if config == nil {
		config = DefaultUploadConfig()
	} else {
		if err := config.Validate(); err != nil {
			return nil, err
		}
	}

	if config.Header == nil {
		config.Header = make(http.Header)
	}

	if config.HttpClient == nil {
		config.HttpClient = &http.Client{}
	}

	return &UploadClient{
		Config:  config,
		Url:     url,
		Version: ProtocolVersion,
		Header:  config.Header,

		client: config.HttpClient,
	}, nil
}

func (c *UploadClient) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.Header {
		req.Header[k] = v
	}

	req.Header.Set("Tus-Resumable", ProtocolVersion)

	return c.client.Do(req)
}

// CreateUpload creates a new upload in the server.
func (c *UploadClient) CreateUpload(u *Upload) (*Uploader, error) {
	if u == nil {
		return nil, ErrNilUpload
	}

	if c.Config.Resume && len(u.Fingerprint) == 0 {
		return nil, ErrFingerprintNotSet
	}

	req, err := http.NewRequest("POST", c.Url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(u.size, 10))
	req.Header.Set("Upload-Metadata", u.EncodedMetadata())

	res, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 201:
		location := res.Header.Get("Location")

		newURL, err := c.resolveLocationURL(location)
		if err != nil {
			return nil, err
		}

		if c.Config.Resume {
			if err := c.Config.Store.Set(u.Fingerprint, newURL.String()); err != nil {
				return nil, err
			}
		}

		return NewUploader(c, newURL.String(), u, 0), nil
	case 412:
		return nil, ErrVersionMismatch
	case 413:
		return nil, ErrLargeUpload
	default:
		return nil, newClientError(res)
	}
}

func (c *UploadClient) resolveLocationURL(location string) (*netUrl.URL, error) {
	baseURL, err := netUrl.Parse(c.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL '%s'", c.Url)
	}

	locationURL, err := netUrl.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("invalid Location header value for a Creation '%s': %s", location, err.Error())
	}
	newURL := baseURL.ResolveReference(locationURL)
	return newURL, nil
}

// ResumeUpload resumes the upload if already created, otherwise it will return an error.
func (c *UploadClient) ResumeUpload(u *Upload) (*Uploader, error) {
	if u == nil {
		return nil, ErrNilUpload
	}

	if !c.Config.Resume {
		return nil, ErrResumeNotEnabled
	} else if len(u.Fingerprint) == 0 {
		return nil, ErrFingerprintNotSet
	}

	url, found := c.Config.Store.Get(u.Fingerprint)

	if !found {
		return nil, ErrUploadNotFound
	}

	offset, err := c.getUploadOffset(url)

	if err != nil {
		return nil, err
	}

	return NewUploader(c, url, u, offset), nil
}

// CreateOrResumeUpload resumes the upload if already created or creates a new upload in the server.
func (c *UploadClient) CreateOrResumeUpload(u *Upload) (*Uploader, error) {
	if u == nil {
		return nil, ErrNilUpload
	}

	uploader, err := c.ResumeUpload(u)

	if err == nil {
		return uploader, err
	} else if (err == ErrResumeNotEnabled) || (err == ErrUploadNotFound) {
		return c.CreateUpload(u)
	}

	return nil, err
}

func (c *UploadClient) uploadChunck(url string, body io.Reader, size int64, offset int64) (int64, error) {
	var method string

	if !c.Config.OverridePatchMethod {
		method = "PATCH"
	} else {
		method = "POST"
	}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.Header.Set("Upload-Offset", strconv.FormatInt(offset, 10))

	if c.Config.OverridePatchMethod {
		req.Header.Set("X-HTTP-Method-Override", "PATCH")
	}

	res, err := c.Do(req)

	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 204:
		if newOffset, err := strconv.ParseInt(res.Header.Get("Upload-Offset"), 10, 64); err == nil {
			return newOffset, nil
		} else {
			return -1, err
		}
	case 409:
		return -1, ErrOffsetMismatch
	case 412:
		return -1, ErrVersionMismatch
	case 413:
		return -1, ErrLargeUpload
	default:
		return -1, newClientError(res)
	}
}

func (c *UploadClient) getUploadOffset(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		return -1, err
	}

	res, err := c.Do(req)

	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		i, err := strconv.ParseInt(res.Header.Get("Upload-Offset"), 10, 64)

		if err == nil {
			return i, nil
		} else {
			return -1, err
		}
	case 403, 404, 410:
		// file doesn't exists.
		return -1, ErrUploadNotFound
	case 412:
		return -1, ErrVersionMismatch
	default:
		return -1, newClientError(res)
	}
}

func newClientError(res *http.Response) ClientError {
	body, _ := io.ReadAll(res.Body)
	return ClientError{
		Code: res.StatusCode,
		Body: body,
	}
}
