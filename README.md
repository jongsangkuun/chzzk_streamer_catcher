## chzzk_streamer_catcher 프로젝트 README

# chzzk_streamer_catcher

치지직(Chzzk) 스트리머 정보를 수집하고 분석하여 다양한 통계 자료를 제공하는 사이드 프로젝트입니다.

## 1. 구현 목표

*   **실시간 시청자 수:** 각 스트리머의 실시간 시청자 수 현황 및 변화 추이 분석
*   **평균 방송 시간:** 스트리머별 평균 방송 시간 측정 및 활동성 지표 제공
*   **카테고리별 통계:** 게임, 토크, 음악 등 다양한 카테고리별 스트리머 활동 및 인기 분석
*   **기타:** 스트리머 관련 유용한 통계 데이터 수집 및 제공

## 2. 기술 스택

*   **언어:** Go (v1.24.1)
*   **데이터베이스:** PostgreSQL
*   **API 통신:** 치지직 OpenAPI
*   **컨테이너화:** Docker, Docker Compose
*   **로깅:** Logrus

## 3. 주요 기능

*   **치지직 OpenAPI 연동:** `https://openapi.chzzk.naver.com/open` 엔드포인트를 활용하여 채널, 카테고리, 라이브 방송 정보를 수집합니다. (관련 상수: `const.go` 참조)
*   **데이터 수집 및 저장:**
    *   `live_data` 테이블에 스트리머 및 방송 정보를 주기적으로 저장합니다.
    *   약 한 달간의 데이터를 축적하여 분석의 기반을 마련합니다.
*   **데이터 분석:** 수집된 데이터를 기반으로 SQL 쿼리를 작성하여 평균 시청자 수, 최대 시청자 수, 방송 시간 등 다양한 통계 지표를 계산합니다.
*   **API 서버 구축:** 분석된 데이터를 외부에서 접근하고 활용할 수 있도록 RESTful API를 제공합니다.
*   **결과물 시각화:** (구현 예정) 수집 및 분석된 데이터를 사용자가 이해하기 쉬운 형태로 시각화하여 제공합니다.

## 4. 개발 환경 설정 및 실행

### 4.1. 사전 준비

*   **Docker & Docker Compose 설치:** 로컬 환경에서 애플리케이션을 실행하기 위해 Docker 및 Docker Compose를 설치합니다.
*   **Git 설치:** 소스 코드 관리를 위해 Git을 설치합니다.

### 4.2. 프로젝트 클론

```shell script
git clone https://github.com/jongsangkuun/chzzk_streamer_catcher.git
cd chzzk_streamer_catcher
```


### 4.3. 환경 변수 설정

`.env` 파일을 생성하고 다음 내용을 포함하여 필요한 환경 변수를 설정합니다.

```dotenv
# .env 파일 예시
CHZZK_CLIENT_ID=YOUR_CHZZK_CLIENT_ID
CHZZK_SECRET_ID=YOUR_CHZZK_SECRET_ID
POSTGRES_HOST=postgresql-db
POSTGRES_PORT=5432
TIMEZONE=Asia/Seoul
```


*   `CHZZK_CLIENT_ID`, `CHZZK_SECRET_ID`: 치지직 OpenAPI 사용을 위한 인증 정보 (필요시 발급받아야 함)
*   `POSTGRES_HOST`, `POSTGRES_PORT`: PostgreSQL 데이터베이스 연결 정보
*   `TIMEZONE`: 애플리케이션 및 데이터베이스 타임존 설정

### 4.4. Docker Compose 실행

프로젝트 루트 디렉토리에서 다음 명령어를 실행하여 애플리케이션과 PostgreSQL 데이터베이스를 실행합니다.

```shell script
docker-compose up -d
```


*   `Dockerfile`을 사용하여 Go 애플리케이션 이미지를 빌드하고 컨테이너를 실행합니다. (개발, 빌드, 프로덕션 스테이지 포함)
*   `docker-compose.yml` 설정에 따라 PostgreSQL 데이터베이스 컨테이너가 함께 실행되며, `healthcheck`를 통해 데이터베이스 준비 상태를 확인합니다.

### 4.5. 애플리케이션 확인

애플리케이션이 정상적으로 실행되면, 설정된 환경 변수(`CHZZK_CLIENT_ID`, `CHZZK_SECRET_ID` 등)를 기반으로 치지직 OpenAPI에서 데이터를 수집하고 PostgreSQL 데이터베이스에 저장합니다. API 서버는 기본적으로 실행되어 데이터를 제공할 준비를 합니다.

## 5. 프로젝트 구조

```
chzzk_streamer_catcher/
├── .env                  # 환경 변수 설정 파일
├── .gitignore            # Git 무시 파일 설정
├── Dockerfile            # Docker 이미지 빌드 설정
├── docker-compose.yml    # Docker Compose 설정 파일
├── go.mod                # Go 모듈 의존성 관리
├── go.sum                # Go 모듈 체크섬
├── README.md             # 프로젝트 설명 파일
├── cmd/                  # 애플리케이션 진입점
│   └── main.go
├── internal/             # 내부 로직 및 유틸리티
│   ├── api/              # API 관련 로직
│   ├── database/         # 데이터베이스 관련 로직
│   └── service/          # 비즈니스 로직
├── pkg/                  # 외부에서 재사용 가능한 패키지
│   └── chzzk/            # 치지직 API 관련 클라이언트 및 상수
│       ├── const.go      # API 상수 정의
│       └── client.go     # API 클라이언트 구현
└── web/                  # 프론트엔드 관련 파일 (구현 예정)
```


## 6. 향후 계획

*   **데이터 시각화:** 프론트엔드 개발을 통해 수집된 데이터를 차트 등으로 시각화하여 제공합니다.
*   **고도화된 분석:** 스트리머 활동 패턴 분석, 시청자 참여도 분석 등 심층적인 분석 기능을 추가합니다.
*   **알림 기능:** 특정 스트리머의 방송 시작 시 알림을 제공하는 기능을 구현합니다.
*   **성능 최적화:** 데이터 처리 및 API 응답 속도를 개선합니다.

---

이 README 파일은 프로젝트의 현재 상태와 향후 계획을 간략하게 요약합니다. 개발 과정에서 내용이 업데이트될 수 있습니다.