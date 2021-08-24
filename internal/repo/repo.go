package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/ozoncp/ocp-remind-api/internal/models"
)

type RemindsRepo interface {
	Add(remind []models.Remind) error
	Describe(id uint64) (*models.Remind, error)
	List(limit, offset uint64) ([]models.Remind, error)
	Remove(id, user_id uint64) error
}

type RemindDBRepository struct {
	db *pgx.Conn
}

func (r *RemindDBRepository) Add(reminds []models.Remind) error {
	query := squirrel.Insert("reminds").
		Columns("remind_id", "user_id", "deadline", "message")

	for _, remind := range reminds {
		query = query.Values(remind.Id, remind.UserId, remind.Deadline, remind.Text)
	}
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RemindDBRepository) Describe(id uint64) (*models.Remind, error) {
	var remind models.Remind
	query := squirrel.Select("remind_id", "user_id", "deadline", "message").
		From("reminds").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(context.Background(), sql, args...).Scan(&remind.Id, &remind.UserId, &remind.Deadline, &remind.Text)
	if err != nil {
		return nil, err
	}

	return &remind, nil
}

func (r *RemindDBRepository) List(limit, offset uint64) ([]models.Remind, error) {
	sql, args, err := squirrel.Select("remind_id", "user_id", "deadline", "message").
		From("reminds").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(context.Background(), sql, args...)
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

func (r *RemindDBRepository) Remove(id, user_id uint64) error {
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

	_, err = r.db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewRemindDBRepository(conn *pgx.Conn) (*RemindDBRepository, error) {

	return &RemindDBRepository{db: conn}, nil
}
