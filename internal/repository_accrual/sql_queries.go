package repositoryaccrual

const (
	queryExistOrderID = `
		SELECT 1 FROM orders WHERE order_number = $1 LIMIT 1;
	`
	querySetOrdeData = `
		INSERT INTO orders (order_number, description, price)
		VALUES ($1, $2, $3)
	`
	querySetOrder = `
		"INSERT INTO accrual_orders (id, status) 
		VALUES ($1, $2)"
	`

	queryGetOrderForProcessing = `
		SELECT order_number, status, accrual 
		FROM accrual_orders 
		WHERE status=$1 order by id;
	`
	queryChangeStatusOrder = `
		UPDATE accrual_orders 
		SET status=$1 
		WHERE id=$2
	`

	queryGetReward = `
		SELECT reward, reward_type 
		FROM catalog 
		WHERE $1 ILIKE '%' || match || '%' limit 1",
	`

	queryAddAccrual = `
		UPDATE accrual_orders 
		SET accrual=$1, status=$2 
		WHERE id=$3
	`
)
