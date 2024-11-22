package support

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/1password/onepassword-sdk-go"
	"github.com/joho/godotenv"
)

var (
	defenv   = ".env"
	usrenv   = ".env.local"
	once     sync.Once
	opclient *onepassword.Client
)

// load global .env files
func init() {
	// .env (default)
	godotenv.Load(defenv)

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(usrenv)
}

func Environ(dir string, env ...string) []string {
	osenv := os.Environ()
	dotenv := dotenvFiles(dir)

	return slices.Concat(osenv, dotenv, env)
}

func dotenvFiles(dir string) (res []string) {
	filenames := []string{dir + "/" + defenv, dir + "/" + usrenv}

	for _, filename := range filenames {
		src, err := os.ReadFile(filename)
		if err != nil {
			return
		}

		envMap, err := godotenv.UnmarshalBytes(src)
		if err != nil {
			return
		}

		for k, v := range envMap {
			res = append(res, k+"="+v)
		}
	}

	return
}

func op(ref string) string {
	opEnabled, _ := strconv.ParseBool(os.Getenv("OP_ENABLED"))
	if !opEnabled || !strings.HasPrefix(ref, "op://") {
		return ref
	}

	if opclient == nil {
		once.Do(initop)
	}

	item, err := opclient.Secrets.Resolve(context.Background(), ref)
	if err != nil {
		return ref
	}

	return item
}

func initop() {
	var err error
	opclient, err = onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(os.Getenv("OP_SERVICE_ACCOUNT_TOKEN")),
		onepassword.WithIntegrationInfo("luna op integration", "v0.1.0"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
