package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/10hourlabs/tentn/internal/service/filestorage"
	"github.com/10hourlabs/tentn/util/osutil"
	"github.com/google/uuid"
)

type TalentImageParams struct {
	TalentId uuid.UUID `json:"talent_id"`
	ImageUrl string    `json:"image_url"`
	path     string
}

func NewTalentImageHandler(talentId uuid.UUID, imageUrl string) *TalentImageParams {
	return &TalentImageParams{
		TalentId: talentId,
		ImageUrl: imageUrl,
	}
}

func (t *TalentImageParams) GetImage() (string, error) {
	path := "data/img"
	t.path = path
	fp := fmt.Sprintf("%s/%v.jpg", path, t.TalentId.String())
	// Get the image from url
	resp, err := http.Get(t.ImageUrl)
	fmt.Println(t.ImageUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if osutil.InDevMode() {
		link, err := t.saveToLocal(fp, resp.Body)
		if err != nil {
			return "", err
		}
		return link, nil
	}

	link, err := t.saveToGoogle(fp, resp.Body)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (t *TalentImageParams) saveToGoogle(fpath string, f io.ReadCloser) (string, error) {
	g := filestorage.NewGoogleBucketFileStorage(fpath)
	link, err := g.Upload(f)

	if err != nil {
		return "", err
	}

	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	p := fmt.Sprintf("https://storage.googleapis.com%v", u.Path)
	return p, nil
}

func (t *TalentImageParams) saveToLocal(fpath string, f io.ReadCloser) (string, error) {
	err := os.MkdirAll(t.path, os.ModePerm)
	if err != nil {
		return "", err
	}
	d, err := os.Create(fpath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(d, f)
	if err != nil {
		return "", err
	}
	return fpath, nil
}
