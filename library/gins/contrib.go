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
	ID           uint                   `json:"id"`
	ExternalID   string                 `json:"external_id"`
	Time         int64                  `json:"time_in_unix"`
}

// pass by address
func (m *Claim) CheckByObject(DB *gorm.DB, obj interface{}) (err error) {
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
	return DB.Debug().Scopes(scopes...).Find(obj).Error
}

// pass by address
func (m *Claim) CheckByCustomFields(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Debug().Where(m.CustomFields).Find(obj).Error
}

// pass by address
func (m *Claim) CheckByExternalID(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Debug().Where(map[string]interface{}{"external_id": m.ExternalID}).Find(obj).Error
}

// pass by address
func (m *Claim) CheckByID(DB *gorm.DB, obj interface{}) (err error) {
	return DB.Debug().Where(map[string]interface{}{"id": m.ID}).Find(obj).Error
}

func (m *Claim) IsExpired() bool {
	return time.Now().After(time.Unix(m.ExpiresAt, 0))
}

func JWT(userID uint, expireAfterInMin int, externalId, issuer, signedString string, customFields map[string]interface{}) (string, error) {
	var (
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
			CustomFields: customFields,
			ID:           userID,
			ExternalID:   externalId,
			Time:         time.Now().UTC().Unix(),
			StandardClaims: jwt.StandardClaims{
				Issuer:    issuer,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: time.Now().AddDate(0, 0, expireAfterInMin/(60*24)).UTC().Unix(),
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
