package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"context"
	"time"
)

const (
	kUser = "user"
)

type Recommendation struct {
	UserId         int64
	Recommendation string
	SendAt         time.Time
}

type Reminder struct {
	UserId   int64
	SendAt   time.Time
	Reminder string
}

type Ration struct {
	Ration    string
	CreatedAt time.Time
}

func execWrapper(conn *pgxpool.Conn, query string, args ...interface{}) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	defer conn.Release()

	commTag, err := conn.Exec(ctx, query, args...)
	return commTag, err
}

func queryRowWrapper(conn *pgxpool.Conn, query string, args ...interface{}) pgx.Row {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	defer conn.Release()

	return conn.QueryRow(ctx, query, args...)
}

func queryWrapper(conn *pgxpool.Conn, query string, args ...interface{}) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	defer conn.Release()

	rows, err := conn.Query(ctx, query, args...)
	return rows, err
}

func AddRecommendation(conn *pgxpool.Conn, userId int64, recommendation string) error {
	query := `
		INSERT INTO recommendations (user_id, recommendation, send_at)
		SELECT $1, $2, $3
		WHERE NOT EXISTS (
			SELECT 1 FROM banned_users
			WHERE banned_users.user_id = $1
				AND banned_users.banned_section = 'recommendation')`

	_, err := execWrapper(conn, query, userId, recommendation, time.Now())

	return err
}

func GetUserRoleByUserId(conn *pgxpool.Conn, userId int64) (string, error) {
	query := `
		SELECT role FROM users WHERE user_id = $1`
	
	var role string
	err := queryRowWrapper(conn, query, userId).Scan(&role)

	return role, err
}

func GetRecommendations(conn *pgxpool.Conn, limit int) ([]Recommendation, error) {
	query := `
		SELECT user_id, recommendation, send_at FROM recommendations
		ORDER BY send_at DESC LIMIT $1`
	
	rows, err := queryWrapper(conn, query, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var recommendations []Recommendation
	for rows.Next() {
		var rec Recommendation
		if err := rows.Scan(&rec.UserId, &rec.Recommendation, &rec.SendAt); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, rec)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return recommendations, nil
}

func BanUser(conn *pgxpool.Conn, userId int64, reason string, banned_section string) error {
	query := `
		INSERT INTO banned_users (user_id, reason, banned_at, banned_section)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			reason = EXCLUDED.reason,
			banned_at = EXCLUDED.banned_at,
			banned_section = EXCLUDED.banned_section`
	
	_, err := execWrapper(conn, query, userId, reason, time.Now(), banned_section)

	return err
}

func IsUserBannedAll(conn *pgxpool.Conn, userId int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM banned_users
			WHERE user_id = $1 AND banned_section = 'all')`
	
	var exists bool
	err := queryRowWrapper(conn, query, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func AssignRole(conn *pgxpool.Conn, userId int64, role string) error {
	query := `
		UPDATE users SET role = $2 WHERE user_id = $1`

	_, err := execWrapper(conn, query, userId, role)

	return err
}

func Registration(conn *pgxpool.Conn, userId int64, timezone string) error {
	query := `
		INSERT INTO users (user_id, timezone, registration_at, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			timezone = EXCLUDED.timezone,
			registration_at = EXCLUDED.registration_at`
	
	_, err := execWrapper(conn, query, userId, timezone, time.Now(), kUser)

	return err
}

func IsRegistrated(conn *pgxpool.Conn, userId int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE user_id = $1)`
	
	var exists bool
	err := queryRowWrapper(conn, query, userId).Scan(&exists)
	return exists, err
}

func GetTimezoneByPrimaryKey(conn *pgxpool.Conn, userId int64) (string, error) {
	query := `
		SELECT timezone FROM users
		WHERE user_id = $1`
	
	var timezone string
	err := queryRowWrapper(conn, query, userId).Scan(&timezone)

	return timezone, err
}

func InsertReminder(conn *pgxpool.Conn, userId int64, sendAt time.Time, expireAt time.Time, reminder string) error {
	query := `
		INSERT INTO reminders (user_id, send_at, expire_at, reminder)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, send_at)
		DO UPDATE SET
			reminder = EXCLUDED.reminder`

	_, err := execWrapper(conn, query, userId, sendAt, expireAt, reminder)

	return err
}

func GetActualReminders(conn *pgxpool.Conn) ([]Reminder, error) {
	now := time.Now()
	query := `
		SELECT user_id, send_at, reminder FROM reminders
		WHERE $1 >= send_at AND $1 < expire_at`
	
	rows, err := queryWrapper(conn, query, now)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var rem Reminder
		if err := rows.Scan(&rem.UserId, &rem.SendAt, &rem.Reminder); err != nil {
			return nil, err
		}
		reminders = append(reminders, rem)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return reminders, nil
}

func DeleteSentReminders(conn *pgxpool.Conn, reminders []Reminder) error {
	userIds := make([]int64, 0, len(reminders))
	sendAts := make([]time.Time, 0, len(reminders))
	for _, reminder := range reminders {
		userIds = append(userIds, reminder.UserId)
		sendAts = append(sendAts, reminder.SendAt)
	}

	query := `
		DELETE FROM reminders
		WHERE (user_id, send_at) IN (
		SELECT * FROM unnest($1::bigint[], $2::timestamptz[]) AS t(user_id, send_at))`

	_, err := execWrapper(conn, query, userIds, sendAts)
	return err
}

func DeleteExpiredReminders(conn *pgxpool.Conn) error {
	now := time.Now()

	query := `
		DELETE FROM reminders
		WHERE $1 >= expire_at`

	_, err := execWrapper(conn, query, now)
	return err
}

func GetActualRemindersByUserId(conn *pgxpool.Conn, userId int64) ([]Reminder, error) {
	now := time.Now()
	query := `
		SELECT user_id, send_at, reminder FROM reminders
		WHERE $1 = user_id AND $2 <= send_at`
	
	rows, err := queryWrapper(conn, query, userId, now)

	if err != nil {
		return nil, err
	}

	var reminders []Reminder
	for rows.Next() {
		var rem Reminder
		if err := rows.Scan(&rem.UserId, &rem.SendAt, &rem.Reminder); err != nil {
			return nil, err
		}
		reminders = append(reminders, rem)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return reminders, nil
}

func GetGoal(conn *pgxpool.Conn, userId int64) (int, error) {
	query := `
		SELECT goal FROM goals
		WHERE $1 = user_id`

	var goal int
	err := queryRowWrapper(conn, query, userId).Scan(&goal)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return goal, nil
}

func GetUserRation(conn *pgxpool.Conn, userId int64, date time.Time) ([]string, error) {
	query := `
		SELECT ration FROM rations
		WHERE user_id = $1 AND created_at::date = $2::date`

	rows, err := queryWrapper(conn, query, userId, date)
	if err != nil {
		return nil, err
	}

	var rations []string
	for rows.Next() {
		var ration string
		if err := rows.Scan(&ration); err != nil {
			return nil, err
		}
		rations = append(rations, ration)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return rations, nil
}

func InsertRation(conn *pgxpool.Conn, userId int64, ration string, created_at time.Time) error {
	query := `
		INSERT INTO rations (user_id, ration, created_at)
		VALUES ($1, $2, $3)`
	
	_, err := execWrapper(conn, query, userId, ration, created_at)

	return err
}

func InsertGoal(conn *pgxpool.Conn, userId int64, goal int) error {
	query := `
		INSERT INTO goals (user_id, goal)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			goal = EXCLUDED.goal`

	_, err := execWrapper(conn, query, userId, goal)
	return err
}

func GetRationsForTheLastTime(conn *pgxpool.Conn, userId int64, startPeriod time.Time) ([]Ration, error) {
	query := `
		SELECT ration, created_at FROM rations
		WHERE user_id = $1 AND created_at > $2
		ORDER BY created_at`

	rows, err := queryWrapper(conn, query, userId, startPeriod)
	if err != nil {
		return nil, err
	}

	var rations []Ration
	for rows.Next() {
		var ration Ration
		if err := rows.Scan(&ration.Ration, &ration.CreatedAt); err != nil {
			return nil, err
		}
		rations = append(rations, ration)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return rations, err
}
