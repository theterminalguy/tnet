package task

import (
	"fmt"

	"github.com/10hourlabs/tentn/internal/tokgen"
	"github.com/10hourlabs/tentn/util"
)

type TokGen struct {
}

func NewTokGen() *TokGen {
	return &TokGen{}
}

func (t *TokGen) Run(params string) error {
	m := util.StringParamsToMap(params)
	toktype, ok := m["type"]
	if !ok {
		return fmt.Errorf("provide a type, could be one of: %v", []string{"recruiter", "talent"})
	}
	if toktype != "recruiter" && toktype != "talent" {
		return fmt.Errorf("provide a type, could be one of: %v", []string{"recruiter", "talent"})
	}
	meta := &tokgen.JWTMeta{
		Audience: "postman",
		Issuer:   "localhost:tentn",
	}
	if toktype == "recruiter" {
		recruiterEmail, ok := m["email"]
		if !ok {
			return fmt.Errorf("provide a recruiter's email address")
		}
		tok, err := tokgen.GenerateRecruiterJWT(recruiterEmail, meta)
		if err != nil {
			return err
		}
		fmt.Println("Recruiter JWT:\n", tok)
		return nil
	}
	talentEmail, ok := m["email"]
	if !ok {
		return fmt.Errorf("provide a talent's email address")
	}
	tok, err := tokgen.GenerateTalentJWT(talentEmail, meta)
	if err != nil {
		return err
	}
	fmt.Println("Talent JWT:\n", tok)
	return nil
}
