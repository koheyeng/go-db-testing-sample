package dbtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	dbConn   *gorm.DB
	resource *dockertest.Resource
	pool     *dockertest.Pool
)

func SetTestData(db *gorm.DB, file string, v interface{}) error {
	// Read test data files
	g, err := ioutil.ReadFile(filepath.Join("testdata", file+".json"))
	if err != nil {
		return fmt.Errorf("failed reading %s.json: %s", file, err)
	}
	if err := json.Unmarshal(g, v); err != nil {
		return fmt.Errorf("%s.json cannot unmarshal to struct", file)
	}

	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.Len(); i++ {
		if err := db.Create(rv.Index(i).Interface()).Error; err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}

func SetUpTestContainerPostgres(args []string) (*gorm.DB, func()) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not creates a new pool: %s", err)
	}

	address := os.Getenv("UT_ATTRDB_ADDRESS")
	if address != "docker" {
		address = "localhost"
	}

	opts := &dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "latest",
		Env:          []string{"POSTGRES_PASSWORD=password"},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: address, HostPort: "5433"},
			},
		},
	}

	resource, err = pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		connStr := fmt.Sprintf("host=%s port=%s user=%s sslmode=%s password=%s", address, "5433", "postgres", "disable", "password")
		dbConn, err = gorm.Open("postgres", connStr)
		if err != nil {
			return err
		}

		dbConn.SingularTable(true)
		dbConn.LogMode(true)

		return dbConn.DB().Ping()
	}); err != nil {
		tearDownTestContainer()
		log.Fatalf("Could not connect to docker: %s", err)
	}

	createTestTable(dbConn, args)

	return dbConn, tearDownTestContainer

}

func tearDownTestContainer() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func createTestTable(db *gorm.DB, args []string) {
	for _, arg := range args {
		f, err := os.Open(filepath.Join("./testdata/ddl", arg+".sql"))
		if err != nil {
			log.Fatalf("Could not open %s table create query: %s\n", arg, err)
		}

		content, _ := ioutil.ReadAll(f)
		db.Exec(string(content))
		f.Close()
	}
}
