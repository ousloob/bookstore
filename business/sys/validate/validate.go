// Package validate contains the support for validating models.
package validate

import "uuid"

// GenerateID generate a unique ID for entities.
func GenerateID() string {
	return uuid.New().String()
}
