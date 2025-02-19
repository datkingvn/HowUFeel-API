package helpers

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(password *string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	hashedPwd := string(bytes)
	return &hashedPwd
}

func VerifyPassword(foundPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPwd), []byte(pwd))
	return err == nil, err
}
