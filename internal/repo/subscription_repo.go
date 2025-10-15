package repo

import (
	"database/sql"

	"github.com/teamcutter/subscriptions-service-task/internal/model"
	"github.com/teamcutter/subscriptions-service-task/internal/utils"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(s *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	startDate, err := utils.ParseDateFromRequest(s.StartDate)
	if err != nil {
		return err
	}

	var endDate interface{}
	if s.EndDate != "" {
		parsedEndDate, err := utils.ParseDateFromRequest(s.EndDate)
		if err != nil {
			return err
		}
		endDate = parsedEndDate
	} else {
		endDate = nil
	}

	_, err = r.db.Exec(
		query, 
		s.ServiceName, 
		s.Price, 
		s.UserID, 
		startDate, 
		endDate)

	return err
}

func (r *SubscriptionRepo) GetAll() ([]model.Subscription, error) {
	rows, err := r.db.Query(
		`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		sub.StartDate, err = utils.ParseDateFromDB(sub.StartDate)
		if err != nil {
			return nil, err
		}

		if sub.EndDate != "" {
			sub.EndDate, err = utils.ParseDateFromDB(sub.EndDate)
			if err != nil {
				return nil, err
			}
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

func (r *SubscriptionRepo) TotalCost(userID, serviceName string, start, end string) (int, error) {
	query := 
	`
	SELECT
	SUM(
    	price * (
        	DATE_PART('year', AGE(COALESCE(end_date, $4), GREATEST(start_date, $3))) * 12 +
        	DATE_PART('month', AGE(COALESCE(end_date, $4), GREATEST(start_date, $3))) + 1
    	)
	) as total_price
	FROM subscriptions
	WHERE user_id = $1
	AND ($2 = '' OR service_name = $2)
	AND start_date <= $4
	AND (end_date >= $3 OR end_date IS NULL);
	`
	var sum sql.NullInt64

	startDate, err := utils.ParseDateFromRequest(start)
	if err != nil {
		return 0, err
	}

	endDate, err := utils.ParseDateFromRequest(end)
	if err != nil {
		return 0, err
	}

	err = r.db.QueryRow(query, userID, serviceName, startDate, endDate).Scan(&sum)
	if sum.Valid {
		return int(sum.Int64), err
	}
	return 0, err
}