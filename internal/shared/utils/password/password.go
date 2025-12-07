package password

import "golang.org/x/crypto/bcrypt"

func Encrypt(p string) (string, error) {
	np, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(np), nil
}

func AreEqual(p string, dbp string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(dbp), []byte(p)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
