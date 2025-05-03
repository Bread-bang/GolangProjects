# Ping

## Introduction

- ping이란?

      상대 컴퓨터로 ICMP Echo Request를 전달하여, 네트워크 상태를 판별하는 명령어이다.

  참고

  *https://www.rfc-editor.org/rfc/rfc792.html*

  *https://pkg.go.dev/golang.org/x/net/icmp#Message*

- 특징

  1. TCP/UDP를 사용하지 않음
  2. IP 패킷의 payload는 ICMP message임

## Roadmap

### v1.0

> 단일 원격 주소에 대해 Ping 통신을 수행합니다.

#### 🛠️ 명령어

```
실행파일명 <목적지 IP>
```

### v2.0

> 다중 원격 주소와의 Ping 통신

#### 🛠️ 명령어

```
실행파일명 <목적지 IP...> -c [MESSAGE 개수]
```

#### 👀 비고

- min/avg/max/stddev 계산 추가

### v3.0

> raw socket으로 단일 원격 주소와의 Ping 통신

### v4.0

> raw socket으로 다중 원격 주소와의 Ping 통신
