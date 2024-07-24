package guid

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pemako/gopkg/guid/snowflake"
)

type Guid struct {
	IDGenerator *snowflake.Node
}

var g *Guid

func init() {
	g = &Guid{}
	node, err := snowflake.NewNode(int64(rand.Intn(1024)))
	if err != nil {
		fmt.Println("get snowflake node error, msg: ", err)
		return
	}
	g.IDGenerator = node
}

func GetInt64(ctx context.Context) (r int64, err error) {
	if g.IDGenerator == nil {
		return 0, fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Int64(), nil
}

func GetString(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.String(), nil
}

func GetBase2(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Base2(), nil
}

func GetBase32(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Base32(), nil
}

func GetBase36(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Base36(), nil
}

func GetBase58(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Base58(), nil
}

func GetBase64(ctx context.Context) (r string, err error) {
	if g.IDGenerator == nil {
		return "", fmt.Errorf("get id generator error")
	}
	gen := g.IDGenerator.Generate()
	return gen.Base64(), nil
}
