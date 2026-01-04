# GopherScript

<div align=center>

[한국어](README.ko.md) | [English](README.md) 

</div>

![GopherScript Logo](gopherScript-logo.png)

Python, Shell 스크립트를 Go 정적 바이너리로 변환하는 CLI 도구입니다.

## 개요

GopherScript는 LLM(Large Language Model)을 활용하여 Python 또는 Shell 스크립트를 관용적인(idiomatic) Go 코드로 변환합니다. 변환된 코드는 단일 정적 바이너리로 컴파일되어 별도의 런타임 의존성 없이 어디서든 실행할 수 있습니다.

### 지원 LLM 프로바이더
- **Google Gemini** (기본값)
- **OpenAI GPT-4o**
- **Anthropic Claude**

## 설치

### 릴리스에서 다운로드 (권장)

[Releases 페이지](https://github.com/bonzonkim/GopherScript/releases)에서 플랫폼에 맞는 최신 릴리스를 다운로드하세요.

**Linux/macOS (한 줄 설치):**
```bash
# Linux (amd64)
curl -sL https://github.com/bonzonkim/GopherScript/releases/latest/download/gopherscript_linux_amd64.tar.gz | tar xz
sudo mv gopherscript /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/bonzonkim/GopherScript/releases/latest/download/gopherscript_darwin_arm64.tar.gz | tar xz
sudo mv gopherscript /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/bonzonkim/GopherScript/releases/latest/download/gopherscript_darwin_amd64.tar.gz | tar xz
sudo mv gopherscript /usr/local/bin/
```

**Windows:**
1. [Releases](https://github.com/bonzonkim/GopherScript/releases)에서 `gopherscript_windows_amd64.zip` 다운로드
2. 압축 해제 후 PATH에 추가

### Go Install 사용

```bash
go install github.com/bonzonkim/GopherScript@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/bonzonkim/GopherScript.git
cd GopherScript
go build -o gopherscript .
```

## 사용법

### 기본 사용법

```bash
# Python 스크립트를 Go로 변환
gopherscript script.py

# Shell 스크립트를 Go로 변환
gopherscript script.sh

# 출력 경로 지정
gopherscript script.py -o output.go

# 변환 후 바이너리 빌드
gopherscript script.py --build

# 커스텀 바이너리 경로로 빌드
gopherscript script.py --build -b ./bin/myapp
```

### LLM 프로바이더 선택

```bash
# OpenAI GPT 사용
gopherscript script.py --provider openai

# Anthropic Claude 사용
gopherscript script.py --provider claude

# Google Gemini 사용 (기본값)
gopherscript script.py --provider gemini
```

### 환경 변수

| 변수명 | 설명 |
|--------|------|
| `LLM_PROVIDER` | 기본 LLM 프로바이더 (gemini/openai/claude) |
| `GEMINI_API_KEY` | Google Gemini API 키 |
| `OPENAI_API_KEY` | OpenAI API 키 |
| `ANTHROPIC_API_KEY` | Anthropic Claude API 키 |
| `API_KEY` | (레거시) Gemini API 키로 폴백 |

### CLI 플래그

| 플래그 | 단축 | 설명 |
|--------|------|------|
| `--output` | `-o` | 생성될 Go 파일 경로 |
| `--binary` | `-b` | 컴파일될 바이너리 경로 (--build 필요) |
| `--build` | | 변환 후 바이너리 빌드 |
| `--provider` | `-p` | 사용할 LLM 프로바이더 |
| `--verbose` | `-v` | 상세 로깅 활성화 |

## ⚠️ 주의사항

### 🔐 민감한 정보 보안

> [!CAUTION]
> **스크립트에 포함된 민감한 정보는 LLM 서버로 전송됩니다!**

GopherScript는 스크립트 내용을 외부 LLM API로 전송하여 변환합니다. 따라서:

1. **API 키, 비밀번호, 토큰 마스킹**
   ```bash
   # ❌ 위험: 실제 API 키가 노출됨
   API_KEY="sk-1234567890abcdef"
   
   # ✅ 안전: 플레이스홀더 사용
   API_KEY="${API_KEY}"  # 또는 API_KEY="YOUR_API_KEY_HERE"
   ```

2. **데이터베이스 연결 정보**
   ```bash
   # ❌ 위험
   DB_PASSWORD="mysecretpassword"
   
   # ✅ 안전
   DB_PASSWORD="${DB_PASSWORD}"
   ```

3. **내부 서버 주소**
   ```bash
   # ❌ 위험: 내부 인프라 정보 노출
   curl http://internal-server.company.com:8080
   
   # ✅ 안전
   curl "${INTERNAL_SERVER_URL}"
   ```

### 📋 변환 전 체크리스트

- [ ] 하드코딩된 비밀번호/API 키를 환경 변수로 교체
- [ ] 내부 IP 주소 및 도메인을 변수화
- [ ] 개인정보(PII)가 포함되어 있지 않은지 확인
- [ ] 회사 기밀 정보가 없는지 확인

### 🔍 자동 마스킹 (권장)

변환 전에 민감한 정보를 마스킹하는 스크립트 예시:

```bash
# 변환 전 마스킹 처리
sed -e 's/password="[^"]*"/password="${PASSWORD}"/g' \
    -e 's/api_key="[^"]*"/api_key="${API_KEY}"/g' \
    script.py > script_masked.py

# 마스킹된 스크립트로 변환
gopherscript script_masked.py
```

### ⚡ 기타 주의사항

1. **LLM 출력 검증**: 생성된 Go 코드는 반드시 검토하세요. LLM이 원본 로직을 완벽하게 변환하지 못할 수 있습니다.

2. **복잡한 스크립트**: 매우 복잡한 스크립트는 여러 번의 시도가 필요할 수 있습니다.

3. **시스템 의존성**: 특정 시스템 명령어나 라이브러리에 의존하는 스크립트는 수동 수정이 필요할 수 있습니다.

4. **API 비용**: LLM API 호출에는 비용이 발생합니다. 대용량 스크립트나 빈번한 변환 시 비용을 확인하세요.

## 예제

### Python 스크립트 변환

원본 (`example.py`):
```python
#!/usr/bin/env python3
import sys

def greet(name):
    return f"Hello, {name}!"

if __name__ == "__main__":
    name = sys.argv[1] if len(sys.argv) > 1 else "World"
    print(greet(name))
```

변환:
```bash
gopherscript example.py --build
```

### Shell 스크립트 변환

원본 (`example.sh`):
```bash
#!/bin/bash
NAME=${1:-"World"}
echo "Hello, $NAME!"
```

변환:
```bash
gopherscript example.sh -o hello.go --build -b hello
```

## 라이선스

MIT License

## 기여

이슈와 풀 리퀘스트를 환영합니다!
