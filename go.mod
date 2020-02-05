module github.com/flant/addon-operator

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/flant/libjq-go v1.0.1-0.20200205115921-27e93c18c17f // branch: master
	github.com/flant/shell-operator v1.0.0-beta.7.0.20200205180330-22622f8aa612 // branch: master
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/go-openapi/spec v0.19.3
	github.com/kennygrant/sanitize v1.2.4
	github.com/onsi/gomega v1.7.0
	github.com/otiai10/copy v1.0.1
	github.com/otiai10/curr v0.0.0-20190513014714-f5a3d24e5776 // indirect
	github.com/peterbourgon/mergemap v0.0.0-20130613134717-e21c03b7a721
	github.com/prometheus/client_golang v1.0.0
	github.com/segmentio/go-camelcase v0.0.0-20160726192923-7085f1e3c734
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/satori/go.uuid.v1 v1.2.0
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	sigs.k8s.io/yaml v1.1.1-0.20191128155103-745ef44e09d6 // branch master, commit 745ef44e09d6, with fixes in yaml.v2.2.7
)

replace github.com/go-openapi/validate => github.com/flant/go-openapi-validate v0.19.4-0.20190926112101-38fbca4ac77f // branch: fix_in_body

//replace github.com/flant/shell-operator => ../shell-operator
