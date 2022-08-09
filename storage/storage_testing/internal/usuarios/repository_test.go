package usuarios

import (
	"log"
	"testing"
	"time"

	"github.com/anesquivel/wave-5-backpack/storage/storage_testing/db"
	"github.com/anesquivel/wave-5-backpack/storage/storage_testing/internal/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestStore(t *testing.T) {
	db.Init()
	repo := NewRepository(db.StorageDB)
	newUser := domain.Usuario{
		Names:       "Ashton",
		LastName:    "Brooke",
		Email:       "ash2@gmail.com",
		Estatura:    1.80,
		IsActivo:    true,
		DateCreated: "2022-08-08",
		Age:         28,
	}

	userResult, err := repo.Store(newUser)

	if err != nil {
		log.Println("----- ERROR- TEST:", err.Error())
	}

	assert.Equal(t, 1, userResult.Id)
	assert.Equal(t, newUser.Names, userResult.Names)
}

func TestByName(t *testing.T) {
	db.Init()
	repo := NewRepository(db.StorageDB)
	name := "Ashton"
	userResult, err := repo.GetByName(name)

	if err != nil {
		log.Println("----- ERROR- TEST:", err.Error())
	}

	assert.Equal(t, 1, userResult.Id)
}

func TestGetAll(t *testing.T) {
	db.Init()
	repo := NewRepository(db.StorageDB)
	totalOfUsers := 4
	userResult, err := repo.GetAll()

	if err != nil {
		log.Println("----- ERROR- TEST:", err.Error())
	}

	assert.Equal(t, totalOfUsers, len(userResult))
}

func TestUpdateLASTAGE(t *testing.T) {
	db.Init()
	repo := NewRepository(db.StorageDB)
	id, lastName, age := 1, "Irwin", 26
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userResult, err := repo.UpdateLastNameAndAge(ctx, id, age, lastName)

	if err != nil {
		log.Println("----- ERROR- TEST:", err.Error())
	}

	assert.Equal(t, age, userResult.Age)
	assert.Equal(t, lastName, userResult.LastName)

}
