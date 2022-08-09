package usuarios

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/anesquivel/wave-5-backpack/storage/storage_testing/internal/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users(names, last_name, email, age, height, is_active, date_created) VALUES( ?, ?, ?, ?, ?, ?, ? )"))
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	userId := 1

	repo := NewRepository(db)

	nwUser := domain.Usuario{
		Id:          int(userId),
		Names:       "Andrea",
		LastName:    "Esquivel",
		Email:       "prueba@gmail.co",
		Estatura:    1.52,
		IsActivo:    true,
		DateCreated: "2022-08-09",
	}

	u, err := repo.Store(nwUser)
	assert.NoError(t, err)
	assert.NotZero(t, u)
	assert.Equal(t, nwUser.Id, u.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := 1
	repo := NewRepository(db)

	var columns = []string{
		"id", "names", "last_name", "email", "age",
	}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(userId, "", "", "", 0)
	mock.ExpectQuery("SELECT id, names, last_name, email, age FROM users where id = ?").WithArgs(userId).WillReturnRows(rows)

	res, err := repo.GetOne(userId)
	assert.NoError(t, err)
	assert.Equal(t, userId, res.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLASTAGE(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id, lastName, age := 1, "Irwin", 26

	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE users SET last_name = ?, age = ? WHERE id = ?"))
	mock.ExpectExec("UPDATE users").WithArgs(lastName, age, id).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	u, err := repo.UpdateLastNameAndAge(ctx, id, age, lastName)
	assert.NoError(t, err)
	assert.NotZero(t, u)
	assert.Equal(t, id, u.Id)
	assert.Equal(t, age, u.Age)
	assert.Equal(t, lastName, u.LastName)
	assert.NoError(t, mock.ExpectationsWereMet())
	defer cancel()

}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := 1

	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM users WHERE id = ?"))
	mock.ExpectExec("DELETE FROM users").WithArgs(userId).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(db)

	err = repo.Delete(userId)

	assert.Nil(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func init() {
	txdb.Register("txdb", "mysql", "root:@/storage")
}

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("txdb", uuid.New().String())
	if err == nil {
		return db, db.Ping()
	}
	return db, err
}

func Test_SqlRepo_Store(t *testing.T) {
	db, err := InitDb()
	assert.NoError(t, err)

	repo := NewRepository(db)
	userId := uuid.New()
	user := domain.Usuario{
		Id: int(userId.ID()),
	}

	res, err := repo.Store(user)
	assert.Error(t, err)

	res, err = repo.GetOne(user.Id)
	assert.Zero(t, res)

}
