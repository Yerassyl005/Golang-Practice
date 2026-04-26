package repository

import (
	"database/sql"
)

type User struct {
	ID int
	Name string
	Email string
	Gender string
	BirthDate string
}

type PaginatedResponse struct {
	Data []User
	TotalCount int
	Page int
	PageSize int
}

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int, name string, email string, order string) (PaginatedResponse, error) {

	offset := (page - 1) * pageSize

	query := "SELECT id,name,email,gender,birth_date FROM users WHERE 1=1"

	if name != "" {
		query += " AND name ILIKE '%" + name + "%'"
	}

	if email != "" {
		query += " AND email ILIKE '%" + email + "%'"
	}

	if order == "" {
		order = "id"
	}

	query += " ORDER BY " + order
	query += " LIMIT $1 OFFSET $2"

	rows, err := r.DB.Query(query, pageSize, offset)
	if err != nil {
		return PaginatedResponse{}, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
		users = append(users, u)
	}

	var count int

	r.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)

	resp := PaginatedResponse{
		Data: users,
		TotalCount: count,
		Page: page,
		PageSize: pageSize,
	}

	return resp, nil
}

func (r *Repository) GetCommonFriends(user1 int, user2 int) ([]User, error) {

	query := `
	SELECT u.id,u.name,u.email,u.gender,u.birth_date
	FROM users u
	JOIN user_friends f1 ON u.id = f1.friend_id
	JOIN user_friends f2 ON u.id = f2.friend_id
	WHERE f1.user_id = $1 AND f2.user_id = $2
	`

	rows, err := r.DB.Query(query, user1, user2)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID,&u.Name,&u.Email,&u.Gender,&u.BirthDate)
		users = append(users,u)
	}

	return users,nil
}
