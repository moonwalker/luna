package support

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/1password/onepassword-sdk-go"
	"github.com/joho/godotenv"
)

var (
	optoken  string
	opclient *onepassword.Client
)

func init() {
	// .env (default)
	godotenv.Load()

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(".env.local")

	if os.Getenv("OP_ENABLED") != "" {
		var err error
		opclient, err = onepassword.NewClient(
			context.Background(),
			onepassword.WithServiceAccountToken(os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")),
			onepassword.WithIntegrationInfo("luna op integration", "v0.1.0"),
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func Environ(env ...string) []string {
	environ := append(os.Environ(), env...)

	for i, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		environ[i] = pair[0] + "=" + op(pair[1])
	}

	return environ
}

func op(ref string) string {
	if opclient == nil || !strings.HasPrefix(ref, "op://") {
		return ref
	}

	item, err := opclient.Secrets.Resolve(context.Background(), ref)
	if err != nil {
		return ref
	}

	return item
}
