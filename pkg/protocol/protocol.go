package protocol

import "net/url"

type ProtocolUrl int

const (
	PDownloadList ProtocolUrl = iota
	PDownloadSpec
	PDownloadMd5
	PDownloadTrucate
	PDownloadDelete
	PDownloadUpdate
	Unknown
)

func (p ProtocolUrl) ServerUrl() string {
	switch p {
	case PDownloadList:
		return "/download/list"
	case PDownloadSpec:
		return "/download/spec"
	case PDownloadMd5:
		return "/download/md5"
	case PDownloadTrucate:
		return "/download/truncate"
	case PDownloadDelete:
		return "/download/delete"
	case PDownloadUpdate:
		return "/download/update"
	default:
		return "/404"
	}
}

func (p ProtocolUrl) ClientUrl(baseurl string, args url.Values) string {
	extra := ""
	if len(args) > 0 {
		extra = "?" + args.Encode()
	}
	switch p {
	case PDownloadList:
		return baseurl + "/download/list" + extra
	case PDownloadSpec:
		return baseurl + "/download/spec" + extra
	case PDownloadMd5:
		return baseurl + "/download/md5" + extra
	case PDownloadTrucate:
		return baseurl + "/download/truncate" + extra
	case PDownloadDelete:
		return baseurl + "/download/delete" + extra
	case PDownloadUpdate:
		return baseurl + "/download/update" + extra
	default:
		return baseurl + "/404"
	}
}
