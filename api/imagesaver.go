package api

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Функция для загрузки изображения в Amazon S3 и возврата ссылки
func uploadImageToS3(fileHeader *multipart.FileHeader) (string, error) {
	// Настройте сессию Amazon S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Укажите ваш регион
	})
	if err != nil {
		return "", err
	}

	// Создайте новый клиент Amazon S3
	svc := s3.New(sess)

	// Откройте файл из заголовка
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Создайте уникальное имя файла на основе времени
	fileName := fmt.Sprintf("images/%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Загрузите файл в Amazon S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("your-s3-bucket-name"), // Замените на имя вашего бакета
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"), // Установите права доступа, как требуется
	})
	if err != nil {
		return "", err
	}

	// Сформируйте URL для доступа к загруженному изображению
	imageURL := fmt.Sprintf("https://your-s3-bucket-name.s3.amazonaws.com/%s", fileName) // Замените на имя вашего бакета

	return imageURL, nil
}
