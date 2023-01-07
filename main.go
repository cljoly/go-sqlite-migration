/*
The APACHE License (APACHE)

Copyright (c) 2023 Clément Joly. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/BurntSushi/migration"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db, err := migration.OpenWith("sqlite3", "./my.db",
		[]migration.Migrator{
			func(tx migration.LimitedTx) error {
				_, err := tx.Exec("CREATE TABLE t(a, b);")
				return err
			},
			func(tx migration.LimitedTx) error {
				_, err := tx.Exec("CREATE TABLE t2(a, b);")
				return err
			},
		},
		func(tx migration.LimitedTx) (version int, err error) {
			err = tx.QueryRow("PRAGMA user_version;").Scan(&version)
			if err != nil {
				return version, err
			}
			return version, nil
		},
		func(tx migration.LimitedTx, version int) error {
			// It is fine to Printf here since we take an int. And pragma
			// can't be prepared, so there is no other option
			_, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", version))
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO t(a, b) VALUES (1, 2);")
	if err != nil {
		log.Fatal(err)
	}
}
