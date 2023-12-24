up:
	docker compose up -d --build --force-recreate

down:
	docker compose down

run_server:
	cd server && go run . || cd -
