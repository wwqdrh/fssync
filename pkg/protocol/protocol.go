package protocol

import "net/url"

type ProtocolUrl int

const (
	PDownloadList ProtocolUrl = iota
	PDownloadSpec
	PDownloadMd5
	PDownloadTrucate
	PDownloadDelete
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
	default:
		return "/404"
	}
}

func (p ProtocolUrl) ClientUrl(baseurl string, args url.Values) string {
	switch p {
	case PDownloadList:
		extra := ""
		if len(args) > 0 {
			extra = "?" + args.Encode()
		}
		return baseurl + "/download/list" + extra
	case PDownloadSpec:
		extra := ""
		if len(args) > 0 {
			extra = "?" + args.Encode()
		}
		return baseurl + "/download/spec" + extra
	case PDownloadMd5:
		extra := ""
		if len(args) > 0 {
			extra = "?" + args.Encode()
		}
		return baseurl + "/download/md5" + extra
	case PDownloadTrucate:
		extra := ""
		if len(args) > 0 {
			extra = "?" + args.Encode()
		}
		return baseurl + "/download/truncate" + extra
	case PDownloadDelete:
		extra := ""
		if len(args) > 0 {
			extra = "?" + args.Encode()
		}
		return baseurl + "/download/delete" + extra
	default:
		return baseurl + "/404"
	}
}
