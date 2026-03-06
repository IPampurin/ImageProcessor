package service

type ResizeParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ImageTask struct {
	ImageID      string        `json:"image_id"`
	ObjectKey    string        `json:"object_key"`
	Bucket       string        `json:"bucket"`
	Thumbnail    bool          `json:"thumbnail"`
	Watermark    bool          `json:"watermark"`
	Resize       *ResizeParams `json:"resize,omitempty"`
	OriginalName string        `json:"original_name,omitempty"`
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
