package user

const (
	insertUser = `
		INSERT INTO users
		(
			fk_role_id,
		 	name,
		 	email,
		 	password,
		 	created_at,
		 	created_by
		)
		VALUES
		(
			:fk_role_id,
		 	:name,
		 	:email,
		 	:password,
		 	:created_at,
		 	:created_by
		) RETURNING *
	`

	readUser = `
		SELECT
		    id,
			fk_role_id,
		 	name,
		 	email,
		 	password,
			status,
			flag,
			meta,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			users
	`

	countUser = `
		SELECT
			COUNT(*)
		FROM
			users
	`

	updateUser = `
		UPDATE
			users
	`
)
