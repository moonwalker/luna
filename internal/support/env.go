package support

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/1password/onepassword-sdk-go"
	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	opclient *onepassword.Client
)

func init() {
	loadEnvFiles("")
}

func Environ(dir string, env ...string) []string {
	loadEnvFiles(dir)

	environ := append(os.Environ(), env...)
	for i, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		environ[i] = pair[0] + "=" + op(pair[1])
	}

	return environ
}

func loadEnvFiles(dir string) {
	defenv := ".env"
	usrenv := ".env.local"

	if dir != "" {
		defenv = dir + "/" + defenv
		usrenv = dir + "/" + usrenv
	}

	// .env (default)
	godotenv.Load(defenv)

	// .env.local # local user specific (usually git ignored)
	godotenv.Overload(usrenv)
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
