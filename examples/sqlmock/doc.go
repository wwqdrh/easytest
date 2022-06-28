package sqlmock

// type Account struct {
// 	ID       int64  `gorm:"primaryKey;column:id"`
// 	UserName string `gorm:"column:username;type:VARCHAR(50)"`
// }

// db, mock, err := SimpleGormMock()
// if err != nil {
// 	t.Fatal(err)
// }
// mock.ExpectBegin()
// mock.ExpectQuery(`INSERT INTO "accounts" (.+) RETURNING`).
// 	WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(int64(1), "kebin"))
// mock.ExpectCommit()
// account := Account{
// 	// ID:       1,
// 	UserName: "kebin",
// }
// err = db.Model(Account{}).Create(&account).Error
// assert.NoError(t, err)
// // verify
// if err := mock.ExpectationsWereMet(); err != nil {
// 	t.Errorf("there were unfulfilled expectations: %s", err)
// }
func DocGormAdd() {

}
