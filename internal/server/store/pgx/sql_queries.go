package pgx

const (
	metricsCreateQuery = "insert into metrics (id, type, delta, value, key) values ($1, $2, $3, $4, $5) " +
		"on conflict (id) do update set id = excluded.id, type = excluded.type, delta = excluded.delta, value = excluded.value, key = excluded.key"
	metricsGetByKeyQuery       = "select id, type, delta, value FROM metrics WHERE key in (?);"
	metricsGetAllQuery         = "SELECT id, type, delta, value from metrics "
	metricsGetByIdAndTypeQuery = "SELECT id, type, delta, value from metrics where id = $1 and type = $2 "
)
