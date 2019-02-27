package gins

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	jwt.StandardClaims
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
	Time       int64  `json:"time_in_unix"`
}

// pass by address
func (m *Claim) CheckByObject(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Where(map[string]interface{}{"external_id": m.ExternalID}).Find(obj).Error
}

func JWT(userID int, externalId, issuer, signedString string) (string, error) {
	var (
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
			ID:         userID,
			ExternalID: externalId,
			Time:       time.Now().UTC().Unix(),
			StandardClaims: jwt.StandardClaims{
				Issuer:    issuer,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: time.Now().AddDate(0, 0, 14).UTC().Unix(),
				Id:        "",
				Audience:  "",
				Subject:   "",
			},
		})
		tkn, err = token.SignedString(signedString)
	)
	if err != nil {
		return "", err
	}
	return tkn, nil
}
