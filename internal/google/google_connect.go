package google

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"net/http"
	"os"
)

type Google struct {
	client       *http.Client
	folderIDList []string
	cfg          config.Config
	prefixPath   string

	srv *drive.Service
}

func NewGoogleClient(cfg config.Config) *Google {
	g := Google{
		cfg: cfg,
	}
	return &g
}

func NewGoogleClientTest(prefixPath string) *Google {
	return &Google{
		prefixPath: prefixPath,
	}
}

func (g *Google) getClient(config *oauth2.Config, path string, code string) (*http.Client, error) {
	tokFile := path
	tok, err := g.tokenFromFile(tokFile)
	if err != nil {
		tok, err = g.getTokenFromWeb(config, code)
		if err != nil {
			return nil, errors.Wrap(err, "you have to update config ")
		}
		g.saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func (g *Google) getTokenFromWeb(config *oauth2.Config, code string) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	//todo will make without config and  will return tok and error
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return tok, nil
}

// Retrieves a token from a local file.
func (g *Google) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open("token.json")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open token's file")
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse file")
	}
	return tok, nil

}

func (g *Google) saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "Unable to cache oauth token")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func (g *Google) Connect(path, code string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Unable to read client secret file")
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return errors.Wrap(err, "Unable to parse client secret file to config")
	}
	g.client, err = g.getClient(config, path, code)
	if err != nil {
		return errors.Wrap(err, "Unable to get client")
	}
	g.srv, err = drive.NewService(context.Background(), option.WithHTTPClient(g.client))
	if err != nil {
		return errors.Wrap(err, "failed to create new service")
	}
	return nil
}
