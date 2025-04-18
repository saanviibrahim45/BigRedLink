package models

type User struct {
	ID           string `firestore:"id"`
	Email        string `firestore:"email"`
	PasswordHash string `firestore:"passwordHash"`
	RefreshToken string `firestore:"refreshToken"`
}
