package benchmark

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/lithammer/shortuuid/v3"
	uuid "github.com/satori/go.uuid"
	"uniqueid/snowflake"
)

// go test -bench=BenchmarkUuid -benchmem
// go test -bench=. -benchmem
func BenchmarkUuid1(b *testing.B) {
	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(uuid.Must(uuid.NewV1()).String(), "-", "")
	}
}

func BenchmarkUuid2(b *testing.B) {
	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(uuid.Must(uuid.NewV2(uuid.DomainPerson)).String(), "-", "")
	}
}

func BenchmarkUuid3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(uuid.NewV3(uuid.NamespaceDNS, "www.example.com").String(), "-", "")
	}
}

func BenchmarkUuid4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
	}
}

func BenchmarkUuid5(b *testing.B) {
	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(uuid.NewV5(uuid.NamespaceDNS, "www.example.com").String(), "-", "")
	}
}

func BenchmarkSnowflake(b *testing.B) {
	//b.ResetTimer()
	sf := snowflake.NewSnowflake(0)
	for i := 0; i < b.N; i++ {
		sf.Generate()
	}
}

func BenchmarkShortUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		shortuuid.New()
	}
}

func BenchmarkJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(struct {
			T1  string
			T2  string
			T3  string
			T4  string
			T5  string
			T6  string
			T7  string
			T8  string
			T9  string
			T10 string
		}{
			T1:  "",
			T2:  "",
			T3:  "",
			T4:  "",
			T5:  "",
			T6:  "",
			T7:  "",
			T8:  "",
			T9:  "",
			T10: "",
		})
	}
}
func BenchmarkBase64Random128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandom()
	}
}

// GenerateRandom 生成随机GUID
func GenerateRandom() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	//return DefaultEncoder.Encode(buf), nil
	return base64.RawStdEncoding.EncodeToString(buf), nil
}
