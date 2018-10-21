package teacup

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type orderByPolarity string
type field string

const (
	DESC orderByPolarity = "DESC"
	ASC  orderByPolarity = "ASC"

	Uid              field = "uid"
	Optional_summary field = "optional_summary"
	Title            field = "title"
	Body             field = "body"
	Post_date        field = "post_date"
)

func SelectContentByUid(uid uint32, table string, dbInfo string) (*PageContents, error) {
	queryRow := func(db sql.DB) (*sql.Row, error) {
		return db.QueryRow(fmt.Sprintf(`
SELECT pa.uid, pa.summary, pa.title, pa.body, pa.post_date 
FROM %s pa
WHERE uid = $1;`, pq.QuoteIdentifier(table)), uid), nil
	}

	return selectContent(queryRow, dbInfo)
}

func SelectContentByTitle(title string, table string, dbInfo string) (*PageContents, error) {
	queryRow := func(db sql.DB) (*sql.Row, error) {
		return db.QueryRow(fmt.Sprintf(`
SELECT pa.uid, pa.summary, pa.title, pa.body, pa.post_date
FROM %s pa
WHERE title = $1;`, pq.QuoteIdentifier(table)), title), nil
	}

	return selectContent(queryRow, dbInfo)
}

func SelectMultipleContents(limit uint32, offset uint32, orderby field, polarity orderByPolarity, table string, dbInfo string) ([]*PageContents, error) {
	queryRows := func(db sql.DB) (*sql.Rows, error) {
		queryStr := fmt.Sprintf(`
SELECT pa.uid, pa.summary, pa.title, pa.body, pa.post_date
FROM %s pa
ORDER BY %s %s
OFFSET $1 LIMIT $2;`, pq.QuoteIdentifier(table), pq.QuoteIdentifier(string(orderby)), string(polarity)) //Don't follow my example, kids.
		return db.Query(queryStr, offset, limit)
	}

	return selectContents(queryRows, dbInfo)
}

func selectContent(queryRow func(db sql.DB) (*sql.Row, error), dbInfo string) (*PageContents, error) {
	db, _ := sql.Open("postgres", dbInfo)
	defer db.Close()

	var p PageContents

	row, err := queryRow(*db)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&p.Uid, &p.Summary, &p.Title, &p.Body, &p.PostDate)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func selectContents(queryRows func(db sql.DB) (*sql.Rows, error), dbInfo string) ([]*PageContents, error) {
	db, _ := sql.Open("postgres", dbInfo)
	defer db.Close()

	rows, err := queryRows(*db)
	if err != nil {
		return nil, err
	}

	var contentArray []*PageContents

	for rows.Next() {
		var p PageContents
		err = rows.Scan(&p.Uid, &p.Summary, &p.Title, &p.Body, &p.PostDate)
		if err != nil {
			return nil, err
		}

		contentArray = append(contentArray, &p)
	}

	return contentArray, nil
}

func (t *teacup) CreateTable(name string, uniqueTitles bool) {
	db, _ := sql.Open("postgres", t.DbInfo)
	defer db.Close()

	unique := ""
	if uniqueTitles {
		unique = "UNIQUE"
	}

	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS ` + name + ` (
		uid serial PRIMARY KEY,
		optional_summary varchar(192),
		title varchar(128) ` + unique + ` NOT NULL,
		body text NOT NULL,
		post_date date DEFAULT CURRENT_DATE NOT NULL
	) WITH (OIDS=FALSE)`)

	if err != nil {
		t.Log.Fatal(err, "\nTeacup could not connect to postgresql database.")
	}

	_, err = db.Exec(`
-- Summary function for body.
CREATE OR REPLACE FUNCTION summary(rec ` + name + `)
  RETURNS varchar(192)
LANGUAGE SQL
AS
$$
SELECT
       CASE WHEN $1.optional_summary IS NULL
                 THEN
           CASE WHEN length($1.body) > 192
                     THEN
               $1.body::varchar(191) || '…'
                ELSE
               $1.body::varchar(192)
               END
            ELSE
           $1.optional_summary
           END
$$;`)

	if err != nil {
		t.Log.Fatal(err)
	}

	t.tables[name] = uniqueTitles
	//TODO("check for duplicates")
}
