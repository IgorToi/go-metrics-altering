package psql

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// var db *sql.DB

// func TestMain(m *testing.M) {
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("could not constuct pool: %s", err)
// 	}

// 	err = pool.Client.Ping()
// 	if err != nil {
// 		log.Fatalf("could not connect to Docker: %s", err)
// 	}

// 	// pulls an image, creates a container based on it and runs it
// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "11",
// 		Env: []string{
// 			"POSTGRES_PASSWORD=secret",
// 			"POSTGRES_USER=user_name",
// 			"POSTGRES_DB=dbname",
// 			"listen_addresses = '*'",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		// set AutoRemove to true so that stopped container goes away by itself
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		log.Fatalf("could not start resource: %s", err)
// 	}

// 	hostAndPort := resource.GetHostPort("5432/tcp")
// 	databaseUrl := fmt.Sprintf("posrgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

// 	log.Println("Connecting to database on url: ", databaseUrl)

// 	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

// 	// exponentioal backoff-retry, because the application in the container might not be ready to accept connections yet
// 	pool.MaxWait = 120 * time.Second
// 	if err = pool.Retry(func() error {
// 		db, err = sql.Open("postgres", databaseUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	// Migrating DB
// 	if err := runMigrations("../../migrations", db); err != nil {
// 		log.Fatalf("Could not migrate db: %s", err)
// 	}

// 	code := m.Run()

// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("could not purge resource: %s", err)
// 	}

// 	os.Exit(code)
// }

// func TestDatabase(t *testing.T) {
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()

//     rows, err := db.QueryContext(ctx, "SELECT 1")
//     assert.NoError(t, err)
//     assert.True(t, rows.Next())
// }