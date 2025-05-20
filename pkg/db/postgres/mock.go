package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPool) Close() {
	m.Called()
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, arguments)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRow) Close() {
	m.Called()
}

func (m *MockRow) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRow) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func (m *MockRow) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRow) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRow) Values() ([]any, error) {
	args := m.Called()
	return args.Get(0).([]any), args.Error(1)
}

func (m *MockRow) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}
