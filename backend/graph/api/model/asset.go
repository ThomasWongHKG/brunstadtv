package model

import (
	"context"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/bcc-code/brunstadtv/backend/common"
	"github.com/samber/lo"
)

type signatureProvider interface {
	SignAzureURL(*url.URL, string) (string, error)
	SignCloudfrontURL(string, string) (string, error)
	SignWithPolicy(string, *sign.Policy) (string, error)
}

type cdnConfig interface {
	GetLegacyVODDomain() string
	GetVOD2Domain() string
}

// FileFrom converts Assetfile rows to the GQL equvivalents
func FileFrom(ctx context.Context, signer signatureProvider, cdnDomain string, file *common.File) *File {
	var subLang *Language
	if file.SubtitleLanguage.Valid {
		l := Language(file.SubtitleLanguage.String)
		subLang = &l
	}

	url := url.URL{
		Path:   file.Path,
		Host:   cdnDomain,
		Scheme: "https",
	}

	policy := sign.NewCannedPolicy(url.String(), time.Now().Add(time.Hour))

	signed, err := signer.SignWithPolicy(url.String(), policy)
	if err != nil {
		panic(err)
	}

	return &File{
		ID:               strconv.Itoa(file.ID),
		URL:              signed,
		FileName:         path.Base(file.Path),
		AudioLanguage:    Language(file.AudioLanguage.String),
		SubtitleLanguage: subLang,
		MimeType:         file.MimeType,
	}
}

// StreamFrom converts Assetfile rows to the GQL equvivalents
func StreamFrom(ctx context.Context, signer signatureProvider, cdn cdnConfig, stream *common.Stream) (*Stream, error) {
	signedURL := ""
	var err error

	if stream.Service == common.StreamServiceAzureMedia {
		streamURL, err := url.Parse(stream.Url)
		if err != nil {
			return nil, err
		}

		streamURL.Host = cdn.GetLegacyVODDomain()
		streamURL.Scheme = "https"

		// This is intentionally hardcoded for now
		manifestURL, _ := url.Parse("https://proxy.brunstad.tv/api/vod/toplevelmanifest")

		q := manifestURL.Query()
		q.Add("playbackUrl", streamURL.String())
		manifestURL.RawQuery = q.Encode()

		signedURL, err = signer.SignAzureURL(manifestURL, stream.EncryptionKeyID.ValueOrZero())
	} else {
		signedURL, err = signer.SignCloudfrontURL(stream.Path, cdn.GetVOD2Domain())
	}

	if err != nil {
		return nil, err
	}

	return &Stream{
		ID:                strconv.Itoa(stream.ID),
		URL:               signedURL,
		AudioLanguages:    lo.Map(stream.AudioLanguages, func(s string, _ int) Language { return Language(s) }),
		SubtitleLanguages: lo.Map(stream.SubtitleLanguages, func(s string, _ int) Language { return Language(s) }),
		Type:              StreamType(stream.Type),
	}, nil
}
