package repository

type Repository struct {
	dynamoDBClient DBInterface
	tableName      string
	logger         LoggerInterface
}

func NewRepository(
	dynamoDBClient DBInterface,
	tableName string,
	logger LoggerInterface,
) *Repository {
	return &Repository{
		dynamoDBClient: dynamoDBClient,
		tableName:      tableName,
		logger:         logger,
	}
}
