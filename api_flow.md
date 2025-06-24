# การทำงานของ API Endpoint `/compare` - คู่มือละเอียด

การทำงานของ API endpoint `/compare` แบบละเอียดทีละขั้นตอน โดยใช้การเปรียบเทียบกับร้านอาหาร:

**[ลูกค้า (Client)] → [เสิร์ฟ (Handler)] → [หัวหน้าครัว (Service)] → [พ่อครัว (Repository)] → [คลังสินค้า (Database)]**

---

## 📱 **[ลูกค้า (Client)] - เริ่มต้นการสั่งอาหาร (ส่งคำขอ)**

ลูกค้า (เช่น แอปพลิเคชันหน้าบ้าน หรือเครื่องมืออย่าง Postman) ส่งคำขอ HTTP POST ไปยัง `/v1/compare/compare`

### ข้อมูลที่ต้องส่ง:
คำขอนี้เป็นแบบ `multipart/form-data` และต้องมี:
- **`excelFile`**: ไฟล์ Excel (.xlsx) ที่ต้องการประมวลผล
- **`columnName`**: ชื่อคอลัมน์ในไฟล์ Excel ที่จะใช้เปรียบเทียบ 
  - ตัวอย่าง: "goods_en", "goods_th", "hs_code"

### การยืนยันตัวตน:
ลูกค้าต้องส่ง **JWT token** ใน Authorization header เนื่องจากเป็น endpoint ที่มีการป้องกัน

---

## 🍽️ **[เสิร์ฟ (Handler: server/compare.go)] - รับออเดอร์และตรวจสอบเบื้องต้น**

### การกำหนดเส้นทางและ Middleware:
เซิร์ฟเวอร์หลัก (`server/server.go`) จะส่งคำขอไปยัง `excelHandler` โดยผ่าน middleware มาตรฐาน:
- Logging
- Request ID
- Timeout
- JWT Authentication
- การเติมข้อมูลบริบทด้วยการเชื่อมต่อฐานข้อมูล

### วิธีการ CompareExcel:

#### 1. Parse Form Data
- แยกวิเคราะห์ข้อมูล multipart form (`r.ParseMultipartForm`)
- ขนาดไฟล์สูงสุด: **10MB** (10 << 20 bytes)

#### 2. ดึงข้อมูล
- **ดึงไฟล์**: แยกไฟล์ Excel ที่อัปโหลด (`r.FormFile("excelFile")`)
- **ดึง columnName**: รับ columnName จาก form values (`r.FormValue("columnName")`)

#### 3. การตรวจสอบเบื้องต้น (ระดับ Handler)
- ✅ ตรวจสอบว่า `columnName` ว่างเปล่าหรือไม่
  - หากว่าง → ส่งคืน **400 Bad Request**
- ✅ ตรวจสอบว่า `columnName` เป็นหนึ่งในค่าที่อนุญาต
  - หากไม่อนุญาต → ส่งคืน **400 Bad Request**

#### 4. ประมวลผลไฟล์
- **อ่านไบต์ไฟล์**: อ่านเนื้อหาทั้งหมดของไฟล์เป็น byte slice (`excelFileBytes`)
- **มอบหมายให้ Service**: เรียก `CompareExcelWithDB` บน service

---

## 👨‍🍳 **[หัวหน้าครัว (Service: tools/compare/service.go)] - ตรรกะทางธุรกิจและการประสานงาน**

### วิธีการ CompareExcelWithDB:

#### 1. การตรวจสอบ columnName (ระดับ Service)
- ✅ ตรวจสอบว่า `columnName` ว่างเปล่าหรือไม่
- ✅ ตรวจสอบ `columnName` กับ `allowedColumns` map
- หากไม่ผ่าน → ส่งคืน error

#### 2. แยกวิเคราะห์ข้อมูล Excel
เรียกฟังก์ชันช่วย `readExcelColumnFromBytes(excelFileBytes, columnName)`:

**กระบวนการ:**
- ใช้ไลบรารี **excelize** เพื่อเปิดและอ่านไฟล์ Excel จาก byte slice
- ค้นหา `columnName` ที่ระบุและคอลัมน์ `hs_code` ในแถวหัวข้อของแผ่นงานแรก
- วนซ้ำผ่านแถว แยกค่าจากคอลัมน์เป้าหมายและ `hs_code` ที่เกี่ยวข้อง

**ผลลัพธ์:**
- ส่งคืน `map[string]ExcelValue`
- คีย์ = ค่าจากคอลัมน์เป้าหมาย
- Value = `ExcelValue` struct (ประกอบด้วยค่าและ HSCode ที่เกี่ยวข้อง)

#### 3. ดึงข้อมูลจากฐานข้อมูล
เรียก `GetValuesFromDB` บน repository (`s.repo.GetValuesFromDB`) โดยส่ง:
- Context
- columnName

#### 4. ประมวลผลผลลัพธ์จากฐานข้อมูล
Repository ส่งคืน `[]*DBDetails` และสร้างแผนที่สองแผนที่เพื่อการค้นหาที่มีประสิทธิภาพ:

**แผนที่ที่ 1 - dbValuesMap:**
- Type: `map[string]DBDetails`
- แมปค่าคอลัมน์เฉพาะ (เช่น goods_en, hs_code) กับ DBDetails struct ทั้งหมด

**แผนที่ที่ 2 - hsCodeMap:**
- Type: `map[string][]DBDetails`
- แมป HSCode กับ slice ของ DBDetails structs ที่ใช้ HSCode เดียวกัน

#### 5. เปรียบเทียบข้อมูล Excel กับข้อมูล DB
วนซ้ำผ่าน `excelVal` แต่ละตัวและสร้าง `ExcelItem`:

##### 🎯 การจับคู่โดยตรง (Direct Match):
- พยายามหาการจับคู่โดยตรงสำหรับ `excelVal.Value` ใน `dbValuesMap`
- หากพบ:
  - `IsMatch = true`
  - `MatchedBy = "column"`
  - `DBDetails` ชี้ไปยังเรคอร์ดฐานข้อมูลที่ตรงกัน

##### 🔄 การจับคู่ HSCode สำรอง (HSCode Fallback Match):
หากไม่มีการจับคู่โดยตรง และการเปรียบเทียบเป็นสำหรับ "goods_en" หรือ "goods_th" และ Excel item มี HSCode:

**ขั้นตอน 1:** ค้นหา `excelVal.HSCode` ใน `hsCodeMap`

**ขั้นตอน 2:** หากพบการจับคู่ HSCode ในฐานข้อมูล:
- **พยายามหาการจับคู่เฉพาะ** ภายในกลุม HSCode นั้น
  - ตรวจสอบว่า Excel value ตรงกับคอลัมน์ฐานข้อมูลที่เกี่ยวข้อง
  - หากพบ: `MatchedBy = "hs_code_specific_en"` หรือ `"hs_code_specific_th"`
- **หากไม่พบการจับคู่เฉพาะ**: ใช้เรคอร์ดแรกจากกลุม HSCode นั้นเป็นการจับคู่สำรอง
  - `MatchedBy = "hs_code_fallback"`

#### 6. สร้างผลลัพธ์
สร้าง `CompareResponse` struct ประกอบด้วย:
- **`TotalExcelRows`**: จำนวนรายการจากไฟล์ Excel
- **`TotalDBRows`**: จำนวนรายการเฉพาะที่โหลดจากฐานข้อมูลสำหรับการจับคู่โดยตรง
- **`MatchedRows`**: จำนวนรายการ Excel ที่ถูกจับคู่
- **`ExcelItems`**: slice ของ ExcelItem structs ที่ประมวลผลแล้ว

---

## 🥘 **[พ่อครัว (Repository: tools/compare/repository.go)] - ชั้นการเข้าถึงข้อมูล**

### วิธีการ GetValuesFromDB:

#### 1. รับการเชื่อมต่อฐานข้อมูล
- ดึงการเชื่อมต่อ PostgreSQL (`*pg.DB`) จากบริบท
- สร้างบริบทใหม่พร้อม timeout สำหรับ database query

#### 2. ตรวจสอบการมีอยู่ของคอลัมน์
- ตรวจสอบว่า `columnName` ที่ร้องขอมีอยู่จริงในตาราง `public.master_hs_code_v2`
- Query `information_schema.columns`
- หากไม่มี → ส่งคืน error

#### 3. สร้าง SQL Query
สร้าง SQL query แบบไดนามิกเพื่อเลือกฟิลด์ที่เกี่ยวข้องทั้งหมดจาก `public.master_hs_code_v2`

**เงื่อนไข WHERE:**
- `columnName` ไม่เป็น null และไม่ว่าง
- `hs_code` ไม่เป็น null และไม่ว่าง

#### 4. ดำเนินการ Query
- ดำเนินการ query ด้วย `db.WithContext(ctxQuery).Query(&dbValues, query)`
- สแกนผลลัพธ์ลงใน `var dbValues []*DBDetails`

#### 5. ส่งคืนผลลัพธ์
ส่งคืน slice ของ `*DBDetails` และ error (หากมี)

---

## 🏪 **[คลังสินค้า (Database: PostgreSQL)] - การจัดเก็บข้อมูล**

### โครงสร้างฐานข้อมูล:
ตาราง `public.master_hs_code_v2` เก็บข้อมูลหลักสำหรับ:
- **รหัส HS** (HS Codes)
- **รายละเอียดสินค้า** (ภาษาอังกฤษและไทย)
- **ภาษี** (Tariffs)
- **รหัสหน่วย** (Unit Codes)
- **อัตราภาษี** (Duty Rates)
- **ข้อมูลอื่นๆ** ที่เกี่ยวข้อง

Repository จะ query ตารางนี้เพื่อดึงข้อมูลที่ต้องการ

---

## 🍽️ **[เสิร์ฟ (Handler)] - ส่งผลลัพธ์กลับ**

### การจัดการข้อผิดพลาด:
หาก service ส่งข้อผิดพลาดกลับมา handler จะ:
- บันทึกข้อผิดพลาด (Log error)
- ส่ง **HTTP 500 Internal Server Error**
- หรือ **400 Bad Request** (หากข้อผิดพลาดเกิดจาก input ที่ไม่ถูกต้อง)

### กรณีสำเร็จ:
หากการเรียก service สำเร็จ:
1. ตั้งค่า `Content-Type` header เป็น `application/json`
2. เขียน HTTP status code **200 OK**
3. เข้ารหัส `CompareResponse` struct เป็น JSON
4. เขียนลงใน HTTP response body

---

## 📱 **[ลูกค้า (Client)] - รับผลลัพธ์**

### ผลลัพธ์ที่ได้รับ:
ลูกค้ารับ HTTP response:

#### ✅ กรณีสำเร็จ:
- **Status Code**: 200 OK
- **Content-Type**: application/json
- **Body**: ผลลัพธ์การเปรียบเทียบในรูปแบบ JSON

#### ❌ กรณีเกิดข้อผิดพลาด:
- **Status Code**: 400 Bad Request หรือ 500 Internal Server Error
- **Body**: ข้อความแสดงข้อผิดพลาดที่เหมาะสม

---

## 📊 **ตัวอย่างผลลัพธ์ JSON Response**

```json
{
  "total_excel_rows": 100,
  "total_db_rows": 5000,
  "matched_rows": 85,
  "excel_items": [
    {
      "value": "Rice",
      "hs_code": "1006.30",
      "is_match": true,
      "matched_by": "column",
      "db_details": {
        "hs_code": "1006.30",
        "goods_en": "Rice",
        "goods_th": "ข้าว",
        "tariff_rate": "5%"
      }
    }
  ]
}
```

---

## 🔄 **สรุปการไหลของข้อมูล**

การไหลแบบละเอียดนี้ครอบคลุมการเดินทางของคำขอผ่านชั้นต่างๆ ของ API:

1. **Client** ส่งคำขอพร้อมไฟล์ Excel และพารามิเตอร์
2. **Handler** รับคำขอ ตรวจสอบเบื้องต้น และส่งต่อ
3. **Service** ประมวลผลตรรกะทางธุรกิจ แยกวิเคราะห์ Excel และเปรียบเทียบข้อมูล
4. **Repository** เข้าถึงฐานข้อมูลและดึงข้อมูลที่จำเป็น
5. **Database** ให้บริการข้อมูลจากตาราง master
6. **ส่งผลลัพธ์กลับ** ผ่านชั้นต่างๆ จนถึงลูกค้า

แต่ละชั้นมีหน้าที่และความรับผิดชอบที่ชัดเจน ทำให้ระบบมีการแยกส่วนที่ดีและง่ายต่อการบำรุงรักษา