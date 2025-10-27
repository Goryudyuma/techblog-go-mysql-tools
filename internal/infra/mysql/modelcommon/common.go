package modelcommon

import (
	"github.com/Goryudyuma/techblog-go-mysql-tools/internal/infra/mysql/modelutil"
)

// Common
// 自動生成されるモデルに、共通の処理を追加するための構造体
type Common struct {
}

func (m *Common) ListDBTagString() []string {
	return modelutil.ListDBTagString(m)
}
