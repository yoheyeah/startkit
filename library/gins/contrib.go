package gins

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	jwt.StandardClaims
	CustomFields map[string]interface{} `json:"custom_fields"`
	ID           uint                   `json:"-"`
	ExternalID   string                 `json:"external_id"`
	Time         int64                  `json:"time_in_unix"`
}

// pass by address
func (m *Claim) FindByObject(DB *gorm.DB, obj interface{}) (err error) {
	var (
		scopes []func(db *gorm.DB) *gorm.DB
	)
	if m.ID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(map[string]interface{}{"id": m.ID}) })
	}
	if m.ExternalID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(map[string]interface{}{"external_id": m.ExternalID}) })
	}
	if len(m.CustomFields) > 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB { return db.Where(m.CustomFields) })
	}
	if len(scopes) <= 0 {
		return errors.New("Invalid JWT Claim")
	}
	return DB.Scopes(scopes...).Find(obj).Error
}

// pass by address
func (m *Claim) FindByCustomFields(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Where(m.CustomFields).Find(obj).Error
}

// pass by address
func (m *Claim) FindByExternalID(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Where(map[string]interface{}{"external_id": m.ExternalID}).Find(obj).Error
}

// pass by address
func (m *Claim) FindByID(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Where(map[string]interface{}{"id": m.ID}).Find(obj).Error
}

func (m *Claim) IsExpired() bool {
	return time.Now().After(time.Unix(m.ExpiresAt, 0))
}

func JWT(userID uint, expireAfterInMin int, externalId, issuer, signedString string, customFields map[string]interface{}) (string, error) {
	var (
		token = jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			Claim{
				CustomFields: customFields,
				ID:           userID,
				ExternalID:   externalId,
				Time:         time.Now().Unix(),
				StandardClaims: jwt.StandardClaims{
					Issuer:    issuer,
					IssuedAt:  time.Now().Unix(),
					NotBefore: time.Now().Unix(),
					ExpiresAt: time.Now().Add(time.Duration(expireAfterInMin) * time.Minute).Unix(),
					Id:        "",
					Audience:  "",
					Subject:   "",
				},
			})
		tkn, err = token.SignedString([]byte(signedString))
	)
	if err != nil {
		return "", err
	}
	return tkn, nil
}
