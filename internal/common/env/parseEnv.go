package env

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	PsqlHost     string
	PsqlPort     string
	PsqlUser     string
	PsqlPassword string
	PsqlDb       string

	ChzzkClientId string
	ChzzkSecretId string
}

func newEnv() Env {
	return Env{}
}

func ParseEnv() (Env, error) {
	env := newEnv()

	// 환경 변수 매핑 (환경변수명 -> 구조체 필드 포인터, 필수여부)
	envMappings := map[string]struct {
		field    *string
		required bool
	}{
		"POSTGRES_HOST":     {&env.PsqlHost, true},
		"POSTGRES_PORT":     {&env.PsqlPort, true},
		"POSTGRES_USER":     {&env.PsqlUser, true},
		"POSTGRES_PASSWORD": {&env.PsqlPassword, true},
		"POSTGRES_DB":       {&env.PsqlDb, true},

		"CHZZK_CLIENT_ID": {&env.ChzzkClientId, true},
		"CHZZK_SECRET_ID": {&env.ChzzkSecretId, true},
	}

	var missingVars []string

	// 모든 환경 변수 처리
	for envName, mapping := range envMappings {
		value := os.Getenv(envName)
		*mapping.field = value

		// 필수 값이 비어있는지 확인
		if mapping.required && value == "" {
			missingVars = append(missingVars, envName)
		}
	}

	// 누락된 환경 변수가 있으면 에러 반환
	if len(missingVars) > 0 {
		return env, fmt.Errorf("다음 필수 환경 변수들이 설정되지 않았습니다: %s",
			strings.Join(missingVars, ", "))
	}

	return env, nil
}
