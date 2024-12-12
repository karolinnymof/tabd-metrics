package utils

import (
	"fmt"

	"github.com/eduardor2m/tabd-metrics/src/models"
	"github.com/go-faker/faker/v4"
)

func GenerateFakeUsers(num int) ([]models.User, error) {
	users := make([]models.User, num)
	if err := faker.FakeData(&users); err != nil {
		return nil, fmt.Errorf("erro ao gerar dados fictícios: %w", err)
	}
	return users, nil
}
