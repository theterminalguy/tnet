package service

import (
	"bytes"
	"html/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFService struct {
}

func NewPDFService() *PDFService {
	return &PDFService{}
}

type CandidateProfile struct {
	FirstName           string
	LastName            string
	PreferredJobTitle   string
	Pronoun             string
	Email               string
	Phone               string
	Location            string
	Timezone            string
	ProfessionalSummary string
	Available           string
	JobPreference       string
	CareerStartDate     string
}

func (p *PDFService) Generate() ([]byte, error) {
	t, err := template.ParseFiles("cmd/generate/templates/candidate-profile.html")
	if err != nil {
		return nil, err
	}
	data := CandidateProfile{
		FirstName:           "John",
		LastName:            "Doe",
		PreferredJobTitle:   "Software Engineer",
		Pronoun:             "He/Him",
		Email:               "johndoe@example.com",
		Phone:               "555-555-5555",
		Location:            "San Francisco, CA",
		Timezone:            "GMT-08:00",
		ProfessionalSummary: "I am a software engineer with a background in electrical engineering and computer science. I am currently working at a startup called tentn.io. I am a self-taught programmer with a passion for learning and problem solving. I am currently looking for a software engineer position in the San Francisco Bay Area.",
		Available:           "Immediately",
		JobPreference:       "Remote",
		CareerStartDate:     "January 1, 2020",
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
