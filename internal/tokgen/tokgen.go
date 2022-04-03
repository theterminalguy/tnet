package tokgen

import (
	"fmt"

	repo "github.com/10hourlabs/tentn/internal/repository"
)

func GenerateRecruiterJWT(email string, m *JWTMeta) (string, error) {
	appInstall := repo.NewSlackAppInstallRepository()
	user, err := appInstall.GetRecruiterByEmail(email)
	if err != nil {
		return "", err
	}
	if user.Role != "recruiter" {
		return "", fmt.Errorf("recruiter not found")
	}
	token, err := NewJWTClaims(user, m).GenerateToken()
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateTalentJWT(email string, m *JWTMeta) (string, error) {
	userRepo := repo.NewUserRepository()
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user.Role != "talent" {
		return "", fmt.Errorf("talent not found")
	}
	token, err := NewJWTClaims(user, m).GenerateToken()
	if err != nil {
		return "", err
	}
	return token, nil
}
