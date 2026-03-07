package domain

const (
	thumbWidht  = 150 // ширина миниатюры по умолчанию
	thumbHeight = 150 // высота миниатюры по умолчанию
)

// ResizeParams описывает требуемые ширину и высоту изображения
type ResizeParams struct {
	Width  int // ширина
	Height int // высота
}

type ImageTask struct {
	ImageID      string        // UUID файла
	ObjectKey    string        // ключ для S3 хранилища
	Bucket       string        // имя бакета в S3 хранилище
	Thumbnail    bool          // требуется ли миниатюра
	Watermark    bool          // нужен ли водяной знак
	Resize       *ResizeParams // требуемые размеры изображения
	OriginalName string        // имя файла
}

/*
task := ImageTask{
    ImageID:   "f47e1b3c-8e5a-4b7d-9c2a-3d1e6f8a9b0c",
    ObjectKey: "images/f47e1b3c-8e5a-4b7d-9c2a-3d1e6f8a9b0c.jpg",
    Bucket:    "original-images",
    Thumbnail: true,
    Watermark: false,
    Resize:    &ResizeParams{Width: 800, Height: 600},
    OriginalName: "vacation.jpg",
}
payload, _ := json.Marshal(task)

outbox := Outbox{
    ID:        uuid.New(),
    Topic:     "image-tasks",
    Key:       task.ImageID, // используем тот же UUID
    Payload:   payload,
    CreatedAt: time.Now(),
}
*/
