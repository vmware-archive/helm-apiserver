package server

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/apiserver"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	super "github.com/kubernetes-incubator/apiserver-builder/pkg/cmd/server"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	openapi "k8s.io/kube-openapi/pkg/common"

	"github.com/bitnami-labs/helm-apiserver/pkg/admission/plugin/impersonate"
)

// This is apiserver-builder/pkg/cmd/server, with additional support
// for admission controllers.  Unfortunately upstream's functions do
// not support injecting a customised config, so there's quite a bit
// of copy+paste here :(
//
// Hopefully much of this will go away with kubernetes/kubernetes#56627
//

// StartApiServer starts an apiserver hosting the provider apis and openapi definitions.
func StartApiServer(etcdPath string, apis []*builders.APIGroupBuilder, openapidefs openapi.GetOpenAPIDefinitions, title, version string) {
	logs.InitLogs()
	defer logs.FlushLogs()

	super.GetOpenApiDefinition = openapidefs

	// To disable providers, manually specify the list provided by getKnownProviders()
	cmd, _ := NewCommandStartServer(etcdPath, os.Stdout, os.Stderr, apis, wait.NeverStop, title, version)
	if logflag := flag.CommandLine.Lookup("v"); logflag != nil {
		level := logflag.Value.(*glog.Level)
		levelPtr := (*int32)(level)
		cmd.Flags().Int32Var(levelPtr, "loglevel", 0, "Set the level of log output")
	}
	cmd.Flags().AddFlagSet(pflag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

// ServerOptions is apiserver-builder's ServerOptions, with support
// for admission controllers.
type ServerOptions struct {
	super.ServerOptions

	Admission *genericoptions.AdmissionOptions
}

func NewServerOptions(etcdPath string, out, errOut io.Writer, b []*builders.APIGroupBuilder) *ServerOptions {
	parent := super.NewServerOptions(etcdPath, out, errOut, b)
	if parent == nil {
		panic("NewServerOptions returned nil")
	}
	o := ServerOptions{
		ServerOptions: *parent,
		Admission:     genericoptions.NewAdmissionOptions(),
	}
	return &o
}

// NewCommandStartMaster provides a CLI handler for 'start master' command
func NewCommandStartServer(etcdPath string, out, errOut io.Writer, builders []*builders.APIGroupBuilder,
	stopCh <-chan struct{}, title, version string) (*cobra.Command, *ServerOptions) {
	o := NewServerOptions(etcdPath, out, errOut, builders)

	// Support overrides
	cmd := &cobra.Command{
		Short: "Launch an API server",
		Long:  "Launch an API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(stopCh, title, version); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.PrintBearerToken, "print-bearer-token", false,
		"Print a curl command with the bearer token to test the server")
	flags.BoolVar(&o.PrintOpenapi, "print-openapi", false,
		"Print the openapi json and exit")
	flags.BoolVar(&o.RunDelegatedAuth, "delegated-auth", true,
		"Setup delegated auth")
	//flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "Kubeconfig of apiserver to talk to.")
	o.RecommendedOptions.AddFlags(flags)
	o.Admission.AddFlags(flags)
	return cmd, o
}

func (o ServerOptions) Validate(args []string) error {
	errors := []error{}
	if err := o.ServerOptions.Validate(args); err != nil {
		errors = append(errors, err)
	}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.Admission.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o ServerOptions) Config() (*apiserver.Config, error) {
	impersonate.Register(o.Admission.Plugins)

	config, err := o.ServerOptions.Config()
	if err != nil {
		return nil, err
	}

	extconf, err := clientcmd.BuildConfigFromFlags("", o.Kubeconfig)
	extclient, err := kubernetes.NewForConfig(extconf)
	if err != nil {
		return nil, err
	}

	glog.Infof("Creating shared informer factory from kubeconfig %s", o.Kubeconfig)
	sharedInformerFactory := informers.NewSharedInformerFactory(extclient, 11*time.Hour)

	if err := o.Admission.ApplyTo(config.GenericConfig, sharedInformerFactory); err != nil {
		return nil, err
	}

	return config, nil
}

func (o *ServerOptions) RunServer(stopCh <-chan struct{}, title, version string) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	if o.PrintBearerToken {
		glog.Infof("Serving on loopback...")
		glog.Infof("\n\n********************************\nTo test the server run:\n"+
			"curl -k -H \"Authorization: Bearer %s\" %s\n********************************\n\n",
			config.GenericConfig.LoopbackClientConfig.BearerToken,
			config.GenericConfig.LoopbackClientConfig.Host)
	}
	o.BearerToken = config.GenericConfig.LoopbackClientConfig.BearerToken

	for _, provider := range o.APIBuilders {
		config.AddApi(provider)
	}

	config.Init()

	config.GenericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(super.GetOpenApiDefinition, builders.Scheme)
	config.GenericConfig.OpenAPIConfig.Info.Title = title
	config.GenericConfig.OpenAPIConfig.Info.Version = version

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	for _, h := range o.PostStartHooks {
		server.GenericAPIServer.AddPostStartHook(h.Name, h.Fn)
	}

	s := server.GenericAPIServer.PrepareRun()
	err = validators.OpenAPI.SetSchema(readOpenapi(server.GenericAPIServer.Handler))
	if o.PrintOpenapi {
		fmt.Printf("%s", validators.OpenAPI.OpenApi)
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	s.Run(stopCh)

	return nil
}

func readOpenapi(handler *genericapiserver.APIServerHandler) string {
	req, err := http.NewRequest("GET", "/swagger.json", nil)
	if err != nil {
		panic(fmt.Errorf("Could not create openapi request %v", err))
	}
	resp := &super.BufferedResponse{}
	handler.ServeHTTP(resp, req)
	return resp.String()
}
