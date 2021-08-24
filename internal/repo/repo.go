package repo

import (
	"context"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"

	"github.com/ozoncp/ocp-remind-api/internal/models"
)

type RemindsRepo interface {
	Add(ctx context.Context, remind []models.Remind) error
	Describe(ctx context.Context, id uint64) (*models.Remind, error)
	List(ctx context.Context, limit, offset uint64) ([]models.Remind, error)
	Remove(ctx context.Context, id, user_id uint64) error
}

type RemindDBRepository struct {
	db *pgx.Conn
}

func (r *RemindDBRepository) Add(ctx context.Context, reminds []models.Remind) error {
	query := squirrel.Insert("reminds").
		Columns("remind_id", "user_id", "deadline", "message")

	for _, remind := range reminds {
		query = query.Values(remind.Id, remind.UserId, remind.Deadline, remind.Text)
	}
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RemindDBRepository) Describe(ctx context.Context, id uint64) (*models.Remind, error) {
	var remind models.Remind
	query := squirrel.Select("remind_id", "user_id", "deadline", "message").
		From("reminds").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, sql, args...).Scan(&remind.Id, &remind.UserId, &remind.Deadline, &remind.Text)
	if err != nil {
		return nil, err
	}

	return &remind, nil
}

func (r *RemindDBRepository) List(ctx context.Context, limit, offset uint64) ([]models.Remind, error) {
	sql, args, err := squirrel.Select("remind_id", "user_id", "deadline", "message").
		From("reminds").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	lists := make([]models.Remind, 0, limit)
	for rows.Next() {
		var remind models.Remind
		err := rows.Scan(&remind.Id,
			&remind.UserId,
			&remind.Deadline,
			&remind.Text)
		if err != nil {
			return nil, err
		}
		lists = append(lists, remind)
	}
	return lists, err
}

func (r *RemindDBRepository) Remove(ctx context.Context, id, user_id uint64) error {
	sql, args, err := squirrel.Delete("reminds").
		Where(
			squirrel.And{
				squirrel.Eq{"remind_id": id},
				squirrel.Eq{"user_id": user_id},
			}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewRemindDBRepository() (*RemindDBRepository, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("REMINDS_DB_URL"))
	if err != nil {
		log.Err(err).Msg("Unable to connect to db")
		return nil, err
	}
	//defer func(conn *pgx.Conn, ctx context.Context) {
	//	err := conn.Close(ctx)
	//	if err != nil {
	//		log.Err(err).Msg("Unable to close connection to db")
	//	}
	//}(conn, context.Background())
	return &RemindDBRepository{db: conn}, nil
}
