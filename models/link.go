package models

import (
	"database/sql"
	"tvbox/db"
)

type Link struct {
	ID    int
	URL   string
	State int
	Name  string
}

func GetActiveLink() (string, error) {
	var url string
	err := db.DB.QueryRow("SELECT url FROM tvbox WHERE state = 1 LIMIT 1").Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到活跃链接
			return "", sql.ErrNoRows
		}
		// 其他数据库错误
		return "", err
	}
	return url, nil
}

func GetAllLinks() ([]Link, error) {
	rows, err := db.DB.Query("SELECT id, url, state, name FROM tvbox")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.URL, &link.State, &link.Name)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func AddLink(url, name string, state int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取当前最大ID
	var maxID int
	err = tx.QueryRow("SELECT COALESCE(MAX(id), 0) FROM tvbox").Scan(&maxID)
	if err != nil {
		return err
	}

	// 插入新记录
	_, err = tx.Exec("INSERT INTO tvbox (id, url, name, state) VALUES (?, ?, ?, ?)", maxID+1, url, name, state)
	if err != nil {
		return err
	}

	// 重置自增ID
	_, err = tx.Exec("DELETE FROM sqlite_sequence WHERE name='tvbox'")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func DeleteLink(id int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除指定链接
	_, err = tx.Exec("DELETE FROM tvbox WHERE id = ?", id)
	if err != nil {
		return err
	}

	// 更新剩余链接的ID
	_, err = tx.Exec("UPDATE tvbox SET id = id - 1 WHERE id > ?", id)
	if err != nil {
		return err
	}

	// 重置自增ID
	_, err = tx.Exec("DELETE FROM sqlite_sequence WHERE name='tvbox'")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ChangeState(id int) error {
	// 查询当前id的状态
	var currentState int
	err := db.DB.QueryRow("SELECT state FROM tvbox WHERE id = ?", id).Scan(&currentState)
	if err != nil {
		return err
	}

	// 如果当前状态为1，则更新为0；否则，将所有状态置为0，然后更新为1
	if currentState == 1 {
		_, err = db.DB.Exec("UPDATE tvbox SET state = 0 WHERE id = ?", id)
	} else {
		_, err = db.DB.Exec("UPDATE tvbox SET state = 0")
		if err != nil {
			return err
		}
		_, err = db.DB.Exec("UPDATE tvbox SET state = 1 WHERE id = ?", id)
	}

	return err
}

func UpdateLink(id int, url, name string) error {
	_, err := db.DB.Exec("UPDATE tvbox SET url = ?, name = ? WHERE id = ?", url, name, id)
	return err
}
