package ipfs

import (
	"bytes"
	"encoding/json"
	shell "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/config"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Connector struct {
	cfg config.Config
}

type ERC721json struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ExternalUrl string `json:"external_url"`
}

//
//type ERC721DescriptionInfo struct {
//	types       string `json:"type"`
//	description string `json:"description"`
//}

func NewConnector(cfg config.Config) *Connector {
	return &Connector{cfg: cfg}
}

func (i Connector) Upload(data []byte) (string, error) {
	ipfs := shell.NewShellWithClient(i.cfg.NetworksConfig().IPFSEndpoint, NewClient(i.cfg.NetworksConfig().IpfsPrId, i.cfg.NetworksConfig().IpfsPrKey))

	fileHash, err := ipfs.Add(bytes.NewReader(data))
	if err != nil {
		return "", errors.Wrap(err, "failed to upload")
	}
	return fileHash, nil

}

func (i Connector) PrepareJSON(tokenName, tokenDescription, imagePath string) ([]byte, error) {
	erc721 := ERC721json{
		Name:        tokenName,
		Description: tokenDescription,
		Image:       "https://ipfs.io/ipfs/" + imagePath,
		ExternalUrl: "https://dlt-academy.com/certificates",
	}

	erc721JSON, err := json.Marshal(erc721)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal erc721")
	}

	return erc721JSON, nil
}

func (i Connector) PrepareImage(imagePath string) ([]byte, error) {
	path, err := filepath.Abs("main.go")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path")
	}
	path = strings.ReplaceAll(path, "main.go", "")

	infile, err := os.Open(path + imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open image")
	}
	img, err := png.Decode(infile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image to []byte")
	}
	return buf.Bytes(), nil
}

func NewClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

// authTransport decorates each request with a basic auth header.
type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}
