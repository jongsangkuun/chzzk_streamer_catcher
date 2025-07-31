# 기본 환경 설정
ENV ?= dev
PROJECT_NAME ?= chzzk_streamer_catcher

# 환경별 빌드
service-build:
ifeq ($(ENV), dev)
	ENV=$(ENV) docker compose build
else ifeq ($(ENV), prod)
	ENV=$(ENV) docker compose -f docker-compose.yml -f docker-compose.prod.yml build
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

# 환경별 서비스 실행
service-up:
ifeq ($(ENV), dev)
	ENV=$(ENV) docker compose up -d
else ifeq ($(ENV), prod)
	ENV=$(ENV) docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

# 환경별 서비스 중지
service-down:
ifeq ($(ENV), dev)
	docker compose down
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

# 환경별 서비스 정리 (볼륨 포함)
service-clean:
ifeq ($(ENV), dev)
	docker compose down -v --remove-orphans
else ifeq ($(ENV), prod)
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down -v --remove-orphans
else
	@echo "Error: Unknown ENV value '$(ENV)'"
	@echo "Please set ENV to 'dev' or 'prod'"
	exit 1
endif

# 개별 앱 컨테이너 제어
app-up:
	docker start $(PROJECT_NAME)-app-1 || docker start $(PROJECT_NAME)_app_1

app-down:
	docker stop $(PROJECT_NAME)-app-1 || docker stop $(PROJECT_NAME)_app_1

app-restart:
	docker restart $(PROJECT_NAME)-app-1 || docker restart $(PROJECT_NAME)_app_1


# 헬프
help:
	@echo "사용 가능한 명령어:"
	@echo "  ENV=dev|prod make service-build    - 환경별 빌드"
	@echo "  ENV=dev|prod make service-up       - 환경별 서비스 실행"
	@echo "  ENV=dev|prod make service-down     - 환경별 서비스 중지"
	@echo "  ENV=dev|prod make service-clean    - 환경별 서비스 정리"
