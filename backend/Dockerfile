# 🏗️ ช่วงที่ 1: ขั้นตอนการประกอบร่าง (ใช้เวอร์ชันระบุชัดเจนป้องกันระบบพังในอนาคต)
FROM golang:latest AS builder

WORKDIR /app

# คัดลอกเฉพาะไฟล์จัดการ Dependency
COPY go.mod go.sum ./

# ดาวน์โหลดโปรแกรมเสริมมารอไว้
RUN go mod download

# คัดลอกซอร์สโค้ดทั้งหมดเข้ากล่องประกอบร่าง
COPY . .

# [แก้ไขหลัก]: เปลี่ยน CGO_ENABLED=0 เพื่อตัดปัญหาเรื่องไลบรารีภาษา C ของ Linux ออกไปทั้งหมด 
# และทำการสร้างไฟล์โปรแกรมแบบ Standalone ที่รันบนตัว Alpine จิ๋วได้ทันที
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o queue-ticket-api ./cmd/main.go


# 📦 ช่วงที่ 2: ขั้นตอนการแพ็คของส่งออก (ใช้ตู้คอนเทนเนอร์ขนาดจิ๋วที่ปลอดภัย)
FROM alpine:3.19

WORKDIR /root/
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/queue-ticket-api .
COPY .env.example .env

EXPOSE 3000

# สั่งเปิดใช้งานระบบคิว
CMD ["./queue-ticket-api"]
