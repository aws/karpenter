package test

import (
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"

	f "github.com/ellistarn/karpenter/pkg/utils/functional"
	"github.com/ellistarn/karpenter/pkg/utils/project"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

var (
	YAMLDocumentDelimiter = regexp.MustCompile(`\n---\n`)
)

type Namespace v1.Namespace

// Returns a test namespace
func NewNamespace() *Namespace {
	return &Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(randomdata.SillyName()),
		},
	}
}

// Instantiates a test resource from YAML
func (n *Namespace) ParseResource(path string, object runtime.Object) error {
	data, err := ioutil.ReadFile(project.RelativeToRoot(path))
	if err != nil {
		return errors.Wrapf(err, "reading file %s", path)
	}
	if err := parseFromYaml(data, object); err != nil {
		return errors.Wrapf(err, "parsing yaml")
	}

	if field := reflect.ValueOf(object).Elem().FieldByName("Namespace"); field.IsValid() {
		field.SetString(n.Name)
	}
	return nil
}

func parseFromYaml(data []byte, object runtime.Object) error {
	errs := []error{}
	for _, document := range YAMLDocumentDelimiter.Split(string(data), -1) {
		if err := yaml.UnmarshalStrict([]byte(document), object); err != nil {
			errs = append(errs, err)
		} else {
			return nil
		}
	}
	return errors.Wrap(f.FirstNonNilError(errs), "parsing YAML")
}
