package utils

import (
	"io"
	"log"
)

//Review ! rktkov: Зачем создан этот пакет, если в большинстве мест, где осуществляется логгирование, используется базовый пакет log.
//Либо надо удалять этот пакет, либо оборачивать до конца использование пакета log и вызывать методы этого пакете во всем коде.
//Review Answer ! anton: I use this package because this implementation was recommended another mentor.
func LogErr(n int, err error) int {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}

	return n
}

func ResourceClose(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Can't close resource: %v", err)
	}
}
