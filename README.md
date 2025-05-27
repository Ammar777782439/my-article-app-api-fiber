# ميزة: إضافة نماذج المؤلفين والعلاقات

يقدم هذا الفرع ميزة جديدة لإدارة المؤلفين ككيانات منفصلة وربطهم بالمقالات. لقد تم تحديث هيكل البيانات وواجهات برمجة التطبيقات لتعكس هذه العلاقة.

## **🚀 الميزات الجديدة**

* **إدارة المؤلفين:** إضافة عمليات CRUD كاملة للمؤلفين (Create, Read, Update, Delete).
* **علاقة "متعدد إلى واحد" (Many-to-One):** ربط كل مقال بمؤلف واحد محدد.
* **تحميل مسبق للبيانات (Preloading):** عند جلب المقالات، يتم الآن تلقائيًا تضمين بيانات المؤلف المرتبط، مما يقلل من عدد استعلامات قاعدة البيانات.

## **🛠️ التغييرات الرئيسية في الكود**

* **`internal/models/author.go` (جديد):**
    * تعريف بنية `Author` مع حقول مثل `Name` و `Email`.
    * تضمين `gorm.Model` للمفاتيح الأساسية وحقول التواريخ.
    * إضافة حقل `Articles []Article` لتحديد علاقة "واحد إلى متعدد" من جانب `Author`.
* **`internal/models/article.go` (مُعدّل):**
    * تم إزالة حقل `Author` (سلسلة نصية).
    * إضافة `AuthorID uint` ليكون المفتاح الأجنبي الذي يشير إلى `Author`.
    * إضافة `Author Author` مع علامة `gorm:"foreignKey:AuthorID"` لتحديد العلاقة الصريحة.
    * إضافة `validate:"-"` إلى حقل `Author` في بنية `Article` لتجنب التحقق من صحة كائن `Author` المتداخل عند إرسال `AuthorID` فقط في طلبات المقالات.
* **`internal/repository/author_repository.go` (جديد):**
    * تنفيذ طرق CRUD لنموذج `Author` (على غرار `ArticleRepository`).
    * استخدام `Preload("Articles")` في `FindByID` لجلب مقالات المؤلف عند الاستعلام عنه.
* **`internal/repository/article_repository.go` (مُعدّل):**
    * تم تعديل دوال `FindAll()` و `FindByID()` لاستخدام `Preload("Author")`، مما يضمن تحميل بيانات المؤلف تلقائيًا مع المقالات.
* **`internal/handlers/author_handler.go` (جديد):**
    * تنفيذ معالجات HTTP لعمليات CRUD للمؤلفين (على غرار `ArticleHandler`).
    * تضمين التحقق من صحة البيانات لـ `Name` و `Email` للمؤلفين.
* **`internal/handlers/article_handler.go` (مُعدّل):**
    * تعديل `CreateArticle` لقبول `AuthorID` بدلاً من حقل `Author` النصي.
    * التعامل مع الاستجابات التي ستتضمن الآن كائن `Author` المتداخل في المقالات.
* **`internal/database/gorm.go` (مُعدّل):**
    * تم تحديث `db.AutoMigrate()` ليشمل نموذج `models.Author{}`، مما يضمن إنشاء جدول `authors` في قاعدة البيانات.
* **`cmd/api/main.go` (مُعدّل):**
    * تهيئة `AuthorRepository` و `AuthorHandler` في وظيفة `main`.
    * إضافة مجموعة مسارات جديدة `/api/v1/authors` لمعالجة طلبات المؤلفين.

## **🚀 واجهات برمجة التطبيقات (API Endpoints) الجديدة للمؤلفين**

بالإضافة إلى مسارات المقالات الموجودة، تم تقديم المسارات التالية لإدارة المؤلفين:

| الطريقة | المسار          | الوصف                          | جسم الطلب (Request Body)                   |
| :---- | :-------------- | :----------------------------- | :----------------------------------------- |
| POST  | `/api/v1/authors` | **إنشاء** مؤلف جديد.            | `{"name": "...", "email": "..."}`         |
| GET   | `/api/v1/authors` | **جلب جميع** المؤلفين.         | لا يوجد                                    |
| GET   | `/api/v1/authors/:id` | **جلب** مؤلف واحد حسب المعرف (ID). | لا يوجد                                    |
| PUT   | `/api/v1/authors/:id` | **تحديث** مؤلف موجود حسب المعرف (ID). | `{"name": "...", "email": "..."}`         |
| DELETE | `/api/v1/authors/:id` | **حذف** مؤلف حسب المعرف (ID).  | لا يوجد                                    |

## **📝 كيفية الاختبار**

1.  **تشغيل قاعدة البيانات:**
    ```bash
    docker-compose up -d
    ```
2.  **تثبيت التبعيات:**
    ```bash
    go mod tidy
    ```
3.  **تشغيل التطبيق:**
    ```bash
    cd cmd/api
    go run main.go
    ```
4.  **استخدام `curl` أو Postman:**
    * **أولاً، أنشئ مؤلفًا:**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"name": "أحمد جمال", "email": "ahmed.gamal@example.com"}' http://localhost:3000/api/v1/authors
        ```
        احفظ الـ `ID` الخاص بالمؤلف الذي تم إنشاؤه من الاستجابة.
    * **ثم، أنشئ مقالًا باستخدام `author_id` الذي حصلت عليه:**
        ```bash
        curl -X POST -H "Content-Type: application/json" -d '{"title": "مقدمة إلى Go Fiber", "content": "مقال رائع يشرح أساسيات Go Fiber.", "author_id": 1}' http://localhost:3000/api/v1/articles
        ```
    * **جلب المقالات للتحقق من تضمين المؤلف:**
        ```bash
        curl http://localhost:3000/api/v1/articles
        ```
        يجب أن ترى بيانات المؤلف متداخلة داخل كل كائن مقال.