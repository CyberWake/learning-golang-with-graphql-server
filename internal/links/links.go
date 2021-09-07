package links

import (
	"database/sql"
	"log"

	database "github.com/CyberWake/meetmeup/internal/pkg/db/mysql"
	"github.com/CyberWake/meetmeup/internal/users"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

//#2
func (link Link) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

func UpdateLink(link Link) (Link, error) {
	linkInDb, err := GetLink(link.ID)
	if err != nil {
		return Link{}, &LinkNotPresent{}
	}
	if linkInDb.User.ID != link.User.ID {
		return Link{}, &LinkUpdationRightMissing{}
	}
	if len(link.Address) > 0 {
		log.Print(len(link.Address))
		linkInDb.Address = link.Address
	}
	if len(link.Title) > 0 {
		linkInDb.Title = link.Title
	}
	stmt, err := database.Db.Prepare("UPDATE Links SET TITLE = ?,ADDRESS = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(linkInDb.Title, linkInDb.Address, link.ID)
	if err != nil {
		log.Fatal(err)
	}
	return linkInDb, nil
}

func DeleteLink(linkID string, userId string) (string, error) {
	linkInDb, err := GetLink(linkID)
	if err != nil {
		return "", &LinkNotPresent{}
	}
	if linkInDb.User.ID != userId {
		return "", &LinkUpdationRightMissing{}
	}
	stmt, err := database.Db.Prepare("DELETE FROM Links where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(linkID)
	if err != nil {
		log.Fatal(err)
	}
	return "Successfully deleted Link with id ", nil
}

func GetLink(id string) (Link, error) {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID WHERE L.id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	var link Link
	var userID string
	var username string
	err = row.Scan(&link.ID, &link.Title, &link.Address, &userID, &username) // changed
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return Link{}, err
	}
	link.User = &users.User{
		ID:       userID,
		Username: username,
	}
	return link, nil
}

func LinksByUserID(userID string) []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID WHERE U.ID = ?") // changed
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	var username string
	var id string
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username) // changed
		if err != nil {
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		} // changed
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID") // changed
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	var username string
	var id string
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username) // changed
		if err != nil {
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		} // changed
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
