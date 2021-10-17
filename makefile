run:
	cd services/auth/keys &&  docker-compose up
	cp services/auth/keys/public.key services/backend/keys
	docker-compose up -d
