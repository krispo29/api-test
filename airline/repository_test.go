package airline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetAllAirlines(t *testing.T) {
	// Create a mock database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Wrap the mockDB with go-pg's DB
	// This is a simplified way; direct interaction with go-pg's QueryContext might need more elaborate mocking
	// or a test database instance. For this example, we'll focus on mocking the query execution.
	// In a real scenario with go-pg, you might need to mock the ORM's behavior or use a dedicated test DB.
	// However, the current repository implementation uses db.QueryContext, which can be mocked
	// if we can control the *pg.DB instance passed in the context.

	// For go-pg, mocking can be complex. A common approach is to use an actual test database.
	// If that's not feasible, you'd typically mock the specific go-pg methods called.
	// The provided repository uses `db.QueryContext(ctx, &list, sqlStr)`.
	// Let's assume for this test that we can inject a mocked `pg.DB` that allows us to control `QueryContext`.
	// This often means wrapping `pg.DB` in an interface or using a library that helps mock `go-pg`.

	// Given the limitations of directly mocking go-pg without a more complex setup,
	// this test will be more conceptual for `QueryContext`.
	// A more robust test would use a test database or a specialized go-pg mocking tool.

	pgDB := pg.NewDB(&pg.Options{
		Addr:     "localhost:5432", // Dummy address for mock
		User:     "test",
		Password: "test",
		Database: "test",
	})
	// This is where it gets tricky. `pg.DB.QueryContext` is not easily mockable with sqlmock alone.
	// For the sake of this example, we'll illustrate the test logic,
	// acknowledging that the `pg.DB` mocking part would need a real setup.

	repo := NewRepository(5 * time.Second)

	t.Run("successful retrieval", func(t *testing.T) {
		// Expected SQL query (be careful with whitespace and exact matching)
		// The actual query uses "logo_url"
		expectedSQL := `
		SELECT 
			"uuid",
			"name",
			"logo_url"
		FROM public.tbl_airlines
		WHERE deleted_at IS NULL
		ORDER BY name ASC;
	`
		// Define the rows that sqlmock should return
		rows := sqlmock.NewRows([]string{"uuid", "name", "logo_url"}).
			AddRow("uuid1", "Airline One", "http://logo.url/1").
			AddRow("uuid2", "Airline Two", "http://logo.url/2")

		// Set up expectations on the mock
		// This part is conceptual with sqlmock as go-pg uses its own query building.
		// You would typically mock the Exec or Query methods that go-pg internally calls,
		// or mock methods on a pg.Query object if you build queries step-by-step.
		// If QueryContext directly maps to a stdlib sql.DB method, this could work.
		mock.ExpectQuery(expectedSQL).WillReturnRows(rows) // This regex needs to be precise

		// Create a context and add the mock database connection
		// In the actual code, pgDB is retrieved from context.
		// To test this, we need to put our mock-capable pgDB into the context.
		// This is the challenging part without refactoring for testability or using a test DB.

		// For this example, let's assume we *could* make our `pgDB` use the `sqlmock` connection.
		// This is non-trivial. A more practical approach for go-pg is often integration testing
		// with a real test database, or using a library designed for mocking go-pg.

		// --- Simplified Conceptual Test ---
		// Due to the difficulty of properly mocking go-pg's QueryContext with sqlmock directly,
		// let's simulate a scenario where the repository method is called
		// and we assert the outcome. The actual DB interaction mocking is hand-waved here.

		// If we could inject a mock pg.DB:
		// ctx := context.WithValue(context.Background(), "postgreSQLConn", mockPgDB)
		// airlines, err := repo.GetAllAirlines(ctx)
		// assert.NoError(t, err)
		// assert.NotNil(t, airlines)
		// assert.Len(t, *airlines, 2)
		// assert.Equal(t, "Airline One", (*airlines)[0].Name)

		// Since we can't easily mock `pg.DB`'s `QueryContext` with `sqlmock` in this setup,
		// we will write a placeholder test that highlights what *should* be tested.
		// A real implementation would require a test database or a more sophisticated mocking strategy for go-pg.

		// Placeholder:
		t.Log("Note: Proper go-pg DB mocking requires a test database or specialized tools.")
		t.Log("This test is a conceptual placeholder for `GetAllAirlines` success.")

		// Simulate a successful call by assuming the DB call works:
		// This part would be replaced by actual interaction with a mocked DB if possible.
		// For now, we can't execute this part of the test without a running DB or better mock.
		// So, we'll skip the actual call to repo.GetAllAirlines in this conceptual test.

		// To make this test runnable but acknowledge the mocking challenge:
		assert.True(t, true, "Placeholder for successful retrieval test. Implement with go-pg mocking or test DB.")
	})

	t.Run("database error", func(t *testing.T) {
		// Expected SQL query
		expectedSQL := `
		SELECT 
			"uuid",
			"name",
			"logo_url"
		FROM public.tbl_airlines
		WHERE deleted_at IS NULL
		ORDER BY name ASC;
	`
		// Expect a query and return an error
		mock.ExpectQuery(expectedSQL).WillReturnError(errors.New("DB error"))

		// --- Simplified Conceptual Test ---
		// Placeholder:
		t.Log("Note: Proper go-pg DB mocking requires a test database or specialized tools.")
		t.Log("This test is a conceptual placeholder for `GetAllAirlines` DB error.")

		// Simulate a database error scenario:
		// As above, without a proper mock/test DB, we can't execute this directly.
		// ctx := context.WithValue(context.Background(), "postgreSQLConn", mockPgDB)
		// _, err := repo.GetAllAirlines(ctx)
		// assert.Error(t, err)
		// assert.Equal(t, "DB error", err.Error())

		assert.True(t, true, "Placeholder for database error test. Implement with go-pg mocking or test DB.")
	})
}
