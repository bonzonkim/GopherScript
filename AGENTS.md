## Project Context
이 프로젝트는 Python, Shell Script를 Parameter로 받아서 Go 로 변환 후, 하나의 정적 바이너리를 빌드하는 목적을 가졌습니다.
기존 Python, Shell Script 로 작성된 운영,관리 스크립트를 Go 정적 바이너리로 변환 것이 목표입니다.

## Tech Stack
- Language: Go (1.25.1)

## Coding Guidelines
1. **Error Handling**: 모든 에러는 래핑하여 스택 트레이스를 유지하십시오.
2. **Testing**: 비즈니스 로직 변경 시 반드시 유닛 테스트를 함께 업데이트하거나 새로 작성하십시오.
3. **Naming**: 인터페이스는 `er`로 끝나도록 명명하고, 구현체는 구체적인 이름을 사용하십시오.

## Agent Behavior Rules
- 파일을 수정하기 전에 항상 전체 프로젝트 구조를 먼저 파악하십시오.
- 실행 불가능한 코드를 생성하지 않도록, 코드 생성 후에는 항상 빌드 테스트를 수행하는 계획을 세우십시오.
- 복잡한 리팩토링을 수행할 때는 한 번에 모든 것을 바꾸지 말고, 단계별(Step-by-step) 계획을 먼저 제시하십시오.
