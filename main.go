package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sinmetal/spshovel/spanner"
)

type Param struct {
	Project     string
	Instance    string
	Database    string
	SqlFilePath string
}

func main() {
	param, err := getFlag()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	db := fmt.Sprintf("projects/%s/instances/%s/databases/%s", param.Project, param.Instance, param.Database)
	fmt.Println(db)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed get working dir. err=%+v\n", err)
	}
	fmt.Println(wd)

	sql, err := ReadSQL(param.SqlFilePath)
	if err != nil {
		fmt.Printf("failed read sql file. err=%+v\n", err)
		os.Exit(1)
	}
	fmt.Println(sql)

	ctx := context.Background()
	sc := spanner.NewClient(ctx, db)
	s := spanner.NewSpannerEntityService(sc)
	cn, data, err := s.Query(ctx, sql)
	if err != nil {
		fmt.Printf("failed query to spanner. err=%+v\n", err)
	}

	var records [][]string
	records = append(records, cn)
	for _, v := range data {
		records = append(records, v)
	}

	if err := Write(wd, records); err != nil {
		fmt.Printf("failed write file. err=%+v\n", err)
	}
}

func getFlag() (*Param, error) {
	var (
		project     = flag.String("project", "", "project is spanner project")
		instance    = flag.String("instance", "", "instance is spanner insntace")
		database    = flag.String("database", "", "database is spanner database")
		sqlFilePath = flag.String("sql-file-path", "", "sql-file-path is sql file path")
	)
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	var emsg string
	if len(*project) < 1 {
		emsg += "project is required\n"
	}
	if len(*instance) < 1 {
		emsg += "instance is required\n"
	}
	if len(*database) < 1 {
		emsg += "database is required\n"
	}
	if len(*sqlFilePath) < 1 {
		emsg += "sql-file-path is required\n"
	}

	if len(emsg) > 0 {
		return nil, errors.New(emsg)
	}

	return &Param{
		Project:     *project,
		Instance:    *instance,
		Database:    *database,
		SqlFilePath: *sqlFilePath,
	}, nil
}