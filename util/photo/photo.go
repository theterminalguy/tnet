package photo

import "fmt"

func GenerateDefaultPhoto(firstName, lastName string) string {
	return fmt.Sprintf(
		"https://ui-avatars.com/api/?name=%s+%s&background=random&size=64",
		firstName,
		lastName,
	)
}
