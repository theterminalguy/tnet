package service

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/theterminalguy/tnet/ent"
)

type PDFService struct {
	Location *GeoLocationService
}

func NewPDFService() *PDFService {
	return &PDFService{
		Location: NewGeoLocationService(),
	}
}

type TalentProfile struct {
	FirstName           string
	LastName            string
	PreferredJobTitle   string
	Pronoun             string
	Email               string
	Phone               string
	Location            string
	Timezone            string
	ProfessionalSummary string
	Available           bool
	JobPreference       string
	CareerStartDate     time.Time
	Skills              []*ent.Skill
	Portfolios          []*ent.PortfolioLink
	WorkExperiences     []*ent.WorkExperience
	Educations          []*ent.Education
}

func (p *PDFService) Generate(u *ent.Talent) ([]byte, error) {
	t, err := template.ParseFiles("cmd/generate/templates/talent-profile.html")
	if err != nil {
		return nil, err
	}
	countryName, err := p.Location.GetCountryNameByCode(u.CountryCode)
	if err != nil {
		return nil, err
	}
	data := TalentProfile{
		FirstName:           u.FirstName,
		LastName:            u.LastName,
		PreferredJobTitle:   u.PreferredJobTitle,
		Pronoun:             u.Pronoun,
		Email:               u.Email,
		Phone:               u.Phone,
		Location:            fmt.Sprintf("%s, %s", u.City, countryName),
		Timezone:            u.Timezone,
		ProfessionalSummary: u.ProfessionalSummary,
		Available:           u.IsAvailable,
		JobPreference:       u.JobPreference.String(),
		CareerStartDate:     u.ProfessionalStartDate,
		Skills:              u.Edges.Skills,
		Portfolios:          u.Edges.Portfoliolinks,
		WorkExperiences:     u.Edges.WorkExperiences,
		Educations:          u.Edges.Educations,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return nil, err
	}

	// initalize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)

	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// magic
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}
	return pdfg.Bytes(), nil
}
